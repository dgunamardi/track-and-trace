package main

import (
	"regexp"
	"strconv"
	"strings"
)

type TransactionData struct {
	EventId      string   `bson:"event_id"`
	EventType    int32    `bson:"event_type"`
	InputGTIN    []string `bson:"input_gtin"`
	OutputGTIN   []string `bson:"output_gtin"`
	SerialNumber string   `bson:"serial_number"`
	EventTime    string   `bson:"event_time"`
	EventLoc     string   `bson:"event_loc"`
	LocationName string   `bson:"location_name"`
	CompanyName  string   `bson:"company_name"`
}

func (txData *TransactionData) Populate(rawByte []byte) {
	stringArr := CleanStringArr(string(rawByte))
	if len(stringArr) < 9 {
		return
	}

	// Event ID
	txData.EventId = stringArr[0]

	// Event Type
	eventType, err := strconv.Atoi(stringArr[1])
	if err != nil {
		panic(err)
	}
	txData.EventType = int32(eventType)

	// GTINS
	txData.InputGTIN = CleanGTIN(stringArr[2])
	txData.OutputGTIN = CleanGTIN(stringArr[3])

	// Serial Number
	txData.SerialNumber = stringArr[4]

	txData.EventTime = stringArr[5]
	txData.EventLoc = stringArr[6]
	txData.LocationName = stringArr[7]
	txData.CompanyName = stringArr[8]

}

func CleanStringArr(rawString string) (res []string) {
	//fmt.Println("raw string= ", rawString)
	trimBrackets := strings.Trim(rawString, "{}")

	// Between GTIN
	mGTIN := regexp.MustCompile(", ")
	replacedCommaGTIN := mGTIN.ReplaceAllString(trimBrackets, "|")

	// Between Coordinate
	mCoord := regexp.MustCompile("[[:alnum:]],[[:alnum:]]")
	replacedCommaCoord := mCoord.ReplaceAllString(replacedCommaGTIN, ";")

	// Remove All Quotes
	cleanQuotes := strings.ReplaceAll(replacedCommaCoord, "\"", "")
	//fmt.Println(cleanQuotes)

	stringRawArr := strings.Split(cleanQuotes, ",")

	for _, stringRaw := range stringRawArr {
		m := regexp.MustCompile("^(.*):")
		clean := m.ReplaceAllString(stringRaw, "")

		res = append(res, clean)
	}
	return res
}

func CleanGTIN(rawString string) (res []string) {
	trimBrackets := strings.Trim(rawString, "[]")

	stringRawArr := strings.Split(trimBrackets, "|")
	for _, stringRaw := range stringRawArr {
		res = append(res, strings.Trim(stringRaw, "'"))
	}
	//fmt.Println(res)

	return res
}
