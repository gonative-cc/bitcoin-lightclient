package main

import (
	"errors"
	"fmt"
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

// func (lc *BTCLightClient) SetHeader(height int64, header wire.BlockHeader) error {
// 	return lc.btcStore.SetHeader(height, header)
// }

// We assume we always insert valid header. Acctually, Cosmos can revert a state
// when module return error so this assumtion is reasonable
func (lc *BTCLightClient) InsertHeaders(header wire.BlockHeader) error {
	previous := header.PrevBlock

	previousBlock := lc.btcStore.LightBlockByHash(previous)

	if previousBlock == nil {
		return errors.New("Block doesn't belong to any fork!")
	}
	// we need to handle 2 cases:

	// extend the exist fork
	if lc.btcStore.RemindFork(previous) {
		previousHeader := previousBlock.Header
		if err := lc.CheckHeader(previousHeader, header); err != nil {
			return err

		}
		
		lc.btcStore.AddBlock(previousBlock, header)
		lc.btcStore.SetLatestBlockOnFork(previous, false)
		lc.btcStore.SetLatestBlockOnFork(header.BlockHash(), true)
		return nil
	}

	// create a new fork
	parent := lc.btcStore.LightBlockByHash(previous)
	return lc.CreateNewFork(parent, header)
}


func (lc *BTCLightClient) extractFork(lastBlock chainhash.Hash) ([]*LightBlock, error) {
	checkpoint := lc.btcStore.LatestCheckPoint()
	checkpointHash := checkpoint.Header.BlockHash()
	count := 0
	fork := make([]*LightBlock, 0)

	for count <= MaxForkAge {
		if lastBlock.IsEqual(&checkpointHash) {
			return fork, nil
		}
		count++
	}

	return nil, errors.New("Fork age invalid")
}

// func (lc *BTCLightClient) insertHeaderStartAtHeight(startHeight uint64, headers []wire.BlockHeader) error {
// 	height := startHeight + 1
// 	for i, header := range headers {
// 		if err := lc.CheckHeader(header); err != nil {
// 			return NewInvalidHeaderErr(header.BlockHash().String(), i)
// 		}

// 		// override a height
// 		lc.SetHeader(int64(height), header)
// 		height++
// 	}
// 	return nil
// }

// func (lc *BTCLightClient) sumTotalWork(startBlock *LightBlock, headers []wire.BlockHeader) *big.Int {
// 	totalWork := startBlock.TotalWork
// 	for _, header := range headers {
// 		totalWork.Add(totalWork, blockchain.CalcWork(header.Bits))
// 	}
// 	return totalWork
// }

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
	fork, _ := lc.extractFork(parent.BlockHash())
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
	latestBlock := lc.btcStore.LatestLightBlock()
	fmt.Println(latestBlock.Height)
}

func NewBTCLightClientWithData(params *chaincfg.Params, headers []wire.BlockHeader, start int) *BTCLightClient {
	lc := NewBTCLightClient(params)
	lb := NewLightBlock(int32(start), headers[0])
	for _, header := range headers[1:] {
		lc.btcStore.AddBlock(lb, header)

	}
	return lc
}
