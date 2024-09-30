package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/btcsuite/btcd/wire"
)

/**
 * This type is input of insert Header 
 */
type BTCHeaderBytes []byte

const BTCHeaderLen = 80

func (m BTCHeaderBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MarshalHex())
}

func (m *BTCHeaderBytes) UnmarshalJSON(bz []byte) error {
	var headerHexStr string

	if err := json.Unmarshal(bz, &headerHexStr); err != nil {
		return err
	}

	return m.UnmarshalHex(headerHexStr)
}

func (m BTCHeaderBytes) Marshal() ([]byte, error) {
	return m, nil
}

func (m *BTCHeaderBytes) Unmarshal(data []byte) error {
	if len(data) != BTCHeaderLen {
		return errors.New("invalid header length(80)")
	}

	// Verify data <=> wire.BlockHeader
	if _, err := NewBlockHeader(data); err != nil {
		return errors.New("bytes do not correspond to a *wire.BlockHeader object")
	}

	*m = data
	return nil
}

func (m BTCHeaderBytes) MarshalHex() string {
	// TODO: what happend when this have error
	btcdHeader, _ := m.NewBlockHeaderFromBytes()

	var buf bytes.Buffer

	if err := btcdHeader.Serialize(&buf); err != nil {
		panic("Block header object cannot be converted to hex")
	}
	return hex.EncodeToString(buf.Bytes())
}

func (m *BTCHeaderBytes) UnmarshalHex(header string) error {
	// Decode the hash string from hex
	decoded, err := hex.DecodeString(header)

	if err != nil {
		return err
	}

	return m.Unmarshal(decoded)
}

func (m BTCHeaderBytes) MarshalTo(data []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, bz)
	return len(data), nil
}

func (m *BTCHeaderBytes) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

// return bytes
func (headerBytes BTCHeaderBytes) NewBlockHeaderFromBytes() (*wire.BlockHeader, error) {
	return NewBlockHeader(headerBytes)
}

// creates a block header from bytes.
func NewBlockHeader(data []byte) (*wire.BlockHeader, error) {
	var header wire.BlockHeader
	reader := bytes.NewReader(data)
	err := header.Deserialize(reader)

	return &header, err
}

// Return Blockheader in bytes format
func ByteFromBlockHeader(header *wire.BlockHeader) ([]byte, error) {
	var data bytes.Buffer

	err := header.Serialize(&data)

	return data.Bytes(), err
}
