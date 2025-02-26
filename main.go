package main

import (
	"os"

	"github.com/gonative-cc/bitcoin-lightclient/btclightclient"
	"github.com/gonative-cc/bitcoin-lightclient/rpcserver"

	"github.com/gonative-cc/bitcoin-lightclient/data"

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
	sampleFilename := os.Args[1]

	if _, err := os.Stat(sampleFilename); os.IsNotExist(err) {
		log.Error().Msgf("Sample file does not exist: %s", sampleFilename)
		return
	}

	networkParams, startHeight, blockHeaders, err := data.ReadJSON(sampleFilename)
	if err != nil {
		log.Error().Msgf("Error reading data file: %s", err)
		return
	}
	headers := make([]wire.BlockHeader, len(blockHeaders))

	for id, headerStr := range blockHeaders {
		h, _ := btclightclient.BlockHeaderFromHex(headerStr)
		headers[id] = h
	}

	btcLC := btclightclient.NewBTCLightClientWithData(networkParams, headers, int(startHeight))
	btcLC.Status()

	err = rpcserver.StartRPCServer(btcLC)
	if err != nil {
		log.Error().Msgf("Error creating RPC server: %s", err)
		return
	}
}
