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
			log.Err(err).Msgf("Failed to insert block header %s", blockHeader.BlockHash())

			return err
		} else {
			log.Info().Msgf("Inserted block header %s", blockHeader.BlockHash())
		}

		// fork cleanup updates last checkpointed height
		if err := h.btcLC.CleanUpFork(); err != nil {
			log.Err(err).Msgf(
				"Failed to update fork choice after inserting block header %s",
				blockHeader.BlockHash(),
			)

			return err
		} else {
			log.Info().Msgf(
				"Updated fork choice after inserting block header %s",
				blockHeader.BlockHash(),
			)
		}
	}

	return nil
}

func (h *RPCServerHandler) ContainsBTCBlock(blockHash *chainhash.Hash) (bool, error) {
	return h.btcLC.IsBlockPresent(*blockHash), nil
}

// GetHeaderChainTip returns the latest finalized block stored in light client
func (h *RPCServerHandler) GetHeaderChainTip() (Block, error) {
	latestFinalizedBlockHeight := h.btcLC.LatestFinalizedBlockHeight()
	latestFinalizedBlockHash := h.btcLC.LatestFinalizedBlockHash()

	latestFinalizedBlock := Block{
		Height: latestFinalizedBlockHeight,
		Hash:   &latestFinalizedBlockHash,
	}

	return latestFinalizedBlock, nil
}

// Verify SPV proof for bitcoin transaction inclusive to a block
func (h *RPCServerHandler) VerifySPV(spvProof *btclightclient.SPVProof) (btclightclient.SPVStatus, error) {
	log.Debug().Msgf("Recieved spvProof %v", spvProof)
	checkSPV := h.btcLC.VerifySPV(*spvProof)

	log.Info().Msgf("SPV proof: %v status %v", spvProof, checkSPV)

	return checkSPV, nil
}

// Verify a bunch SPV proof for bitcoin transactions inclusive to a block
func (h *RPCServerHandler) VerifySPVs(spvProofs []btclightclient.SPVProof) ([]btclightclient.SPVStatus, error) {
	log.Debug().Msgf("Received list of SPV %v", spvProofs)
	status := h.btcLC.VerifySPVs(spvProofs)

	log.Info().Msgf("SPV status status: %v status %v", spvProofs, status)

	return status, nil
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
	rpcServer.AliasMethod("get_header_chain_tip", "RPCServerHandler.GetHeaderChainTip")
	rpcServer.AliasMethod("verify_spv", "RPCServerHandler.VerifySPV")
	rpcServer.AliasMethod("verify_spvs", "RPCServerHandler.VerifySPVs")

	server := &http.Server{
		Addr:         ":9797",
		Handler:      rpcServer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info().Msgf("RPC server running at: %s", server.Addr)

	return server.ListenAndServe()
}
