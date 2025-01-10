package btclightclient

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// tx bytes
// address

func (lc *BTCLightClient) GetBalance(tx *wire.MsgTx, addr *chainhash.Hash) (uint64, error) {
	for _, output := range tx.TxOut {
		scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(
			output.PkScript, lc.ChainParams())

		// TODO: Handle other script types: pay pubkey, pay script hash, multi signature...
		if scriptClass != txscript.PubKeyHashTy {
			return 0, errors.New("only support pubkey hash")
		}

		fmt.Println(addresses)
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
