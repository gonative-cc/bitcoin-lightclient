package types

import (
	"fmt"
	"testing"
	"time"
	
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func TestNewBlockHeaderFromByte(t *testing.T) {

}

func TestByteFromBlockHeader(t *testing.T) {
  	prevHash, _ := chainhash.NewHashFromStr(fmt.Sprintf("%x", 10000))
	merketRoot, _ := chainhash.NewHashFromStr(fmt.Sprintf("%x", 10000))

	timestamp := time.Unix(1727624087, 0)
	
	blockHeader := wire.BlockHeader{
		Version:   1,
		PrevBlock: *prevHash,
		MerkleRoot: *merketRoot,
		Timestamp: timestamp,
		Bits: 20240924,
		Nonce: 20000,
	}

	data, err := ByteFromBlockHeader(&blockHeader)
	if err == nil {
		fmt.Println(data)		
	} else {
		t.Fatalf("error")
	}

}
