package main

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
)

// type BTCHeaderBytes []byte
// const BTCHeaderSize = 80

func NewBlockHeader(dataStr string) (*wire.BlockHeader, error) {

	data, _ := hex.DecodeString(dataStr)
	var header wire.BlockHeader
	reader := bytes.NewReader(data)
	err := header.Deserialize(reader)

	return &header, err
}

