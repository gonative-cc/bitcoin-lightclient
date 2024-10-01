package types

import (
	"bytes"
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "btclightclient"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_btclightclient"
)

var (
	ParamsKey = []byte("p_btclightclient")
	LatestBlockKey = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}



func HeaderKey(height uint64) ([]byte, error) {
	buf := new(bytes.Buffer)

	if  err := binary.Write(buf, binary.BigEndian, height); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
