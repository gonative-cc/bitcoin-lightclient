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

func (lc *BTCLightClient) SetHeader(height int64, header wire.BlockHeader) error {
	return lc.btcStore.SetHeader(height, header)
}

// We assume we always insert valid header. Acctually, Cosmos can revert a state
// when module return error so this assumtion is reasonable
func (lc *BTCLightClient) InsertHeaders(headers []wire.BlockHeader) error {
	latestHeight := lc.btcStore.LatestHeight()

	return lc.insertHeaderStartAtHeight(uint64(latestHeight), headers)
}


func (lc *BTCLightClient) insertHeaderStartAtHeight(startHeight uint64, headers []wire.BlockHeader) error {
	height := startHeight + 1
	for i, header := range headers {
		if err := lc.CheckHeader(header); err != nil {
			return NewInvalidHeaderErr(header.BlockHash().String(), i)
		}

		// override a height
		lc.SetHeader(int64(height), header)
		height++
	}
	return nil
}

func (lc *BTCLightClient) sumTotalWork(startBlock *LightBlock, headers []wire.BlockHeader) *big.Int {
	totalWork := startBlock.TotalWork
	for _, header := range headers {
		totalWork.Add(totalWork, blockchain.CalcWork(header.Bits))
	}
	return totalWork
}

func (lc *BTCLightClient) HandleFork(headers []wire.BlockHeader) error {
	// find the light block match with first header
	firstHeader := headers[0]
	lightBlock := lc.btcStore.LightBlockByHash(firstHeader.BlockHash())
	if lightBlock == nil {
		return errors.New("Header doesn't belong to the chain")
	}
		currentTotalWork := lc.btcStore.LatestLightBlock()
		otherForkTotalWork := lc.computeTotalWorkFork(lightBlock, headers[1:])

		// comapre with current fork
		if otherForkTotalWork.Cmp(currentTotalWork.TotalWork) > 0 {
			return lc.insertHeaderStartAtHeight(uint64(lightBlock.Height), headers[1:])
		} else {
			return errors.New("Invalid fork")
		}

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

func (lc *BTCLightClient) CheckHeader(header wire.BlockHeader) error {
	noFlag := blockchain.BFNone
	latestLightBlock := lc.btcStore.LatestLightBlock()

	if err := blockchain.CheckBlockHeaderContext(&header, NewHeaderContext(latestLightBlock, lc.btcStore), noFlag, lc, true); err != nil {
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
	lcStore := NewMemStore()

	lc := &BTCLightClient{
		params:   params,
		btcStore: lcStore,
	}
	for id, header := range headers {
		lc.SetHeader(int64(id+start), header)
	}
	return lc
}
