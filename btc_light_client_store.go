package main

import (
	"github.com/btcsuite/btcd/wire"
)

type Store interface {
	LightBlockAtHeight(int64) *LightBlock
	LatestHeight() int64
	LatestLightBlock() *LightBlock
	AddHeader(height int64, header wire.BlockHeader) error
}

type MemStore struct {
	latestHeight  int64
	lightblockMap map[int64]*LightBlock
}

func NewMemStore() *MemStore {
	return &MemStore{
		latestHeight:  0,
		lightblockMap: make(map[int64]*LightBlock),
	}
}

func (lcStore *MemStore) LightBlockAtHeight(height int64) *LightBlock {
	return lcStore.lightblockMap[height]
}

func (lcStore *MemStore) LatestHeight() int64 {
	return lcStore.latestHeight
}

func (lcStore *MemStore) LatestLightBlock() *LightBlock {
	return lcStore.lightblockMap[lcStore.LatestHeight()]
}

func (lcStore *MemStore) AddHeader(height int64, header wire.BlockHeader) error {
	lightBlock := NewLightBlock(int32(height), header)
	lcStore.latestHeight = height
	lcStore.lightblockMap[height] = lightBlock
	return nil
}
