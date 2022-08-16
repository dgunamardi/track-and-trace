package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cfg "earhart.com/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"

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

func populateWallet(wallet *gateway.Wallet) error {
	credPath := "/home/tkgoh/Sandbox/track-and-trace/ccp/4f08db41ded98093a7266580a4a2ae3ce62ce74a.peer/msp"

	//certName := "cert.pem"
	certName := "Admin@4f08db41ded98093a7266580a4a2ae3ce62ce74a.peer-4f08db41ded98093a7266580a4a2ae3ce62ce74a.default.svc.cluster.local-cert.pem"
	certPath := filepath.Join(credPath, "signcerts", certName)
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("4f08db41ded98093a7266580a4a2ae3ce62ce74aMSP", string(cert), string(key))

	err = wallet.Put("admin", identity)
	if err != nil {
		return err
	}
	// fmt.Println("wallet done")
	return nil
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
	case "getOwnerCredit":
		GetOwnerCredit(client, args[1])

	default:
		panic("argument is not available. Available Arguments:\n- insertData\n- getOwnerCredit")
	}
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
