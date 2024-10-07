package main

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
)

// We don't want to make it compilicated yet :D
// type BTCHeaderBytes []byte
// const BTCHeaderSize = 80

func BlockHeaderFromHex(hexStr string) (*wire.BlockHeader, error) {
	data, _ := hex.DecodeString(hexStr)
	var header wire.BlockHeader
	reader := bytes.NewReader(data)
	err := header.Deserialize(reader)
	return &header, err
}
