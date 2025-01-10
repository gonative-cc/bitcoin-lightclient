package btclightclient

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/wire"
)


func TestOutputExtraction(t *testing.T) {
	txData := "020000000001017e10a9aacb82d088bbb9acbadfbd0a544d328bec62aea1b4135140aa06f86af31600000017160014ab454d6a2cc66c550a0d5cedb1a80164b31f163dffffffff01792f00000000000017a91424f4377fa4f486495beae33b64f915266d7fd1ba87024830450221009feec844f6556cfee7f3a9afdc96b3bc87ad23b43e92b4955e1ad1ba287ec360022013eee664a2d04f917d057f4d279822b7d1d20f9e8a168f7934616f0ee8b5c00201210357c4b37f213d01b463f6e9f1c816b6ff8ff6400cd4d7fcd1cd5db139b697bc0e00000000"
	txBuf, _ := hex.DecodeString(txData)
	rbuf := bytes.NewReader(txBuf)

	var tx wire.MsgTx
	tx.Deserialize(rbuf)
	fmt.Println(tx)
}
