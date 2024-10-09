package main

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type Store interface {
	LightBlockAtHeight(int64) *LightBlock
	LatestHeight() int64
	LatestLightBlock() *LightBlock
	AddHeader(height int64, header wire.BlockHeader) error
	CurrentTotalWork() *big.Int
	LightBlockByHash(hash chainhash.Hash) *LightBlock 
}

type MemStore struct {
	latestHeight  int64
	lightblockMap map[int64]*LightBlock
	lightBlockByHashMap map[chainhash.Hash] *LightBlock
}

func NewMemStore() *MemStore {
	return &MemStore{
		latestHeight:  0,
		lightblockMap: make(map[int64]*LightBlock),
		lightBlockByHashMap: make(map[chainhash.Hash]*LightBlock),
	}
}

func (s *MemStore) LightBlockAtHeight(height int64) *LightBlock {
	return s.lightblockMap[height]
}

func (s *MemStore) LatestHeight() int64 {
	return s.latestHeight
}

func (s *MemStore) LatestLightBlock() *LightBlock {
	return s.lightblockMap[s.LatestHeight()]
}


func (s *MemStore) LightBlockByHash(hash chainhash.Hash) *LightBlock {
	return s.lightBlockByHashMap[hash]
}

func (s *MemStore) AddHeader(height int64, header wire.BlockHeader) error {
	lightBlock := NewLightBlock(int32(height), header)
	
	if previousBlock := s.LightBlockAtHeight(height - 1); previousBlock != nil {
		lightBlock.TotalWork.Add(lightBlock.TotalWork, previousBlock.TotalWork)
	}

	
	s.lightblockMap[height] = lightBlock
	
	headerHash := header.BlockHash()
	
	s.lightBlockByHashMap[headerHash] = lightBlock
	if s.latestHeight < height {
		s.latestHeight = height
	}
	
	return nil
}


func (s *MemStore) CurrentTotalWork() *big.Int{
	return s.LatestLightBlock().TotalWork
}
