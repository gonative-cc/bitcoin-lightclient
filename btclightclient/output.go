package btclightclient

import (
	"bytes"
	"encoding/hex"
	"errors"
	// "fmt"

	// "github.com/btcsuite/btcd/btcutil"
	// "github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func MsgTxFromHex(txData string) (wire.MsgTx, error) {
	var tx wire.MsgTx
	txBuf, err := hex.DecodeString(txData)
	if err != nil {
		return tx, err
	}

	rbuf := bytes.NewReader(txBuf)
	err = tx.Deserialize(rbuf)
	if err != nil {
		return tx, err
	}
	return tx, nil
}

// Get balance of addr from tx
func (lc *BTCLightClient) GetBalance(tx *wire.MsgTx, addr string) (int64, error) {
	balance := int64(0)
	for _, output := range tx.TxOut {
		_, addresses, _, err := txscript.ExtractPkScriptAddrs(
			output.PkScript, lc.ChainParams())
		if err != nil {
			return 0, err
		}

		// TODO: Handle other script types: pay pubkey, pay script hash, multi signature...
		if addresses[0].String() != addr {
			break
		}
		balance = balance + output.Value
	}

	return balance, nil
}

// Verify addr balance in this block. We check:
// - tx valid
// - return balance from tx
func (lc *BTCLightClient) VerifyBalance(tx *wire.MsgTx, addr string, spv SPVProof) (int64, error) {
	txID := tx.TxID()

	if spv.txId != txID {
		return 0, errors.New("TX ID not match")
	}

	if spvStatus := lc.VerifySPV(spv); spvStatus != ValidSPVProof {
		return 0, errors.New("spv not valid")
	}

	balance, err := lc.GetBalance(tx, addr)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
