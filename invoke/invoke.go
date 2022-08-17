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
	case "submitTransaction":
		SubmitTransaction(client, args[1:])
	case "getOwnerCredit":
		GetOwnerCredit(client, args[1])

	default:
		panic("argument is not available. Available Arguments:\n- insertData\n- getOwnerCredit")
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
