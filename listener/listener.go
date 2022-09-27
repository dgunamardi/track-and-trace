package main

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"os"
	"strconv"

	cfg "earhart.com/config"
	parser "earhart.com/parser"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
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
	URI         string
	dbClient    *mongo.Client
	Name        string
	Collections []string
}

var (
	listenArgs = ListenArgs{
		SeekType:   seek.Newest,
		StartBlock: 0,
	}

	dbVars = DBVars{
		URI:  "mongodb://localhost:27017/food_safety",
		Name: "food_sagety",
		Collections: []string{
			"track_trace",
			"import",
		},
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

func ListenToBlockEvents(channelProvider context.ChannelProvider) {
	evClient, err := eventClient.New(
		channelProvider,
		eventClient.WithBlockEvents(),
		eventClient.WithSeekType(listenArgs.SeekType),
		eventClient.WithBlockNum(listenArgs.StartBlock),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create event client: %v", err))
	}

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

		blockNumber := events.Block.GetHeader().GetNumber()
		log.Println(blockNumber)

		parsedBlock := parser.Block{}
		parsedBlock.Init(events.Block)

		/*
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
		*/

		// === UNWRAP THE ENVELOPE ===
		envelopes := parsedBlock.BlockData.Envelopes
		for _, envelope := range envelopes {

			// * SKIP IF NOT TX
			txType := common.HeaderType(envelope.Payload.Header.ChannelHeader.ChannelHeaderProto.GetType())
			if txType != common.HeaderType_ENDORSER_TRANSACTION {
				continue
			}

			// * FOR EVERY ACTION
			txActions := envelope.Payload.Transaction.TransactionActions
			for _, txAction := range txActions {

				// * CHECK INVOCATION SPEC FOR FUNCTION NAME TO DETERMINE COLLECTION / OBJECT TYPE
				invocationSpec := txAction.ChaincodeActionPayload.ChaincodeProposalPayload.ChaincodeInvocationSpec.ChaincodeInvocationSpecProto
				collectionName := GetCollectionName(invocationSpec.GetChaincodeSpec())

				log.Println(collectionName)

				nsRWSets := txAction.ChaincodeActionPayload.ChaincodeEndorsedAction.ProposalResponsePayload.Extension.Results.NsReadWriteSets
				for _, nsRWSet := range nsRWSets {
					// * SKIP IF LSCC
					if nsRWSet.Namespace == "lscc" {
						continue
					}
					kvWrites := nsRWSet.RWSet.KVWrites
					InsertToDB(kvWrites, collectionName)
				}
			}
		}
	}
}

func GetCollectionName(ccSpec *peer.ChaincodeSpec) (collectionName string) {
	fcnName := string(ccSpec.GetInput().GetArgs()[0])
	if strings.Contains(fcnName, "IMP") {
		return "import"
	}
	if strings.Contains(fcnName, "TNT") {
		return "track_trace"
	}
	return ""
}

func InsertToDB(kvWrites []parser.KVWrite, collectionName string) {

	//coll := dbClient.Database(dbVars.Name).Collection(collectionName)

	for _, kvWrite := range kvWrites {

		var data parser.ObjectData

		if collectionName == "track_trace" {
			data = &parser.EventData{}
			err := json.Unmarshal(kvWrite.Value, data)
			if err != nil {
				panic(err)
			}
		}

		if collectionName == "import" {
			data = &parser.ImportData{}
			err := json.Unmarshal(kvWrite.Value, data)
			if err != nil {
				panic(err)
			}
		}

		if data.IsValid() {
			log.Println(data)

			/*
				result, err := coll.InsertOne(ctx.TODO(), data)
				if err != nil {
					panic(fmt.Errorf("failed to insert document to collection: %v", err))
				}

				log.Printf("Inserted document with _id: %v\n", result.InsertedID)
			*/
		}
	}

}
