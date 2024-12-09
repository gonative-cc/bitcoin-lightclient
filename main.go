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
	startHeight := 204100
	headerStrings := []string{
		`020000004cdba1415b2c6e7808c1b3c18df1374238454f7104203475bf01000000000000c17ea9d06015dc83902911cd24837a8ba4bdc0c1d72b873f906d921e06e48d2f984a8250ef75051a72a8061a`,
		`01000000a7e267eae437f7c99115682d3e104999b6b5b01f34f3f5d763030000000000004beb49a1caa3941adc853cd605dd030e5ab690fcd58ac76233397653c49f07ecdb4a8250ef75051a5e58c1fe`,
		`0100000032e89326a365919ac08d1f0bd98d097240fd5c8dd21b7c38e5030000000000003d3645c6e6008583570d618041b2452497d55c0cb936529af92f8191e0dbdcc9764b8250ef75051ac733c738`,
		`010000006cb4d23bba90d14f11d0719a02bdb0218e7b9a82a4453e0a51010000000000007caa984f729a85cec6ce000b87894bb30897195d36bc6491be92e7628de00d02e54d8250ef75051ac9b86274`,
		`01000000a5a2223dc3f0a3f7901d14d064216d6ecf9d5dabfa44199f2b02000000000000e877c5cb9c52e27e3b2b9ff5d3b1658bd5c49ce2e3a1ef608759b44eea5f832b424e8250ef75051a47a25741`,
		`0100000096bacf353c5534bc30044268588fb752df9c164a3ad14bd64d04000000000000a25ec7252a9a994d864b664bb4c5a2b1d8c1d688dc82ab78128eae34481110f416518250ef75051a1ac7889f`,
		`020000009c37b58365fedc02e5ce0fa879e14d3ec110a3700b6476b941000000000000007e82e8388e953e200ab3e41777e4f9982ce3fe3aa46cc29a8c8f974706f76f852b508250ef75051a9e6f81c1`,
		`010000001f7db2278bb6aeac530fa48e0f8e75cae092a89b6cbedb838500000000000000d7dfbe2005869b3319244192ba2fef1b0c9f1cf999da51be7d8f1a7f5bb02c6bf0538250ef75051a82d85a5f`,
		`01000000575cdd3be997ebc63b377e78ecde7f14ef354669ff7d62543703000000000000a7b0a3354071262e02c246ad1d49d6823f6d864aef247263a07641a1d419089b5a5a8250ef75051a1e7d568e`,
		`01000000fb49bbe46f780f8b0aa4b99ceddf1202e76d12fa4f5da1acdf000000000000003ec200da269c721728c820e305d01d4656e4ba0f6784c4fdeea05a904e0f0aab395d8250ef75051ac332b7b5`,
		`02000000418574a33d0657d095a8ff4723e8e1fafedcabe32b93d1ccf6010000000000006f95d48f47f27715f6c1bcc466c4c715a17f9e2be39ce1c41106baed2d146d546f5f8250ef75051a77fc3e40`,
		`01000000d19ffbe9a876f329acb05feadc90dae27578762e3f538f5a08050000000000005c8524057100e644fb3680b5b5af3d79f24fc4e51f374d9205e926e9b9cafd9c38648250ef75051a20379982`,
		`02000000670ffede95831fb41c09d0d104285d6182ce1c8577da40506405000000000000e54435f50bfc776b8f3d9ac047963ee6bdddd8d40b69236b4d97acb52a1fdce41e678250ef75051a88842656`,
	}

	headers := make([]wire.BlockHeader, len(headerStrings))

	for id, headerStr := range headerStrings {
		h, _ := btclightclient.BlockHeaderFromHex(headerStr)
		headers[id] = h
	}

	headerInsert, _ := btclightclient.BlockHeaderFromHex("01000000ea0ec14effa5f7f2a1a9f4431588b63b575d167a261c1d93b604000000000000c1844859aa7bb44251cf04a19098169f657e4bd91ebeb3f2a028211f1f8bde271c6e8250ef75051a7dc08785")

	btcLC := btclightclient.NewBTCLightClientWithData(&chaincfg.MainNetParams, headers, startHeight)
	btcLC.Status()

	rpcService := rpcserver.NewRPCServer(btcLC)
	log.Info().Msgf("RPC server running at: %s", rpcService.URL)

	if err := btcLC.InsertHeader(headerInsert); err != nil {
		log.Err(err).Msg("Failed to insert block header")
	} else {
		log.Info().Msgf("Inserted block header %s", headerInsert.BlockHash())
	}

	err := btcLC.CleanUpFork()
	if err != nil {
		log.Err(err).Msg("Failed to clean up fork")
	}
	btcLC.Status()

	// Create channel to listen for interrupt signal
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
