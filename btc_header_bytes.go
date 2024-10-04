package main

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
)

// type BTCHeaderBytes []byte
// const BTCHeaderSize = 80

func NewBlockHeader(data []byte) (*wire.BlockHeader, error) {
	var header wire.BlockHeader
	reader := bytes.NewReader(data)
	err := header.Deserialize(reader)

	return &header, err
}

