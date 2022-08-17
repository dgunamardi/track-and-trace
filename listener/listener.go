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

	"go.mongodb.org/mongo-driver/bson"
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

type TransactionData struct {
	Event_Id      string `json:"event_id"`
	Event_Type    int32  `json:"event_type"`
	Input_GTIN    string `json:"input_gtin"`
	Ouput_GTIN    string `json:"output_gtin"`
	Serial_Number string `json:"serial_number"`
	Event_Time    string `json:"event_time"`
	Event_Loc     string `json:"event_loc"`
	Location_Name string `json:"location_name"`
	Company_Name  string `json:"company_name"`
}

var (
	listenArgs = ListenArgs{
		SeekType:   seek.Newest,
		StartBlock: 340,
	}

	//parsedBlock parser.Block

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

	for events := range blockEvents {
		if skipEvent {
			skipEvent = false
			continue
		}

		parsedBlock := parser.Block{}
		parsedBlock.Init(events.Block)

		log.Println(events.Block.GetHeader().GetNumber())

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

						docBson := &TransactionData{}
						err = bson.UnmarshalExtJSON(kvWrite.Value, true, docBson)
						if err != nil {
							panic(fmt.Errorf("failed to unmarshal json to bson: %v", err))
						}

						//log.Printf("\nDocBson:\n%v", docBson)

						if docBson.Event_Type != 0 {
							result, err := coll.InsertOne(ctx.TODO(), docBson)
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
