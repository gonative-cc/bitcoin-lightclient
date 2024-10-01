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
	LatestBlockKey = []byte("lastBlock")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}



func HeaderKey(height uint64) ([]byte, error) {
	buf := new(bytes.Buffer)

	if  err := binary.Write(buf, binary.BigEndian, height); err != nil {
		return nil, err
	}

	key := make([]byte, 1)
	key[0] = 1;
	
	value := append(key, buf.Bytes()...)
	return value, nil
}
