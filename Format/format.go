package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"fmt"

	"io/ioutil"

	"earhart.com/parser"
)

func main() {
	path := "..\\data-samples\\product.json"
	//pathcsv := "..\\data-samples\\generate_import_data.csv"

	fileExt := filepath.Ext(path)
	fmt.Println(fileExt)

	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteArray, _ := ioutil.ReadAll(jsonFile)

	var recalls []parser.RecallData
	json.Unmarshal(byteArray, &recalls)

	//fmt.Println(recalls)

	//dataset := parser.JSONToData(path, parser.RECALL_DATA)
	//fmt.Println(dataset)

	//datasetcsv := parser.CSVToData(pathcsv, parser.IMPORT_DATA)
	//fmt.Println(datasetcsv)

	dset := parser.JSONToData(path, parser.PRODUCT_DATA)
	//p := &parser.ProductData{}

	fmt.Println(dset[0])
}
