package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	pmsp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	PrivateKey string
	User       pmsp.SigningIdentity
	Cert       []byte
)

func InitializeUserIdentity() {
	mspClient, err := msp.New(Sdk.Context(), msp.WithOrg(CVars.Credentials.OrgId))
	if err != nil {
		panic(errors.Errorf("Error creating MSP client: %s", err))
	}

	if PrivateKey == "" {
		User, err = mspClient.GetSigningIdentity(CVars.Credentials.UserName)
		if err != nil {
			panic(errors.Errorf("GetSigningIdentity returned error: %v", err))
		}

	} else {
		if Cert == nil {
			if Cert, err = GetSigncertsBytes(CVars.Credentials.OrgId); err != nil {
				panic(errors.Errorf("Nope: %v", err))
			}
		}

		// pvtKey must be decrypted when it is passed in function CreateSigningIdentity.
		User, err = mspClient.CreateSigningIdentity(pmsp.WithCert(Cert), pmsp.WithPrivateKey([]byte(PrivateKey)))
		if err != nil {
			panic(errors.Errorf("CreateSigningIdentity returned error: %v", err))
		}
	}

}

// GetCryptoPath get msp directory from sdk config file's path
func GetCryptoPath(ordId string) string {
	cryptoPath := SdkFile.Get("organizations").Get(ordId).Get("cryptoPath").MustString()
	return cryptoPath
}

// GetSigncertsBytes can get Signcerts from sdk config file's path
// cryptoPath is get from function GetCryptoPath
func GetSigncertsBytes(orgId string) ([]byte, error) {
	cryptoPath := GetCryptoPath(orgId)
	signcertsPathdir := filepath.Join(cryptoPath, "signcerts")
	files, _ := ioutil.ReadDir(signcertsPathdir)
	if len(files) != 1 {
		return nil, errors.Errorf("file count invalid in the directory [%s]", signcertsPathdir)
	}

	f, err := ioutil.ReadFile(filepath.Join(signcertsPathdir, files[0].Name()))
	if err != nil {
		return nil, errors.Errorf("read signcerts from [%s] fail", files[0].Name())
	} else if f == nil {
		return nil, errors.Errorf("result of read signcerts file [%s] is null", files[0].Name())
	}

	return f, nil
}

// GetTlsCryptoKeyPath get tlsCryptoKeyPath from sdk config file's path with orgId
func GetTlsCryptoKeyPath(orgId string) string {
	tlsCryptoKeyPath := SdkFile.Get("organizations").Get(orgId).Get("tlsCryptoKeyPath").MustString()
	return tlsCryptoKeyPath
}

// GetTlsCryptoKey can get tlsCryptoKey content from tlsCryptoKeyPath
// tlsCryptoKeyPath is get from function GetTlsCryptoKeyPath
func GetTlsCryptoKey(orgId string) (string, error) {
	tlsCryptoKeyPath := GetTlsCryptoKeyPath(orgId)
	f, err := ioutil.ReadFile(tlsCryptoKeyPath)
	if err != nil {
		return "", errors.Errorf("read tlsCryptoKey from [%s] fail", tlsCryptoKeyPath)
	} else if f == nil {
		return "", errors.Errorf("result of read tlsCryptoKey file [%s] is null", tlsCryptoKeyPath)
	}

	return string(f), nil
}

// GetPrivateKeyBytes can get privateKey from sdk config file's path
// cryptoPath is get from function GetCryptoPath
func GetPrivateKeyBytes(orgId string) ([]byte, error) {
	cryptoPath := GetCryptoPath(orgId)
	keystorePathdir := filepath.Join(cryptoPath, "keystore")
	files, _ := ioutil.ReadDir(keystorePathdir)
	if len(files) != 1 {
		return nil, errors.Errorf("file count invalid in the directory [%s]", keystorePathdir)
	}

	f, err := ioutil.ReadFile(filepath.Join(keystorePathdir, files[0].Name()))
	if err != nil {
		return nil, errors.Errorf("read signcerts from [%s] fail", files[0].Name())
	} else if f == nil {
		return nil, errors.Errorf("result of read keystore file [%s] is null", files[0].Name())
	}

	return f, nil
}

// SetPrivateKey update the privateKey used in ChannelClient
func SetPrivateKey(key string) {
	PrivateKey = key
}

// ClearPrivateKey set privateKey empty
func ClearPrivateKey() {
	PrivateKey = ""
}

// SetClientTlsKey update the tls key in fabric-sdk
func SetClientTlsKey(tlsKey string) {
	vc := viper.New()
	vc.SetConfigFile(CVars.ConfigPath.Yaml)
	err := vc.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to read configFile: %s", CVars.ConfigPath.Yaml))
	}

	orgID := vc.GetString("client.originalOrganization")
	if orgID == "" {
		orgID = vc.GetString("client.organization")
	}

	//SetTlsClientKey can be used to update the tlskey in fabric-sdk
	fab.SetTlsClientKey(orgID, tlsKey)
}

// ClearClientTlsKey reset the tls key in fabric-sdk with tls key file specified in config
func ClearClientTlsKey() {
	vc := viper.New()
	vc.SetConfigFile(CVars.ConfigPath.Yaml)
	err := vc.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to read configFile: %s", CVars.ConfigPath.Yaml))
	}

	orgID := vc.GetString("client.originalOrganization")
	if orgID == "" {
		orgID = vc.GetString("client.organization")
	}
	fab.ResetTlsClientKeyWithOrgID(orgID)
}
