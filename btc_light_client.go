package main

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

var _ blockchain.HeaderCtx = (*LightBlock)(nil)

type BTCLightClient struct {
	params      *chaincfg.Params
	btc_storage BTCLightClientStorage
}

func NewBTCLightClient(params *chaincfg.Params) *BTCLightClient {
	return &BTCLightClient{
		params:      params,
		btc_storage: NewLCStorage(),
	}
}
func (lc BTCLightClient) ChainParams() *chaincfg.Params {
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

func (lc *BTCLightClient) InsertHeaders(headers []*wire.BlockHeader) error {
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

func (b *BlockMedianTimeSource) AddTimeSample(_ string, _time time.Time) {
	// We only verify header, so we don't need do anything here
}

func (b *BlockMedianTimeSource) Offset() time.Duration {
	// don't need to update any
	return 0
}

func (lc *BTCLightClient) CheckHeader(header *wire.BlockHeader) error {
	noFlag := blockchain.BFNone
	if err := blockchain.CheckBlockHeaderContext(header, lc.btc_storage.LatestLightBlock(), noFlag, lc, false); err != nil {
		return err
	}

	if err := blockchain.CheckBlockHeaderSanity(header, lc.params.PowLimit, newBlockMedianTimeSource(header), noFlag); err != nil {
		return err
	}
	return nil
}

func (lc BTCLightClient) Status() {
	fmt.Println(lc.params.Net)
	fmt.Println(lc.params.Name)
	fmt.Println("Status of BTC light client")
}
