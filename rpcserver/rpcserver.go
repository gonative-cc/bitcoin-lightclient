package rpcserver

import (
	"net/http"
	"time"

	"github.com/gonative-cc/bitcoin-lightclient/btclightclient"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/rs/zerolog/log"
)

type Block struct {
	Hash   *chainhash.Hash
	Height int64
}

// Have a type with some exported methods
type RPCServerHandler struct {
	btcLC *btclightclient.BTCLightClient
}

func (h *RPCServerHandler) Ping(in int) int {
	return in
}

// txn to insert bitcoin block headers to light client
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

		// update last checkpointed block height
		h.btcLC.CleanUpFork()
	}

	return nil
}

func (h *RPCServerHandler) ContainsBTCBlock(blockHash *chainhash.Hash) (bool, error) {
	return h.btcLC.IsBlockPresent(*blockHash), nil
}

// returns the block height and hash of tip block stored in light client
func (h *RPCServerHandler) GetBTCHeaderChainTip() (Block, error) {
	latestFinalizedBlockHeight := h.btcLC.LatestFinalizedBlockHeight()
	latestFinalizedBlockHash := h.btcLC.LatestFinalizedBlockHash()

	latestFinalizedBlock := Block{
		Height: latestFinalizedBlockHeight,
		Hash:   &latestFinalizedBlockHash,
	}

	return latestFinalizedBlock, nil
}

// returns if the spvProof is valid or not
func (h *RPCServerHandler) VerifySPV(spvProof btclightclient.SPVProof) (btclightclient.SPVStatus, error) {
	checkSPV := h.btcLC.VerifySPV(spvProof)

	return checkSPV, nil
}

// NewRPCServer creates a new instance of the rpcServer and starts listening
func StartRPCServer(btcLC *btclightclient.BTCLightClient) error {
	rpcServer := jsonrpc.NewServer()
	serverHandler := &RPCServerHandler{
		btcLC: btcLC,
	}
	rpcServer.Register("RPCServerHandler", serverHandler)

	rpcServer.AliasMethod("ping", "RPCServerHandler.Ping")
	rpcServer.AliasMethod("insert_headers", "RPCServerHandler.InsertHeaders")
	rpcServer.AliasMethod("contains_btc_block", "RPCServerHandler.ContainsBTCBlock")
	rpcServer.AliasMethod("get_btc_header_chain_tip", "RPCServerHandler.GetBTCHeaderChainTip")
	rpcServer.AliasMethod("verify_spv", "RPCServerHandler.VerifySPV")

	server := &http.Server{
		Addr:         ":9797",
		Handler:      rpcServer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info().Msgf("RPC server running at: %s", server.Addr)

	return server.ListenAndServe()
}
