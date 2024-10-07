package main

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)

type LightBlock struct {
	height  int32
	header  *wire.BlockHeader
	lcStore Store
}

func (lb LightBlock) Height() int32 {
	return lb.height
}

func (lb LightBlock) Bits() uint32 {
	return lb.header.Bits
}

func (lb LightBlock) Timestamp() int64 {
	return lb.header.Timestamp.Unix()
}

func (lb LightBlock) Parent() blockchain.HeaderCtx {
	return lb.RelativeAncestorCtx(1)
}

func (lb *LightBlock) RelativeAncestorCtx(
	distance int32) blockchain.HeaderCtx {
	if distance <= lb.Height() {
		ancestorHeight := lb.Height() - distance
		return lb.lcStore.LightBlockAtHeight(int64(ancestorHeight))
	}
	return nil
}

func NewLightBlock(height int32, header *wire.BlockHeader, lcStore Store) *LightBlock {
	return &LightBlock{
		height: height,
		header: header,
		lcStore: lcStore,
	}
}
