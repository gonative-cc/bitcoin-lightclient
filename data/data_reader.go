package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
)

type Sample struct {
	Network      string   `json:"network"`
	StartHeight  int64    `json:"start_height"`
	BlockHeaders []string `json:"blockheaders"`
}

// map network name string to chaincfg param object
var NetworkMap = map[string]*chaincfg.Params{
	"mainnet":       &chaincfg.MainNetParams,
	"testnet3":      &chaincfg.TestNet3Params,
	"simnet":        &chaincfg.SimNetParams,
	"signet":        &chaincfg.SigNetParams,
	"regressionnet": &chaincfg.RegressionNetParams,
}

func ReadJSON(jsonFilePath string) (*chaincfg.Params, int64, []string, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, 0, nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, 0, nil, err
	}

	var dataContent Sample
	if err := json.Unmarshal(byteValue, &dataContent); err != nil {
		return nil, 0, nil, err
	}

	networkParams, ok := NetworkMap[dataContent.Network]
	if !ok {
		err := fmt.Errorf("network %s not found", dataContent.Network)
		return nil, 0, nil, err
	}

	return networkParams, dataContent.StartHeight, dataContent.BlockHeaders, nil
}
