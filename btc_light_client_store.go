package main

import "github.com/btcsuite/btcd/blockchain"

type BTCLightClientStorage interface {
	LightBlockAtHeight(int32) blockchain.HeaderCtx
	LatestHeight() uint64
	LatestLightBlock() blockchain.HeaderCtx
}





