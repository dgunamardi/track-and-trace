package config

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/ghodss/yaml"
)

type CompositePath struct {
	Json string
	Yaml string
}

type Credentials struct {
	OrgId    string
	UserName string
}

type ConfigVariables struct {
	ConfigPath  CompositePath
	Credentials Credentials

	ChaincodeId string
	ChannelId   string
}

var (
	CVars = ConfigVariables{
		ConfigPath: CompositePath{
			Json: "/home/tkgoh/Sandbox/track-and-trace/ccp-bcs2/bcs2-channel-sdk-config.json",
			Yaml: "/home/tkgoh/Sandbox/track-and-trace/ccp-bcs2/bcs2-channel-sdk-config.yaml",
		},
		Credentials: Credentials{
			OrgId:    "e104613781c697e3e9ec6b02f6876d4d42604f93",
			UserName: "Admin",
		},
	}
	SdkFile *simplejson.Json
)

func LoadConfig() {
	data, err := ReadFile(CVars.ConfigPath.Yaml)
	if err != nil {
		panic(err)
	}
	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		panic(err)
	}
	SdkFile, err = simplejson.NewJson(data)
	CVars.ChannelId = GetDefaultChannelId()
	CVars.ChaincodeId = GetDefaultChaincodeId()
}

func GetDefaultChannelId() string {
	channels := SdkFile.Get("channels").MustMap()
	for k := range channels {
		return k
	}
	return ""
}

func GetDefaultChaincodeId() string {
	chaincodes := SdkFile.Get("channels").Get(CVars.ChannelId).Get("chaincodes").MustArray()
	if str, ok := chaincodes[0].(string); ok {
		return strings.Split(str, ":")[0]
	}
	return ""
}

// ReadFile reads the file named by filename and returns the contents.
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64 = bytes.MinRead

	if fi, err := f.Stat(); err == nil {
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return readAll(f, n)
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.

	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if errors, ok := e.(error); ok && errors == bytes.ErrTooLarge {
			err = errors
		} else {
			panic(e)
		}
	}()
	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
