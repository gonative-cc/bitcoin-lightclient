package btclightclient

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/wire"
)

// We don't want to make it compilicated yet :D
// type BTCHeaderBytes []byte

const BTCHeaderSize = 80

func BlockHeaderFromHex(hexStr string) (wire.BlockHeader, error) {
	var header wire.BlockHeader

	fmt.Println(len(hexStr))
	if len(hexStr) != BTCHeaderSize*2 {
		return header, errors.New("invalid header size, must have 80 bytes")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return header, err
	}

	reader := bytes.NewReader(data)
	err = header.Deserialize(reader)
	return header, err
}
