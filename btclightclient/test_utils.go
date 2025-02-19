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
	transactions []*btcutil.Tx,
	difficulty int,
) (*wire.MsgBlock, error) {
	block := &wire.MsgBlock{}
	targetPrefix := strings.Repeat("0", difficulty)

	for nonce := 0; ; nonce++ {
		txnStrings := ""
		var txHashes []*chainhash.Hash
		for _, tx := range transactions {
			txHash := tx.Hash()
			txnStrings += txHash.String()
			txHashes = append(txHashes, txHash)
		}
		hashInput := prevHash.String() + txnStrings + miner + strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(hashInput))
		var blockHash chainhash.Hash
		copy(blockHash[:], hash[:])

		if strings.HasPrefix(blockHash.String(), targetPrefix) {
			block.Header.PrevBlock = prevHash
			block.Header.Nonce = uint32(nonce)
			merkleRoot := CalculateMerkleRoot(txHashes)
			block.Header.MerkleRoot = *merkleRoot[len(merkleRoot)-1]
			for _, tx := range transactions {
				block.Transactions = append(block.Transactions, tx.MsgTx())
			}
			return block, nil
		}
	}
}

// CalculateMerkleRoot creates a merkle tree from a slice of transaction hashes.
func CalculateMerkleRoot(transactions []*chainhash.Hash) []*chainhash.Hash {
	if len(transactions) == 0 {
		return []*chainhash.Hash{new(chainhash.Hash)}
	}
	if len(transactions) == 1 {
		return []*chainhash.Hash{transactions[0]}
	}

	// Calculate how many entries are required to hold the binary merkle tree
	nextPoT := nextPowerOfTwo(len(transactions))
	arraySize := nextPoT*2 - 1
	merkles := make([]*chainhash.Hash, arraySize)

	// Create the base transaction hashes and fill in any empty ones with the
	// hash of an empty string
	copy(merkles, transactions)
	for i := len(transactions); i < nextPoT; i++ {
		merkles[i] = transactions[len(transactions)-1]
	}

	// Start the array offset after the last transaction and adjusted to the
	// next power of two.
	offset := nextPoT
	for i := 0; i < arraySize-1; i += 2 {
		switch {
		case merkles[i] == nil:
			merkles[offset] = nil

		case merkles[i+1] == nil:
			newHash := HashMerkleBranches(merkles[i], merkles[i])
			merkles[offset] = newHash

		default:
			newHash := HashMerkleBranches(merkles[i], merkles[i+1])
			merkles[offset] = newHash
		}
		offset++
	}

	return merkles
}

// nextPowerOfTwo returns the next highest power of two from a given number if
// it is not already a power of two.
func nextPowerOfTwo(n int) int {
	// Return the number if it's already a power of 2.
	if n&(n-1) == 0 {
		return n
	}

	// Figure out and return the next power of two.
	exponent := uint(0)
	for n != 0 {
		n >>= 1
		exponent++
	}
	return 1 << exponent
}

// HashMerkleBranches takes two merkle tree branches as hash.Hash and returns the hash of their concatenation
func HashMerkleBranches(left, right *chainhash.Hash) *chainhash.Hash {
	// Concatenate the left and right nodes.
	var hash [chainhash.HashSize * 2]byte
	copy(hash[:chainhash.HashSize], left[:])
	copy(hash[chainhash.HashSize:], right[:])

	newHash := chainhash.DoubleHashH(hash[:])
	return &newHash
}
