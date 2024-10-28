package main

import (
	// "math/big"

	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

const MaxForkAge = 10

type Store interface {
	LightBlockAtHeight(int64) *LightBlock
	LatestHeight() int64
	LatestLightBlock() *LightBlock
	// SetHeader(height int64, header wire.BlockHeader) error
	LightBlockByHash(hash chainhash.Hash) *LightBlock
	RemindFork(latestBlockInFork chainhash.Hash) bool
	LatestBlockOfFork() []chainhash.Hash
	LatestCheckPoint() *LightBlock
	AddBlock(parent *LightBlock, header wire.BlockHeader) error
	SetLatestBlockOnFork(bh chainhash.Hash, latest bool) error
	TotalWorkAtBlock(bh chainhash.Hash) *big.Int
	SetBlock(lb *LightBlock, perviousPower *big.Int) 
}

type MemStore struct {
	latestHeight          int64
	lightblockMap         map[int64]*LightBlock
	lightBlockByHashMap   map[chainhash.Hash]*LightBlock
	latestBlockHashOfFork map[chainhash.Hash]struct{}
	totalWorkMap          map[chainhash.Hash]*big.Int
	latestcheckpoint *LightBlock
}

func NewMemStore() *MemStore {
	return &MemStore{
		latestHeight:          0,
		lightblockMap:         make(map[int64]*LightBlock),
		lightBlockByHashMap:   make(map[chainhash.Hash]*LightBlock),
		latestBlockHashOfFork: make(map[chainhash.Hash]struct{}),
		latestcheckpoint: nil,
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

func (s *MemStore) removeBlockByHash(hash chainhash.Hash) bool {
	if block := s.lightBlockByHashMap[hash]; block == nil {
		return false
	}

	delete(s.lightBlockByHashMap, hash)
	return true
}

// // this remove the block at current height/by hash and
// // this override a new light block by height/hash
// func (s *MemStore) SetHeader(height int64, header wire.BlockHeader) error {
// 	lightBlock := NewLightBlock(int32(height), header)

// 	if previousBlock := s.LightBlockAtHeight(height - 1); previousBlock != nil {
// 		lightBlock.TotalWork.Add(lightBlock.TotalWork, previousBlock.TotalWork)
// 	}

// 	headerHash := header.BlockHash()

// 	// remove the old hash if this exist in storage
// 	s.removeBlockByHash(headerHash)

// 	s.lightblockMap[height] = lightBlock
// 	s.lightBlockByHashMap[headerHash] = lightBlock

// 	if s.latestHeight < height {
// 		s.latestHeight = height
// 	}

// 	return nil
// }



func (s *MemStore) SetBlock(lb *LightBlock, previousPower *big.Int) {
	blockHash := lb.Header.BlockHash()
	s.lightBlockByHashMap[blockHash] = lb
	power := previousPower.Add(previousPower, lb.CalcWork())
	s.totalWorkMap[blockHash] = power
}

func (s *MemStore) AddBlock(parent *LightBlock, header wire.BlockHeader) error{
	height := parent.Height + 1

	newBlock := NewLightBlock(height, header)
	blockHash := header.BlockHash()
	// TODO: handle case block exist in db when add
	s.lightBlockByHashMap[blockHash] = newBlock
	prevTotalWork := s.TotalWorkAtBlock(parent.Header.BlockHash())

	s.SetBlock(newBlock, prevTotalWork)
	
	return nil
}

func (s *MemStore) SetLatestBlockOnFork(bh chainhash.Hash, latest bool) error {
	if latest {
		s.latestBlockHashOfFork[bh] = struct {}{}
	} else {
		delete(s.latestBlockHashOfFork, bh)
	}
	
	return nil
}

func(s *MemStore) TotalWorkAtBlock(hash chainhash.Hash) *big.Int{
	return s.totalWorkMap[hash]
}

func (s *MemStore) LatestBlockOfFork() []chainhash.Hash {
	return []chainhash.Hash{}
}

func (s *MemStore) LatestCheckPoint() *LightBlock {
	return s.latestcheckpoint
} 

func (s *MemStore) RemindFork(h chainhash.Hash) bool {
	_, ok := s.latestBlockHashOfFork[h]	
	return ok
}

