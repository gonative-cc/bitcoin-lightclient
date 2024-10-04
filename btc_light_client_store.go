package main

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)


var _ blockchain.HeaderCtx = (*LightBlock)(nil)

type BTCLightClientStorage interface {
	LightBlockAtHeight(int64) blockchain.HeaderCtx
	LatestHeight() int64
	LatestLightBlock() blockchain.HeaderCtx
	AddHeader(height int64, header *wire.BlockHeader) error 
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

func (lcStore *LCStorage) AddHeader(height int64, header *wire.BlockHeader) error {
	lightBlock := NewLightBlock(int32(height), header, lcStore)
	lcStore.latestHeight = height
	lcStore.lightblockMap[height] = lightBlock
	return nil
}
