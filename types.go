package main

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type BTCLightClient struct {
	params       *chaincfg.Params
	headers      []*wire.BlockHeader
	latestHeight uint64
}

func (lc *BTCLightClient) InsertHeaders(headers []*wire.BlockHeader) error {
	// check hash chain
	latestHeader := lc.LatestHeader()
	for _, header := range headers {
		prevBlockHash := latestHeader.BlockHash()
		if !header.PrevBlock.IsEqual(&prevBlockHash) {
			return errors.New("Invalid Headers")
		}
		latestHeader = header
	}
	return nil
}

func (lc *BTCLightClient) LatestHeader() *wire.BlockHeader {
	return lc.headers[lc.latestHeight]
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
