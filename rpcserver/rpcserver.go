package rpcserver

import (
	"net/http/httptest"

	"github.com/gonative-cc/bitcoin-lightclient/btclightclient"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/rs/zerolog/log"
)

// Have a type with some exported methods
type RPCServerHandler struct {
	btcLC *btclightclient.BTCLightClient
}

func (h *RPCServerHandler) Ping(in int) int {
	return in
}

// txn to insert bitcoin block headers to babylon chain
func (h *RPCServerHandler) InsertHeaders(
	blockHeaders []*wire.BlockHeader,
) error {
	for _, blockHeader := range blockHeaders {
		if err := h.btcLC.InsertHeader(*blockHeader); err != nil {
			log.Err(err).Msg("Failed to insert block header")

			return err
		} else {
			log.Info().Msgf("Inserted block header %s", blockHeader.BlockHash())
		}
	}

	return nil
}

// returns the block height and hash of tip block stored in babylon chain
func (h *RPCServerHandler) GetBTCHeaderChainTip() (*chainhash.Hash, error) {
	latestBlockHash := h.btcLC.LatestBlockHash()

	return &latestBlockHash, nil
}

// NewRPCServer creates a new instance of the rpcServer and starts listening
func NewRPCServer(btcLC *btclightclient.BTCLightClient) *httptest.Server {
	rpcServer := jsonrpc.NewServer()
	serverHandler := &RPCServerHandler{
		btcLC: btcLC,
	}
	rpcServer.Register("RPCServerHandler", serverHandler)

	rpcServer.AliasMethod("ping", "RPCServerHandler.Ping")
	rpcServer.AliasMethod("insert_headers", "RPCServerHandler.InsertHeaders")
	rpcServer.AliasMethod("get_btc_header_chain_tip", "RPCServerHandler.GetBTCHeaderChainTip")

	return httptest.NewServer(rpcServer)
}
