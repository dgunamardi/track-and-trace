package main

import (
	ctx "context"
	"fmt"
	"log"

	"os"
	"strconv"

	cfg "earhart.com/config"
	parser "earhart.com/parser"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	eventClient "github.com/hyperledger/fabric-sdk-go/pkg/client/event"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListenArgs struct {
	SeekType   seek.Type
	StartBlock uint64
}

type DBVars struct {
	URI      string
	dbClient *mongo.Client
}

var (
	listenArgs = ListenArgs{
		SeekType:   seek.FromBlock,
		StartBlock: 2001,
	}

	dbVars = DBVars{
		URI: "mongodb://localhost:27017/track_trace",
	}
)

func main() {
	cfg.LoadConfig()
	cfg.InitializeSDK()
	cfg.InitializeUserIdentity()

	session := cfg.Sdk.Context(fabsdk.WithIdentity(cfg.User))
	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(session, cfg.CVars.ChannelId)
	}

	args := os.Args[1:]
	SetListenerArgs(args)

	ListenToBlockEvents(channelProvider)
}

func ConnectToDB() {
	client, err := mongo.Connect(ctx.TODO(), options.Client().ApplyURI(dbVars.URI))
	if err != nil {
		panic(fmt.Errorf("failed to create client: %v", err))
	}
	dbVars.dbClient = client

}

func SetListenerArgs(args []string) {
	if len(args) == 0 || len(args) > 2 {
		listenArgs.SeekType = seek.Newest
		return
	}
	switch args[0] {
	case "oldest":
		listenArgs.SeekType = seek.Oldest
	case "newest":
		listenArgs.SeekType = seek.Newest
	case "from":
		if args[1] == "" {
			panic("not enough arguments. 'from' should be followed by a number indicating the starting block\n")
		}
		listenArgs.SeekType = seek.FromBlock

		sb, err := strconv.Atoi(args[1])
		if err != nil {
			panic(fmt.Errorf("error in arg to int conversion: %v", err))
		}
		listenArgs.StartBlock = uint64(sb)
	default:
		listenArgs.SeekType = seek.Newest
	}
}

// CHECKPOINT:
// - GTIN IN OUT: 833
// - INVOKE TESTING:
// - PROPER CP:
//

func ListenToBlockEvents(channelProvider context.ChannelProvider) {

	// connect to fabric ahnnel
	evClient, err := eventClient.New(
		channelProvider,
		eventClient.WithBlockEvents(),
		eventClient.WithSeekType(listenArgs.SeekType),
		eventClient.WithBlockNum(listenArgs.StartBlock),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create event client: %v", err))
	}

	// register events
	eventRegister, blockEvents, err := evClient.RegisterBlockEvent()
	defer evClient.Unregister(eventRegister)

	log.Println("--- start listening to events ---")

	// skip event once when seek.newest is called to prevent duplicate of latest block to db in case the service is restart
	skipEvent := false
	if listenArgs.SeekType == seek.Newest {
		skipEvent = true
	}
	//skipEvent = false

	for events := range blockEvents {
		if skipEvent {
			skipEvent = false
			continue
		}

		blockNumber := events.Block.GetHeader().GetNumber()
		/*
			if blockNumber > 2001 {
				break
			}
		*/
		log.Println(blockNumber)

		parsedBlock := parser.Block{}
		parsedBlock.Init(events.Block)

		// connect to mongo DB
		dbClient, err := mongo.Connect(ctx.TODO(), options.Client().ApplyURI(dbVars.URI))
		if err != nil {
			panic(fmt.Errorf("failed to create client: %v", err))
		}
		defer func() {
			if err = dbClient.Disconnect(ctx.TODO()); err != nil {
				panic(err)
			}
		}()

		// get collection
		coll := dbClient.Database("track_trace").Collection("event")

		// unwrap
		envelopes := parsedBlock.BlockData.Envelopes
		for _, envelope := range envelopes {
			txActions := envelope.Payload.Transaction.TransactionActions
			for _, txAction := range txActions {
				nsRWSets := txAction.ChaincodeActionPayload.ChaincodeEndorsedAction.ProposalResponsePayload.Extension.Results.NsReadWriteSets
				for _, nsRWSet := range nsRWSets {
					kvWrites := nsRWSet.RWSet.KVWrites

					for _, kvWrite := range kvWrites {
						//log.Println("kvWrite:\n" + kvWrite.ValueString)

						var eventData parser.EventData
						eventData.Populate(kvWrite.Value)

						if eventData.EventType != 0 {
							log.Println(eventData)

							result, err := coll.InsertOne(ctx.TODO(), eventData)
							if err != nil {
								panic(fmt.Errorf("failed to insert document to collection: %v", err))
							}

							log.Printf("Inserted document with _id: %v\n", result.InsertedID)
						}

					}

				}
			}
		}
	}
}
