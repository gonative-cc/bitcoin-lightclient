package main

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)

type LightBlock struct {
	height int32
	header *wire.BlockHeader
	store  BTCLightClientStorage
}

func (lb LightBlock) Height() int32 {
	return lb.height
}

func (lb LightBlock) Bits() uint32 {
	return lb.header.Bits
}

func (lb LightBlock) TimeStamp() int64 {
	return lb.header.Timestamp.Unix()
}

func (lb LightBlock) Parent() blockchain.HeaderCtx {
	if lb.height > 0 {
		parentHeight := lb.height - 1
		return lb.store.LightBlockAtHeight(parentHeight)
	}
	return nil
}

func (l *LightBlock) RelativeAncestorCtx(
	distance int32) blockchain.HeaderCtx {
	return nil
}
