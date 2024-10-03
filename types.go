package main

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)


type BTCLightClientStorage interface {
	LightBlockAtHeight(int32) blockchain.HeaderCtx
}

type BTCLightClient struct {
	params       *chaincfg.Params
	headers      []*wire.BlockHeader
	lightBlocks  map[uint64]*LightBlock
	latestHeight uint64
}

func (lc *BTCLightClient) InsertHeaders(headers []*wire.BlockHeader) error {

	return nil
}

func (lc *BTCLightClient) GetLightBlock(height uint64) *LightBlock {
	if height > lc.latestHeight {
		return nil
	}

	return lc.lightBlocks[height]
}

func (lc *BTCLightClient) CheckHeader(header *wire.BlockHeader) bool {

	return true
}

func NewBTCLightClient(params chaincfg.Params) *BTCLightClient {
	return &BTCLightClient{
		params:       &params,
		headers:      []*wire.BlockHeader{},
		latestHeight: 0,
	}
}

func (lc BTCLightClient) Status() {
	fmt.Println(lc.params.Net)
	fmt.Println(lc.params.Name)
	fmt.Println("Status of BTC light client")
}
