package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	cfg "earhart.com/config"
	"earhart.com/parser"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	clientChannel "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

//TODO
// Data Array Enum

func main() {
	cfg.LoadConfig()
	cfg.InitializeSDK()
	cfg.InitializeUserIdentity()

	session := cfg.Sdk.Context(fabsdk.WithIdentity(cfg.User))

	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(session, cfg.CVars.ChannelId)
	}

	args := os.Args[1:]
	Invoke(channelProvider, args)
}

func Invoke(channelProvider context.ChannelProvider, args []string) {
	client, err := clientChannel.New(channelProvider)
	if err != nil {
		panic(fmt.Errorf("failed to create channel client: %v", err))
	}

	if len(args) == 0 {
		panic("argument cannot be empty")
	}

	switch args[0] {
	case "-q":
		GetOwnerCredit(client, args[1])
	case "-se":
		SubmitData(client, args[1], parser.EVENT_DATA)
	case "-si":
		SubmitImportData(client, args[1], parser.IMPORT_DATA)
	case "-sp":
		SubmitProductData(client, args[1], parser.PRODUCT_DATA)
	case "-sr":
		SubmitRecallData(client, args[1], parser.RECALL_DATA)
	default:
		panic("argument is not available. Available Arguments: -q, -se, -si, -sp, -sr'\n")
	}
}

// csv version
func SubmitData(client *clientChannel.Client, args string, objectType parser.ObjectType) {
	dataset := parser.CSVToData(args, objectType)

	var fcn string

	switch objectType {
	case parser.EVENT_DATA:
		fcn = "AddTNTData"
	case parser.IMPORT_DATA:
		fcn = "AddIMPData"
	}

	for _, data := range dataset {
		key := uuid.New().String()
		accId := fmt.Sprint(data.GetId())

		dataJson, _ := json.Marshal(data)
		dataString := string(dataJson)

		fmt.Println(dataString)
		byteArgs := [][]byte{
			[]byte(key),
			[]byte(accId),
			[]byte(dataString),
		}

		SubmitTransaction(client, byteArgs, fcn)
	}
}

func SubmitProductData(client *clientChannel.Client, args string, objectType parser.ObjectType) {
	jsonFile, err := os.Open(args)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteArray, _ := ioutil.ReadAll(jsonFile)

	var products []parser.ProductData
	json.Unmarshal(byteArray, &products)

	fcn := "AddPROData"

	for _, data := range products {
		key := uuid.New().String()
		accId := fmt.Sprint(data.GetId())

		dataJson, _ := json.Marshal(data)
		dataString := string(dataJson)

		fmt.Println(dataString)
		byteArgs := [][]byte{
			[]byte(key),
			[]byte(accId),
			[]byte(dataString),
		}
		SubmitTransaction(client, byteArgs, fcn)
	}
}

func SubmitRecallData(client *clientChannel.Client, args string, objectType parser.ObjectType) {
	jsonFile, err := os.Open(args)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteArray, _ := ioutil.ReadAll(jsonFile)

	var recalls []parser.RecallData
	json.Unmarshal(byteArray, &recalls)

	fcn := "AddRECData"

	for _, data := range recalls {
		key := uuid.New().String()
		accId := fmt.Sprint(data.GetId())

		dataJson, _ := json.Marshal(data)
		dataString := string(dataJson)

		fmt.Println(dataString)
		byteArgs := [][]byte{
			[]byte(key),
			[]byte(accId),
			[]byte(dataString),
		}
		SubmitTransaction(client, byteArgs, fcn)
	}
}

func SubmitImportData(client *clientChannel.Client, args string, objectType parser.ObjectType) {
	jsonFile, err := os.Open(args)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteArray, _ := ioutil.ReadAll(jsonFile)

	var imports []parser.ImportData
	json.Unmarshal(byteArray, &imports)

	fcn := "AddRECData"

	for _, data := range imports {
		key := uuid.New().String()
		accId := fmt.Sprint(data.GetId())

		dataJson, _ := json.Marshal(data)
		dataString := string(dataJson)

		fmt.Println(dataString)
		byteArgs := [][]byte{
			[]byte(key),
			[]byte(accId),
			[]byte(dataString),
		}
		SubmitTransaction(client, byteArgs, fcn)
	}
}

func SubmitTransaction(client *clientChannel.Client, byteArgs [][]byte, fcn string) {
	response, err := client.Execute(clientChannel.Request{
		ChaincodeID: cfg.CVars.ChaincodeId,
		Fcn:         fcn,
		Args:        byteArgs,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("insert response: %v\n", string(response.Payload))

}

func GetOwnerCredit(client *clientChannel.Client, ownerId string) {
	if ownerId == "" {
		panic("ownerId cannot be empty")
	}
	args := [][]byte{
		[]byte(ownerId),
	}
	response, err := client.Query(clientChannel.Request{
		ChaincodeID: cfg.CVars.ChaincodeId,
		Fcn:         "ReadAsset",
		Args:        args,
	})
	if err != nil {
		panic(err)
	}
	// to pass somewhere?
	fmt.Printf("query response: %v\n", string(response.Payload))
}
