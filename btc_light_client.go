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

func (lc *BTCLightClient) AddHeader(height int64, header wire.BlockHeader) error {
	return lc.btcStore.AddHeader(height, header)
}

// We assume we always insert valid header. Acctually, Cosmos can revert a state
// when module return error so this assumtion is reasonable
func (lc *BTCLightClient) InsertHeaders(headers []wire.BlockHeader) error {
	latestHeight := lc.btcStore.LatestHeight()

	return lc.insertHeadersWithPosition(uint64(latestHeight), headers)
}

func (lc *BTCLightClient) insertHeadersWithPosition(height uint64, headers []wire.BlockHeader) error {
	insertHeight := height + 1
	for _, header := range headers {
		if err := lc.CheckHeader(header); err != nil {
			return err
		}

		lc.AddHeader(int64(insertHeight), header)
		insertHeight = insertHeight + 1
	}
	return nil
}

func (lc *BTCLightClient) HandleFork(headers []wire.BlockHeader) error {
	// find the light block match with first header
	firstHeader := headers[0]
	if lightBlock := lc.btcStore.LightBlockByHash(firstHeader.BlockHash()); lightBlock != nil {
		latestBlock := lc.btcStore.LatestLightBlock()
		totalWorkOnSecondChain := lightBlock.TotalWork
		for _, header := range headers[1:] {
			totalWorkOnSecondChain.Add(totalWorkOnSecondChain, blockchain.CalcWork(header.Bits))
		}

		if totalWorkOnSecondChain.Cmp(latestBlock.TotalWork) > 0 {
			return lc.insertHeadersWithPosition(uint64(lightBlock.Height), headers[1:])
		} else {
			return errors.New("Invalid fork")
		}
	}
	return errors.New("Header not belong to chain")

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
		lc.AddHeader(int64(id+start), header)
	}
	return lc
}
