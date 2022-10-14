package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"

	"earhart.com/parser"
)

func main() {
	path := "..\\data-samples\\recall.json"

	fileExt := filepath.Ext(path)
	fmt.Println(fileExt)

	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteArray, _ := ioutil.ReadAll(jsonFile)
	var products []parser.RecallData

	json.Unmarshal(byteArray, &products)

	fmt.Println(products[0])
}
