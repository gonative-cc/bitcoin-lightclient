package types

import (

	"bytes"

	"github.com/btcsuite/btcd/wire"
)

type BTCLightBlock struct {
	height uint64
	header *wire.BlockHeader
}


// covert bytes to BlockHeader
func NewBlockHeaderFromBytes(data []byte) (*wire.BlockHeader, error) {
	header := &wire.BlockHeader{}

	reader := bytes.NewReader(data)

	err := header.Deserialize(reader)

	return header, err
}

func ByteFromBlockHeader(header *wire.BlockHeader) ([]byte, error){
	var data bytes.Buffer
	
	err := header.Serialize(&data)

	return data.Bytes(), err
}
