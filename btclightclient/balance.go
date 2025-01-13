package btclightclient

import (
	"bytes"
	"encoding/hex"
	"errors"

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
func (lc *BTCLightClient) GetBalance(tx wire.MsgTx, addr string) (int64, error) {
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

type BalanceStatus int

const (
	// return when user provide invalid proof of utxo
	InvalidBalance BalanceStatus = iota
	// return when block are waiting for block finalize 
	WaitingForConfirmBalance
	// return when balance valid (block finalize)
	ValidBalance
)

type BalanceReport struct {
	balance int64
	status  BalanceStatus
}

func newInvalidBalance() BalanceReport {
	return BalanceReport{
		balance: 0,
		status:  InvalidBalance,
	}
}

func newValidBalance(balance int64) BalanceReport {
	return BalanceReport{
		balance: balance,
		status:  ValidBalance,
	}
}

func newWaitingConfirmBalance(balance int64) BalanceReport {
	return BalanceReport{
		balance: balance,
		status:  WaitingForConfirmBalance,
	}
}

// Verify addr balance in this block. We check:
// - tx valid
// - return balance from tx
func (lc *BTCLightClient) VerifyBalance(tx wire.MsgTx, addr string, spv SPVProof) (BalanceReport, error) {
	txID := tx.TxID()

	if spv.TxId != txID {
		return newInvalidBalance(), errors.New("TX ID not match")
	}

	spvStatus := lc.VerifySPV(spv)
	if spvStatus == InvalidSPVProof {
		return newInvalidBalance(), errors.New("spv invalid")
	}
	
	balance, err := lc.GetBalance(tx, addr);
	if err != nil {
		return newInvalidBalance(), err
	}
	

	if spvStatus == PartialValidSPVProof {
		return newWaitingConfirmBalance(balance), nil
	}

	
	return newValidBalance(balance), nil
}
