package main

import "github.com/btcsuite/btcd/blockchain"


var _ blockchain.HeaderCtx = (*LightBlock)(nil)

type BTCLightClientStorage interface {
	LightBlockAtHeight(int64) blockchain.HeaderCtx
	LatestHeight() int64
	LatestLightBlock() blockchain.HeaderCtx
}

type LCStorage struct {
	latestHeight  int64
	lightblockMap map[int64]*LightBlock
}

func NewLCStorage() *LCStorage {
	return &LCStorage{
		latestHeight:  0,
		lightblockMap: make(map[int64]*LightBlock),
	}
}


func (lcStore *LCStorage) LightBlockAtHeight(height int64) blockchain.HeaderCtx {
	return lcStore.lightblockMap[height]
}

func (lcStore *LCStorage) LatestHeight() int64 {
	return lcStore.latestHeight
}

func (lcStore *LCStorage) LatestLightBlock() blockchain.HeaderCtx {
	return lcStore.lightblockMap[lcStore.LatestHeight()]
}
