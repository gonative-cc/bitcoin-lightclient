package data

import (
	"encoding/json"
	"io"
	"os"
)

type Sample struct {
	StartHeight   int64    `json:"start_height"`
	HeaderStrings []string `json:"header_strings"`
}

func ReadJson(jsonFilePath string) (int64, []string, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return 0, nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return 0, nil, err
	}

	var dataContent Sample
	if err := json.Unmarshal(byteValue, &dataContent); err != nil {
		return 0, nil, err
	}

	return dataContent.StartHeight, dataContent.HeaderStrings, nil
}
