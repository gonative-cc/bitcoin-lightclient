package btclightclient

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)


func TestOutputExtraction(t *testing.T) {
	// txID = f82a618b2f6212f2134d2dc66aad26c0aa0a2ed7fa53a298c53736d0510ad69f
	txData := "01000000014260bfba493220006a8141b162e5785c863149c043d7abf0030276c10d2a72e5010000008b483045022100d1424374bd4e6a0264dd55baf7de69a17bbe2416b20a55a1475fd3779799d13b02202e9249feb3b43cd82c8f097a1ffb5f98b06ca27595697bbc2acee4a5c255fd830141045408a52d4b3cdc9c78c14418c38d4a0fd6e0ed396f966d3f51c86164351bbd0e38426b922c9c1f79bca05dfeb72e08ba9d41eb5972e48db217301a9bf795aeadffffffff020065cd1d000000001976a914de526c004ab53c6d7ed3ff98789bca013797d33f88acf0ca8744000000001976a9148fd0e06ce84a4f108b1e259f5aeef21a16a2272c88ac00000000"
	txBuf, _ := hex.DecodeString(txData)
	rbuf := bytes.NewReader(txBuf)

	var tx wire.MsgTx
	tx.Deserialize(rbuf)

	
	pk, netID, err := base58.CheckDecode("1MGXkMpTTA4Ue4wSDh4kGBKfLSwj93MHTq");
	fmt.Println(pk, netID, err)
	addr, _ := btcec.ParsePubKey(pk)
	
	fmt.Println(addr)
	uaddr, _ := btcutil.NewAddressPubKey(pk, &chaincfg.MainNetParams);
	fmt.Println(uaddr)

	lc := NewBTCLightClient(&chaincfg.MainNetParams)
	balance, err := lc.GetBalance(&tx, "1MGXkMpTTA4Ue4wSDh4kGBKfLSwj93MHTq");
	fmt.Println(balance, err)
	
}
