package btclightclient

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"gotest.tools/assert"
)

func TestGenerateBlock(t *testing.T) {
	tests := []struct {
		name            string
		transactionsRaw []string
		miner           string
		previousHash    string
		difficulty      int
		expectedError   error
	}{
		{
			name:            "Valid block with one transaction",
			transactionsRaw: []string{"0100000001d95f1a3f947ba8f60df95e47147dd918446d74604a2829fc7f02aa40c1f7a0c8010000006b483045022100f3abbe1b0d622cc80a72f5760d78a21b33c2c65a82b37a7a76cd3a47b38e597e02207eab42fd194fa869d1b9a76e6e7a76d289bd2391c5c5a15c83e5b42d339f96da01210311b8ce3e832cc4cd22749587cb8a5cb1e053bd8691bfc9f80297bd50990bfddcffffffff0280d1f008000000001976a9144620b10b5d2cdde246fc82a94a5827ea0ee5426188ac40420f00000000001976a914ad5b5d5f9c69f5b05cfe08769c2675c0f7446a4a88ac00000000"},
			miner:           "miner1",
			difficulty:      1,
			previousHash:    "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			expectedError:   nil,
		},
		{
			name: "Valid block with two transactions",
			transactionsRaw: []string{
				"0100000001d95f1a3f947ba8f60df95e47147dd918446d74604a2829fc7f02aa40c1f7a0c8010000006b483045022100f3abbe1b0d622cc80a72f5760d78a21b33c2c65a82b37a7a76cd3a47b38e597e02207eab42fd194fa869d1b9a76e6e7a76d289bd2391c5c5a15c83e5b42d339f96da01210311b8ce3e832cc4cd22749587cb8a5cb1e053bd8691bfc9f80297bd50990bfddcffffffff0280d1f008000000001976a9144620b10b5d2cdde246fc82a94a5827ea0ee5426188ac40420f00000000001976a914ad5b5d5f9c69f5b05cfe08769c2675c0f7446a4a88ac00000000",
				"0100000001a3b2c1d4e5f67890abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef000000006a47304402207f2e3d1c0b9a8f7e6d5c4b3a291827161514131211100f0e0d0c0b0a0908070602205f4e3d2c1b0a9f8e7d6c5b4a392817161514131211100f0e0d0c0b0a09080706012102abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdffffffff0280d1f008000000001976a914abcdefabcdefabcdefabcdefabcdefabcdefabcdef88ac10270000000000001976a91489abcdef0123456789abcdef0123456789abcdef012345678988ac00000000",
			},
			miner:         "miner2",
			difficulty:    1,
			previousHash:  "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			expectedError: nil,
		},
		{
			name:            "Empty transaction list",
			transactionsRaw: []string{},
			previousHash:    "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			miner:           "miner1",
			difficulty:      1,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prevHash, _ := chainhash.NewHashFromStr(tt.previousHash)

			transactions := make([]*btcutil.Tx, 0)
			for _, txnRaw := range tt.transactionsRaw {
				rawTxBytes, err := hex.DecodeString(txnRaw)
				assert.NilError(t, err)

				// Create an empty wire.MsgTx
				msgTx := wire.NewMsgTx(wire.TxVersion)

				// Deserialize the raw transaction into the MsgTx object
				err = msgTx.Deserialize(bytes.NewReader(rawTxBytes))
				assert.NilError(t, err)

				// Convert wire.MsgTx to btcutil.Tx
				txn := btcutil.NewTx(msgTx)
				transactions = append(transactions, txn)
			}

			block, err := GenerateBlock(*prevHash, tt.miner, transactions, tt.difficulty)

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

			if len(block.Transactions) != len(transactions) {
				t.Errorf("Expected %d transactions, got %d", len(transactions), len(block.Transactions))
			}

			if len(block.Header.BlockHash()) == 0 {
				t.Error("Block hash should not be empty")
			}
		})
	}
}
