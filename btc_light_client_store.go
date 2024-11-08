package main

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

const MaxForkAge = 12

type Store interface {
	LightBlockAtHeight(int64) *LightBlock
	LatestHeight() int64
	LightBlockByHash(hash chainhash.Hash) *LightBlock
	// check hash h is hash of latest block in remind fork sets.
	IsForkHead(h chainhash.Hash) bool
	LatestCheckPoint() *LightBlock
	AddBlock(parent *LightBlock, header wire.BlockHeader) error
	SetLatestBlockOnFork(bh chainhash.Hash, latest bool) error
	TotalWorkAtBlock(bh chainhash.Hash) *big.Int
	SetBlock(lb *LightBlock, perviousPower *big.Int)
	SetLatestCheckPoint(lb *LightBlock)
	SetLightBlockByHeight(lb *LightBlock)
	MostDifficultFork() *LightBlock
}

type MemStore struct {
	lightblockMap         map[int64]*LightBlock
	lightBlockByHashMap   map[chainhash.Hash]*LightBlock
	latestBlockHashOfFork map[chainhash.Hash]struct{}
	totalWorkMap          map[chainhash.Hash]*big.Int
	latestcheckpoint      *LightBlock
	mostDifficultFork     *LightBlock
}

func NewMemStore() *MemStore {
	return &MemStore{
		lightblockMap:         make(map[int64]*LightBlock),
		lightBlockByHashMap:   make(map[chainhash.Hash]*LightBlock),
		latestBlockHashOfFork: make(map[chainhash.Hash]struct{}),
		totalWorkMap:          make(map[chainhash.Hash]*big.Int),
		latestcheckpoint:      nil,
		mostDifficultFork:     nil,
	}
}

func (s *MemStore) SetLightBlockByHeight(lb *LightBlock) {
	s.lightblockMap[int64(lb.Height)] = lb
}

func (s *MemStore) LightBlockAtHeight(height int64) *LightBlock {
	return s.lightblockMap[height]
}

func (s *MemStore) LatestHeight() int64 {
	return int64(s.latestcheckpoint.Height)
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

func (s *MemStore) SetLatestCheckPoint(lb *LightBlock) {
	s.latestcheckpoint = lb
}

func (s *MemStore) SetBlock(lb *LightBlock, previousPower *big.Int) {
	blockHash := lb.Header.BlockHash()
	s.lightBlockByHashMap[blockHash] = lb

	power := big.NewInt(0)
	power = power.Add(previousPower, lb.CalcWork())

	powerForkBlock := s.MostDifficultFork()
	mostPower := big.NewInt(0)
	if powerForkBlock != nil {
		mostPower = s.totalWorkMap[powerForkBlock.Header.BlockHash()]
	}

	if mostPower.Cmp(power) < 0 {
		s.mostDifficultFork = lb
	}

	s.totalWorkMap[blockHash] = power
}

func (s *MemStore) AddBlock(parent *LightBlock, header wire.BlockHeader) error {
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
		s.latestBlockHashOfFork[bh] = struct{}{}
	} else {
		delete(s.latestBlockHashOfFork, bh)
	}

	return nil
}

func (s *MemStore) TotalWorkAtBlock(hash chainhash.Hash) *big.Int {
	return s.totalWorkMap[hash]
}

func (s *MemStore) LatestCheckPoint() *LightBlock {
	return s.latestcheckpoint
}

func (s *MemStore) IsForkHead(h chainhash.Hash) bool {
	_, ok := s.latestBlockHashOfFork[h]
	return ok
}

func (s *MemStore) MostDifficultFork() *LightBlock {
	return s.mostDifficultFork
}
