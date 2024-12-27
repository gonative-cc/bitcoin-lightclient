package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Sample struct {
	StartHeight   int64    `json:"start_height"`
	HeaderStrings []string `json:"header_strings"`
}

func ReadJson(jsonFilePath string) (int64, []string, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return 0, nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return 0, nil, err
	}

	var dataContent DataContent
	if err := json.Unmarshal(byteValue, &dataContent); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
		return 0, nil, err
	}

	return dataContent.StartHeight, dataContent.HeaderStrings, nil
}
