package main

import "github.com/btcsuite/btcd/chaincfg"


func main() {
	btcLC := NewBTCLightClient(&chaincfg.MainNetParams)
	btcLC.Status()
}
