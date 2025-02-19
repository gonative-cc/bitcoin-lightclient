package btclightclient

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"gotest.tools/assert"
)

func TestGenerateBlock(t *testing.T) {
	tests := []struct {
		name          string
		transactions  []*btcutil.Tx
		miner         string
		previousHash  string
		difficulty    int
		expectedError error
	}{
		{
			name:          "Empty transaction list",
			transactions:  []*btcutil.Tx{},
			previousHash:  "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			miner:         "miner1",
			difficulty:    1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prevHash, _ := chainhash.NewHashFromStr(tt.previousHash)

			block, err := GenerateBlock(*prevHash, tt.miner, tt.transactions, tt.difficulty)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			}

			assert.NilError(t, err)

			if block.Header.PrevBlock.String() != tt.previousHash {
				t.Errorf("Expected previous hash %s, got %s", tt.previousHash, block.Header.PrevBlock.String())
			}

			if len(block.Transactions) != len(tt.transactions) {
				t.Errorf("Expected %d transactions, got %d", len(tt.transactions), len(block.Transactions))
			}

			if len(block.Header.BlockHash()) == 0 {
				t.Error("Block hash should not be empty")
			}
		})
	}
}
