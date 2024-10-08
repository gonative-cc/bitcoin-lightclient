package main

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/wire"
)

var _ blockchain.HeaderCtx = (*HeaderContext)(nil)

type LightBlock struct {
	Height int32
	Header wire.BlockHeader
}

type HeaderContext struct {
	lightBlock *LightBlock
	store      Store
}

func (h *HeaderContext) Height() int32 {
	return h.lightBlock.Height
}

func (h *HeaderContext) Bits() uint32 {
	return h.lightBlock.Header.Bits
}

func (h *HeaderContext) Timestamp() int64 {
	return h.lightBlock.Header.Timestamp.Unix()
}

func (h *HeaderContext) Parent() blockchain.HeaderCtx {
	return h.RelativeAncestorCtx(1)
}

func (h *HeaderContext) RelativeAncestorCtx(
	distance int32) blockchain.HeaderCtx {
	if distance <= h.Height() {
		ancestorHeight := h.Height() - distance
		blockAtHeight := h.store.LightBlockAtHeight(int64(ancestorHeight))
		return NewHeaderContext(blockAtHeight, h.store)
	}
	return nil
}

func NewLightBlock(height int32, header wire.BlockHeader) *LightBlock {
	return &LightBlock{
		Height: height,
		Header: header,
	}
}

func NewHeaderContext(lightBlock *LightBlock, store Store) *HeaderContext {
	return &HeaderContext{
		lightBlock: lightBlock,
		store:      store,
	}
}
