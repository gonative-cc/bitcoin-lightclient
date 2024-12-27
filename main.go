package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gonative-cc/bitcoin-lightclient/btclightclient"
	"github.com/gonative-cc/bitcoin-lightclient/rpcserver"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rs/zerolog/log"
)

func main() {
	// read the json file
	// example: ./data/sample.json
	if len(os.Args) < 2 {
		log.Error().Msg("Missing filename.\nUsage: bitcoin-lightclient <sample_file.json>")
		return
	}
	dataFilePath := os.Args[1]

	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		log.Error().Msgf("Data file does not exist: %s", dataFilePath)
		return
	}

	startHeight, headerStrings, err := ReadJson(dataFilePath)
	if err != nil {
		log.Error().Msgf("Error reading data file: %s", err)
		return
	}
	headers := make([]wire.BlockHeader, len(headerStrings))

	for id, headerStr := range headerStrings {
		h, _ := btclightclient.BlockHeaderFromHex(headerStr)
		headers[id] = h
	}

	btcLC := btclightclient.NewBTCLightClientWithData(&chaincfg.MainNetParams, headers, int(startHeight))
	btcLC.Status()

	rpcService := rpcserver.NewRPCServer(btcLC)
	log.Info().Msgf("RPC server running at: %s", rpcService.URL)

	// Create channel to listen for interrupt signal
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
