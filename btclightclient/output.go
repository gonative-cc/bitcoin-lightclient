package btclightclient

import (
	"errors"
	"fmt"
	// "fmt"

	// "github.com/btcsuite/btcd/btcutil"
	// "github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// Get balance of addr from tx
func (lc *BTCLightClient) GetBalance(tx *wire.MsgTx, addr string) (int64, error) {
	balance := int64(0)
	for _, output := range tx.TxOut {
		scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(
			output.PkScript, lc.ChainParams())

		if err != nil {
			return 0, err
		}

		// TODO: Handle other script types: pay pubkey, pay script hash, multi signature...
		if scriptClass != txscript.PubKeyHashTy {
			return 0, errors.New("only support pubkey hash")
		}

		fmt.Println(addresses[0].String())
		if addresses[0].String() != addr {
			break
		}
		balance = balance + output.Value
	}

	return balance, nil
}
