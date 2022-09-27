package parser

import (
	//"fmt"

	"regexp"
	"strconv"
	"strings"
)

type EventData struct {
	EventId       string   `bson:"event_id" json:"event_id"`
	EventType     int32    `bson:"event_type" json:"event_type"`
	InputGTIN     []string `bson:"input_gtin" json:"input_gtin"`
	OutputGTIN    []string `bson:"output_gtin" json:"output_gtin"`
	InputName     []string `bson:"input_name" json:"input_name"`
	OutputName    []string `bson:"output_name" json:"output_name"`
	SerialNumber  string   `bson:"serial_number" json:"serial_number"`
	EventTime     string   `bson:"event_time" json:"event_time"`
	EventLocation string   `bson:"event_location" json:"event_location"`
	GLN           string   `bson:"gln" json:"gln"`
	CompanyName   string   `bson:"company_name" json:"company_name"`
}

// for invoke
func (txData *EventData) PopulateWithMap(record map[string]string) {
	txData.EventId = record["event_id"]

	eventType, _ := strconv.Atoi(record["event_type"])
	txData.EventType = int32(eventType)

	var (
		inputGTIN, outputGTIN string
		inputName, outputName string
	)

	if _, ok := record["gtin"]; ok {
		inputGTIN = record["gtin"]
		outputGTIN = record["gtin"]
	} else {
		inputGTIN = record["input_gtin"]
		outputGTIN = record["output_gtin"]
	}

	if _, ok := record["name"]; ok {
		inputName = record["input_name"]
		outputName = record["output_name"]
	} else {
		inputName = record["input_name"]
		outputName = record["output_name"]
	}

	txData.InputGTIN = CleanGTIN(inputGTIN)
	txData.OutputGTIN = CleanGTIN(outputGTIN)

	txData.InputName = CleanGTIN(inputName)
	txData.OutputName = CleanGTIN(outputName)

	txData.SerialNumber = record["serial_number"]
	txData.EventTime = record["event_time"]
	txData.EventLocation = record["event_location"]
	txData.GLN = record["gln"]
	txData.CompanyName = record["company_name"]

}

// NO LONGER USED BUT KEPT FOR UTILITY
func CleanStringArr(rawString string) (res []string) {
	//fmt.Println("raw string:", rawString)
	cleanString := strings.Trim(rawString, "{}")

	// GLOBAL COLON
	mColon := regexp.MustCompile("(\"):(\"|\\d)")
	cleanString = mColon.ReplaceAllString(cleanString, "$1|$2")
	//fmt.Println("ColonToPipe:", cleanString)
	// ------------

	// GLOBAL COMMA
	mComma := regexp.MustCompile("(\"|\\d),(\")")
	cleanString = mComma.ReplaceAllString(cleanString, "$1;$2")
	//fmt.Println("CommaToSemiColon:", cleanString)
	// ------------

	stringRawArr := strings.Split(cleanString, ";")

	for _, stringRaw := range stringRawArr {
		//fmt.Println("String", stringRaw)
		// CLEAR KEY
		mKey := regexp.MustCompile("^(.+)\\|")
		clean := mKey.ReplaceAllString(stringRaw, "")

		// CLEAR QUOTES
		mQuote := regexp.MustCompile("\"(.+)\"")
		clean = mQuote.ReplaceAllString(clean, "$1")

		res = append(res, clean)
	}

	return res
}

func CleanGTIN(rawString string) (res []string) {
	cleanString := strings.Trim(rawString, "[]")
	//fmt.Println(cleanString)

	stringRawArr := strings.Split(cleanString, ", ")
	//fmt.Println(stringRawArr)

	for _, stringRaw := range stringRawArr {
		res = append(res, strings.Trim(stringRaw, "'"))
	}
	//fmt.Println(res)

	return res
}

func (txData *EventData) IsValid() bool {
	if txData.EventId != "" {
		return true
	}
	return false
}

func (txData *EventData) GetId() string {
	return txData.GLN
}
