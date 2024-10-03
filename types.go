package main

import (
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
	return nil
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
