package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type BTCLightClient struct {
	params   *chaincfg.Params
	btcStore Store
}

func NewBTCLightClient(params *chaincfg.Params) *BTCLightClient {
	return &BTCLightClient{
		params:   params,
		btcStore: NewMemStore(),
	}
}

func (lc *BTCLightClient) ChainParams() *chaincfg.Params {
	return lc.params
}

func (lc *BTCLightClient) BlocksPerRetarget() int32 {
	return int32(lc.params.TargetTimespan / time.Second)
}

func (lc *BTCLightClient) MinRetargetTimespan() int64 {
	return int64(lc.BlocksPerRetarget()) / lc.params.RetargetAdjustmentFactor
}

func (lc *BTCLightClient) MaxRetargetTimespan() int64 {
	return int64(lc.BlocksPerRetarget()) * lc.params.RetargetAdjustmentFactor
}

func (lc *BTCLightClient) VerifyCheckpoint(height int32, hash *chainhash.Hash) bool {
	return false
}

func (lc *BTCLightClient) FindPreviousCheckpoint() (blockchain.HeaderCtx, error) {
	return nil, nil
}

// We assume we always insert valid header. Acctually, Cosmos can revert a state
// when module return error so this assumtion is reasonable
func (lc *BTCLightClient) InsertHeader(header wire.BlockHeader) error {

	if lb := lc.btcStore.LightBlockByHash(header.BlockHash()); lb != nil {
		return errors.New("Parent block not found")
	}

	parentHash := header.PrevBlock
	parent := lc.btcStore.LightBlockByHash(parentHash)
	if parent == nil {
		return errors.New("Block doesn't belong to any fork!")
	}

	// we need to handle 2 cases:
	// extend the exist fork
	if lc.btcStore.IsForkHead(parentHash) {
		if err := lc.CheckHeader(parent.Header, header); err != nil {
			return err

		}

		lc.btcStore.AddBlock(parent, header)
		lc.btcStore.SetLatestBlockOnFork(parentHash, false)
		lc.btcStore.SetLatestBlockOnFork(header.BlockHash(), true)
		return nil
	}

	// create a new fork
	return lc.CreateNewFork(parent, header)
}

// TODO: fix this function in the next PR
// There can be few heads. Example, where we have 3 active forks, with 3 latest blocks (b3, c3, d3):
// b1 <- b2  <- b3
//    |- c2  <- c3
//    |- d2'  <- d3
func (lc *BTCLightClient) findLightBlock(bh chainhash.Hash) ([]*LightBlock, error) {
	checkpoint := lc.btcStore.LatestCheckPoint()
	checkpointHash := checkpoint.Header.BlockHash()
	fork := make([]*LightBlock, 0)
	for count := 0; count <= MaxForkAge; count++ {
		curr := lc.btcStore.LightBlockByHash(bh)
		fork = append(fork, curr)
		if bh.IsEqual(&checkpointHash) {
			return fork, nil
		}
		bh = curr.Header.PrevBlock
	}

	return nil, errors.New("Block not found or too old")
}

// TODO: We will do this in the next PR
// - Remove all fork invalid
// - Update map(height => block) in btcStore
func (lc *BTCLightClient) CleanupFork() {
}

func (lc *BTCLightClient) CreateNewFork(parent *LightBlock, header wire.BlockHeader) error {
	if lc.CheckHeader(parent.Header, header) != nil {
		// add to db
		lc.btcStore.AddBlock(parent, header)
		// update block as latest block in this fork
		lc.btcStore.SetLatestBlockOnFork(header.BlockHash(), true)
	}
	return nil
}

type BlockMedianTimeSource struct {
	h *wire.BlockHeader
}

func newBlockMedianTimeSource(header *wire.BlockHeader) *BlockMedianTimeSource {
	return &BlockMedianTimeSource{
		h: header,
	}
}

func (b *BlockMedianTimeSource) AdjustedTime() time.Time {
	return b.h.Timestamp
}

func (b *BlockMedianTimeSource) AddTimeSample(string, time.Time) {
	// We only verify header, so we don't need do anything here
}

func (b *BlockMedianTimeSource) Offset() time.Duration {
	// don't need to update any
	return 0
}

func (lc *BTCLightClient) CheckHeader(parent wire.BlockHeader, header wire.BlockHeader) error {
	noFlag := blockchain.BFNone
	fork, _ := lc.findLightBlock(parent.BlockHash())
	latestLightBlock := fork[0]

	if err := blockchain.CheckBlockHeaderContext(&header, NewHeaderContext(latestLightBlock, lc.btcStore, fork), noFlag, lc, true); err != nil {
		return err
	}

	if err := blockchain.CheckBlockHeaderSanity(&header, lc.params.PowLimit, newBlockMedianTimeSource(&header), noFlag); err != nil {
		return err
	}
	return nil
}

// query status, use for test
func (lc *BTCLightClient) Status() {
	fmt.Println(lc.params.Net)
	latestBlock := lc.btcStore.LatestCheckPoint()
	fmt.Println(latestBlock.Height)
	fmt.Println(lc.btcStore.MostDifficultFork())
}

func NewBTCLightClientWithData(params *chaincfg.Params, headers []wire.BlockHeader, start int) *BTCLightClient {
	lc := NewBTCLightClient(params)
	lb := NewLightBlock(int32(start), headers[0])
	lc.btcStore.SetBlock(lb, big.NewInt(0))
	lc.btcStore.SetLatestCheckPoint(lb)
	lc.btcStore.SetLightBlockByHeight(lb)
	for i, header := range headers[1:] {
		previousPower := lc.btcStore.TotalWorkAtBlock(header.PrevBlock)
		lb := NewLightBlock(int32(start+i+1), header)
		lc.btcStore.SetBlock(lb, previousPower)

	}
	return lc
}
