package parser

import (
	//"fmt"

	"regexp"
	"strconv"
	"strings"
)

type EventData struct {
	EventId      string   `bson:"event_id" json:"event_id"`
	EventType    int32    `bson:"event_type" json:"event_type"`
	InputGTIN    []string `bson:"input_gtin" json:"input_gtin"`
	OutputGTIN   []string `bson:"output_gtin" json:"output_gtin"`
	SerialNumber string   `bson:"serial_number" json:"serial_number"`
	EventTime    string   `bson:"event_time" json:"event_time"`
	EventLoc     string   `bson:"event_loc" json:"event_loc"`
	LocationName string   `bson:"location_name" json:"location_name"`
	CompanyName  string   `bson:"company_name" json:"company_name"`
}

// NOTES
// IF tags doesnt work in fabric, just all string, parse later

func (txData *EventData) Populate(rawByte []byte) {
	stringArr := CleanStringArr(string(rawByte))
	if len(stringArr) < 9 {
		return
	}
	//fmt.Println(stringArr)

	txData.EventId = stringArr[0]

	eventType, err := strconv.Atoi(stringArr[1])
	if err != nil {
		panic(err)
	}
	txData.EventType = int32(eventType)

	txData.InputGTIN = CleanGTIN(stringArr[2])
	txData.OutputGTIN = CleanGTIN(stringArr[3])

	txData.SerialNumber = stringArr[4]
	txData.EventTime = stringArr[5]
	txData.EventLoc = stringArr[6]
	txData.LocationName = stringArr[7]
	txData.CompanyName = stringArr[8]

}

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
