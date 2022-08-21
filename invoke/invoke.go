package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	cfg "earhart.com/config"

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
	case "-s":
		SubmitTransaction(client, args[1:])
	case "-q":
		GetOwnerCredit(client, args[1])
	case "-sf":
		SubmitTransactionFromFile(client, args[1])
	default:
		panic("argument is not available. Available Arguments:\n-s for 'submitTransaction'\n-q for 'getOwnerCredit'\n-sf for 'submitTransactionFromFile'\n")
	}
}

func SubmitTransaction(client *clientChannel.Client, args []string) {
	var byteArgs [][]byte

	for _, arg := range args {
		byteArgs = append(byteArgs, []byte(arg))
	}

	response, err := client.Execute(clientChannel.Request{
		ChaincodeID: cfg.CVars.ChaincodeId,
		Fcn:         "AddCTEwithAsset",
		Args:        byteArgs,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("insert response: %v\n", string(response.Payload))

}

// 9 indices data format
func SubmitTransactionFromFile(client *clientChannel.Client, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("unable to read file from path: %v", err))
	}
	defer f.Close()

	dmaps := CSVToMap(f)
	for _, dmap := range dmaps {
		//make random key
		previous_key := uuid.New().String()
		new_key := previous_key

		// get data
		event_id := fmt.Sprint(dmap["event_id"])
		event_type := fmt.Sprint(dmap["event_type"])
		event_time := fmt.Sprint(dmap["event_time"])
		generator_gln := fmt.Sprint(dmap["generator_gln"])
		serial_number := fmt.Sprint(dmap["serial_number"])
		event_location := fmt.Sprint(dmap["event_location"])
		location_name := fmt.Sprint(dmap["location_name"])
		company_name := fmt.Sprint(dmap["company_name"])

		// gtin check
		var input_gtin, output_gtin string
		if _, ok := dmap["gtin"]; ok {
			input_gtin = fmt.Sprint(dmap["gtin"])
			output_gtin = input_gtin
		} else {
			input_gtin = fmt.Sprint(dmap["input_gtin"])
			output_gtin = fmt.Sprint(dmap["output_gtin"])
		}

		args := []string{
			previous_key,
			new_key,
			generator_gln,
			event_id,
			event_type,
			input_gtin,
			output_gtin,
			serial_number,
			event_time,
			event_location,
			location_name,
			company_name,
		}
		fmt.Println(args)

		SubmitTransaction(client, args)
	}
}

func CSVToMap(file io.Reader) []map[string]string {
	csvReader := csv.NewReader(file)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
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
