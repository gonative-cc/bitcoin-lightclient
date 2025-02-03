package btclightclient

import (
	"crypto/sha256"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func GenerateBlock(
	prevHash chainhash.Hash,
	miner string,
	txns []*btcutil.Tx,
	difficulty int,
) (*wire.MsgBlock, error) {
	block := &wire.MsgBlock{}
	targetPrefix := strings.Repeat("0", difficulty)

	for nonce := 0; ; nonce++ {
		txnStrings := ""
		for _, tx := range txns {
			txnStrings += tx.Hash().String()
		}
		hashInput := prevHash.String() + txnStrings + miner + strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(hashInput))
		var blockHash chainhash.Hash
		copy(blockHash[:], hash[:])

		if strings.HasPrefix(blockHash.String(), targetPrefix) {
			block.Header.PrevBlock = prevHash
			block.Header.Nonce = uint32(nonce)
			block.Header.MerkleRoot = blockHash // Simulated Merkle Root
			return block, nil
		}
	}
}
