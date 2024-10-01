package types

import (
	"encoding/json"

	"github.com/btcsuite/btcd/wire"
)

// lightblock data
type BTCLightBlock struct {
	Height uint64
	Header *wire.BlockHeader
}

func NewBTCLightBlock(height uint64, header *wire.BlockHeader) BTCLightBlock {
	return BTCLightBlock{
		Height: height,
		Header: header,
	}
}

// use for protobuf
func (btcLightBlock BTCLightBlock) Marshal() ([]byte, error) {
	return json.Marshal(btcLightBlock)
}

// use for protobuf
func (btcLightBlock *BTCLightBlock) Unmarshal(data []byte) error {
	return json.Unmarshal(data, btcLightBlock)
}

// use for protobuf
func (btcLightBlock *BTCLightBlock) Size() int {
	bz, _ := btcLightBlock.Marshal()
	return len(bz)
}

// use for prototbuf
func (btcLightBlock *BTCLightBlock) MarshalTo(data []byte) (int, error) {
	bz, err := btcLightBlock.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, bz)
	return len(data), nil
}
