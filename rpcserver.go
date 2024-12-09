package main

import (
	"net/http/httptest"

	// "github.com/babylonchain/babylon/x/btclightclient/types"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/rs/zerolog/log"
)

// Have a type with some exported methods
type RPCServerHandler struct {
	btcLC *BTCLightClient
}

func (h *RPCServerHandler) Ping(in int) int {
	return in
}

// txn to insert bitcoin block headers to babylon chain
func (h *RPCServerHandler) InsertHeaders(
	blockHeaders []*wire.BlockHeader,
) (*chainhash.Hash, error) {
	for _, blockHeader := range blockHeaders {
		if err := h.btcLC.InsertHeader(*blockHeader); err != nil {
			log.Err(err).Msg("Failed to insert block header")

			return nil, err
		} else {
			log.Info().Msgf("Inserted block header %s", blockHeader.BlockHash())
		}
	}

	return nil, nil
}

// returns the block height and hash of tip block stored in babylon chain
func (h *RPCServerHandler) GetBTCHeaderChainTip() (*chainhash.Hash, error) {
	latestBlockHash := h.btcLC.btcStore.LatestCheckPoint().Header.BlockHash()

	return &latestBlockHash, nil
}

// NewRPCServer creates a new instance of the rpcServer and starts listening
func NewRPCServer(btcLC *BTCLightClient) *httptest.Server {
	// Create a new RPC server
	rpcServer := jsonrpc.NewServer()

	// create a handler instance and register it
	serverHandler := &RPCServerHandler{
		btcLC: btcLC,
	}
	rpcServer.Register("RPCServerHandler", serverHandler)

	rpcServer.AliasMethod("ping", "RPCServerHandler.Ping")
	rpcServer.AliasMethod("insert_headers", "RPCServerHandler.InsertHeaders")
	rpcServer.AliasMethod("get_btc_header_chain_tip", "RPCServerHandler.GetBTCHeaderChainTip")

	return httptest.NewServer(rpcServer)
}
