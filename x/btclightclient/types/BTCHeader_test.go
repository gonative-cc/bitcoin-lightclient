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

	blockHeaderAfter, err := NewBlockHeaderFromBytes(data)
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
	data, err := ByteFromBlockHeader(&blockHeader)
	if err == nil {
		fmt.Println(data)
	} else {
		t.Fatalf("error")
	}

}

func TestNewBlockHeaderFromByteFail(t *testing.T) {
	data := []byte(`abcd`)
	_, err := NewBlockHeaderFromBytes(data)
	if err == nil {
		t.Fatalf("Should failed")
	}
}


// func TestNewBTCLightBlock(t *testing.T) {
// 	blockHeader := sampleBlockHeader()
// 	lightBlock := NewBTCLightBlock(10, &blockHeader)


// }
