package btclightclient

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"gotest.tools/assert"
	"testing"
)

func TestMsgTxFromHex(t *testing.T) {
	txID := "f82a618b2f6212f2134d2dc66aad26c0aa0a2ed7fa53a298c53736d0510ad69f"
	txData := "01000000014260bfba493220006a8141b162e5785c863149c043d7abf0030276c10d2a72e5010000008b483045022100d1424374bd4e6a0264dd55baf7de69a17bbe2416b20a55a1475fd3779799d13b02202e9249feb3b43cd82c8f097a1ffb5f98b06ca27595697bbc2acee4a5c255fd830141045408a52d4b3cdc9c78c14418c38d4a0fd6e0ed396f966d3f51c86164351bbd0e38426b922c9c1f79bca05dfeb72e08ba9d41eb5972e48db217301a9bf795aeadffffffff020065cd1d000000001976a914de526c004ab53c6d7ed3ff98789bca013797d33f88acf0ca8744000000001976a9148fd0e06ce84a4f108b1e259f5aeef21a16a2272c88ac00000000"
	tx, err := MsgTxFromHex(txData)
	assert.NilError(t, err)
	assert.Equal(t, tx.TxID(), txID)
}

func TestVerifyBalance(t *testing.T) {
	type testCase struct {
		name          string
		gettxoutproof string
		txData        string
		txHash        string
		addr          string
		balance       BalanceReport
	}

	run := func(t *testing.T, tc testCase) {
		header, err := BlockHeaderFromHex(tc.gettxoutproof[:160])
		assert.NilError(t, err)

		lc := NewBTCLightClientWithData(&chaincfg.MainNetParams, []wire.BlockHeader{header}, 1000)

		spvProof, err := SPVProofFromHex(tc.gettxoutproof, tc.txHash)
		assert.NilError(t, err)

		tx, _ := MsgTxFromHex(tc.txData)

		balance, err := lc.VerifyBalance(tx, tc.addr, *spvProof)
		assert.NilError(t, err)
		assert.Equal(t, balance, tc.balance)
	}

	testCases := []testCase{
		{
			name:          "Happy test case",
			gettxoutproof: "000000202834abd71bdd0d3298542af4506918ea168ce002936b040000000000000000001da8757e4d756e848245cacf3e103c1b9f6ed2405c6d818a73172c8ec72856d4db3864606fdf0c17dcc1000ccb0300000cb86343fc64abcdab51e530303a4ee2b420fa6b5a12b435c9c76fe953ca5471ca074a0bfaf4462cef0a5665b89fd7fd5e4f8536630cde6824d09b20400b2f65eed9f744b2dc695b0ea0c4afd06310a21b93ddd7270a781acd0ada1afdd23b5750aa59aac6bcb5a037cbc56b9efbfc159a36142a07d23e81c4b89d3dbbc31be1cefe0bb7b0369ffc3b1d530e234987543a2613bbb8b06c86f993a930dee7b9d87f661ef556adc0174c7f180aa28006ee93ce2291302801ecd045c234c00b186ea35ff1e77eac3f113492e2eb12f38b9df452f5831f55c861865ac8f3c7dd06be2377f859ba1d12dea2ec44987796a27d42d5727250c1e0181d6a251f8272f21b9a2034069a2471de43de655619904d43b4665f6ce38741320998dc97838c32c79f1ada066ddf7a441357d55cc42a8906970bff2d5342be694002476733ff593af26f320c10df7ba9a76355438f462c040b598868dfb67c5e88d6d9a426ec8cdd74337d42df6b29e9fb319410848f3ff7228d00dc539e2962d185348ab9663a112a03ff6e00",
			txHash:        "ee652f0b40209bd02468de0c6336854f5efdd79fb865560aef2c46f4fa0b4a07",
			txData:        "010000000001019dcfb29f9a915612fc334232f5bc2b03d710feb85ad4a39e08dc817a426712750100000000ffffffff0270c286000000000017a914a9c923dfcc27c61c114f2170492ccf1155a4484487487b850000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d040047304402202d211d83acb7cb9dc3c45f86960e39fb195aee33f53d9f2bb7a591b0b328acdc02207734781019fa0821773472e886fc4b73255b989ab3ea3c0c5561af7d7855e301014730440220205f9f2da387e2de609e2f7aeeab106930113cf273993bfe09f392255b8ddecd02203fbc7c7ba2ed262868783528a8e76784ac2a930eb54469a028316791c62d6240016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000",
			addr:          "3HAm5RaZPeMWXH3qoAdF5oidXMrg2XCWow",
			balance:       newValidBalance(8831600),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
