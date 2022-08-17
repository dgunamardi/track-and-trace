package config

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
)

var (
	Sdk *fabsdk.FabricSDK
)

func InitializeSDK() {
	configProvider := config.FromFile(CVars.ConfigPath.Yaml)
	var opts []fabsdk.Option

	opts, err := getOptstoInitalizeSDK(CVars.ConfigPath.Yaml)
	if err != nil {
		panic(fmt.Errorf("Failed to create new SDK: %s\n", err))
	}

	Sdk, err = fabsdk.New(configProvider, opts...)
	if err != nil {
		panic(fmt.Errorf("Failed to create new SDK: %s\n", err))
	}
	//fmt.Println("fabric SDK initialized")

}

func getOptstoInitalizeSDK(configPath string) ([]fabsdk.Option, error) {
	var opts []fabsdk.Option

	vc := viper.New()
	vc.SetConfigFile(configPath)
	err := vc.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to create new SDK: %s\n", err))
	}

	org := vc.GetString("client.originalOrganization")
	if org == "" {
		org = vc.GetString("client.organization")
	}

	opts = append(opts, fabsdk.WithOrgid(org))
	opts = append(opts, fabsdk.WithUserName(CVars.Credentials.UserName))
	return opts, nil

}
