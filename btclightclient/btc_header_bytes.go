package btclightclient

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
)

const BTCHeaderSize = 80 // 80 bytes

// Utils for converting hex string to header
// The input must be 80 bytes hex string type
func BlockHeaderFromHex(hexStr string) (wire.BlockHeader, error) {
	var header wire.BlockHeader

	if len(hexStr) != BTCHeaderSize*2 {
		return header, ErrInvalidHeaderSize
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return header, err
	}

	reader := bytes.NewReader(data)
	err = header.Deserialize(reader)
	return header, err
}
