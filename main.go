package main

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

func main() {

	header_0, _ := NewBlockHeader([]byte("0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c"))

	header_1, _ := NewBlockHeader([]byte("010000006fe28c0ab6f1b372c1a6a246ae63f74f931e8365e15a089c68d6190000000000982051fd1e4ba744bbbe680e1fee14677ba1a3c3540bf7b1cdb606e857233e0e61bc6649ffff001d01e36299"))
	headers := []*wire.BlockHeader{header_0}

	btcLC := NewBTCLightClientWithData(&chaincfg.MainNetParams, headers)

	if err := btcLC.InsertHeaders([]*wire.BlockHeader{header_1}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Insert success")
	}
	btcLC.Status()
}
