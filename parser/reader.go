package parser

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// CSV -> []Map[string]string -> Obj -> Json -> Byte(Json) -> Obj
//
//

type ObjectType string

const (
	EVENT_DATA   ObjectType = "EVENT_DATA"
	IMPORT_DATA  ObjectType = "IMPORT_DATA"
	RECALL_DATA  ObjectType = "RECALL_DATA"
	PRODUCT_DATA ObjectType = "PRODUCT_DATA"
)

func CSVToData(filePath string, objectType ObjectType) (dataset []ObjectData) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	records := ToMap(csvFile)

	for _, record := range records {
		switch objectType {
		case EVENT_DATA:
			var event EventData
			event.PopulateWithMap(record)
			dataset = append(dataset, &event)
		case IMPORT_DATA:
			var imp ImportData
			imp.PopulateWithMap(record)
			dataset = append(dataset, &imp)
		}
	}
	return dataset
}

func JSONToData(filepath string, objectType ObjectType) (dataset []ObjectData) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteArray, _ := ioutil.ReadAll(jsonFile)

	switch objectType {
	case IMPORT_DATA:
		imports := []ImportData{}
		json.Unmarshal(byteArray, &imports)
		for _, imp := range imports {
			i := imp
			dataset = append(dataset, &i)
		}
	case PRODUCT_DATA:
		products := []ProductData{}
		json.Unmarshal(byteArray, &products)
		for _, product := range products {
			p := product
			dataset = append(dataset, &p)
		}
	case RECALL_DATA:
		recalls := []RecallData{}
		json.Unmarshal(byteArray, &recalls)
		for _, recall := range recalls {
			r := recall
			dataset = append(dataset, &r)
		}
	}
	return dataset
}

func ToMap(file *os.File) (records []map[string]string) {
	reader := csv.NewReader(file)
	var header []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if header == nil {
			header = record
			continue
		}

		dict := map[string]string{}
		for i := range header {
			dict[header[i]] = record[i]
		}
		records = append(records, dict)
	}
	return records
}
