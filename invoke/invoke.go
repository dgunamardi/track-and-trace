package main

import (
	"fmt"
	"os"

	cfg "earhart.com/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	clientChannel "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
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
	case "insertData":
		SubmitTranscation(client)
	case "getOwnerCredit":
		GetOwnerCredit(client, args[1])

	default:
		panic("argument is not available. Available Arguments:\n- insertData\n- getOwnerCredit")
	}
}

func SubmitTranscation(client *clientChannel.Client) {
	stringArgs := []string{
		"912edf2e-933d-4793-9ba0-2077c57070aq",
		"912edf2e-933d-4793-9ba0-2077c57070aq",
		"1043022868954",
		"42547cba-a4aa-4758-812f-a3699489c1c4",
		"1",
		"69700806964203",
		"210UDXYXNFEKJAIWFVQMW",
		"2022-Aug-09T15:21:40 +0000",
		"-8.1971482,114.4440049",
		"Bali",
		"Bali_factory_2",
	}

	var byteArgs [][]byte

	for _, stringArg := range stringArgs {
		byteArgs = append(byteArgs, []byte(stringArg))
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
	fmt.Printf("query response: %v\n", string(response.Payload))
}
