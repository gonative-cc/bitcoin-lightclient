package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func sampleBlockHeader() wire.BlockHeader {
	prevHash, _ := chainhash.NewHashFromStr(fmt.Sprintf("%x", 10000))
	merketRoot, _ := chainhash.NewHashFromStr(fmt.Sprintf("%x", 10000))

	timestamp := time.Unix(1727624087, 0)

	blockHeader := wire.BlockHeader{
		Version:    1,
		PrevBlock:  *prevHash,
		MerkleRoot: *merketRoot,
		Timestamp:  timestamp,
		Bits:       20240924,
		Nonce:      20000,
	}
	return blockHeader
}

func TestNewBlockHeaderFromByte(t *testing.T) {
	blockHeaderBefore := sampleBlockHeader()
	data, _ := ByteFromBlockHeader(&blockHeaderBefore)
	blockHeaderAfter, err := NewBlockHeader(data)

	if err != nil {
		t.Fatalf("error")
	} else {
		if *blockHeaderAfter != blockHeaderBefore {
			t.Fatalf("not equal")
		}
	}
}

func TestByteFromBlockHeader(t *testing.T) {
	blockHeader := sampleBlockHeader()
	_, err := ByteFromBlockHeader(&blockHeader)
	if err != nil {
		t.Fatalf("error")
	}

}

func TestBlockHeaderHex(t *testing.T) {
	headerHex := `00809c2cc58dd5bb09f12796b3c3a7d69b4901f857d274229ba00000000000000000000093ac838cd2308ba70827bc48bc0a8f62a4156a6e5cf9903ac478e1e44f035153971db164943805177722d461`
	var btcHeaderBytes BTCHeaderBytes
	btcHeaderBytes.UnmarshalHex(headerHex)
	fmt.Println(btcHeaderBytes.NewBlockHeaderFromBytes())
	
}

func TestNewBlockHeaderFromByteFail(t *testing.T) {
	data := []byte(`abcd`)
	_, err := NewBlockHeader(data)
	if err == nil {
		t.Fatalf("Should failed")
	}
}
