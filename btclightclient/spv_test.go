package btclightclient

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"gotest.tools/assert"
)

func hashFromString(t *testing.T, hashStr string) chainhash.Hash {
	h, err := chainhash.NewHashFromStr(hashStr)
	assert.NilError(t, err)
	return *h
}

func TestSPVMerkleRoot(t *testing.T) {
	type testCase struct {
		name        string
		txIdHash    string
		merkleRoot  string
		merkleProof []string
		index       uint32
	}

	run := func(t *testing.T, tc testCase) {
		merklePath := make([]chainhash.Hash, len(tc.merkleProof))

		for i := 0; i < len(merklePath); i++ {
			merklePath[i] = hashFromString(t, tc.merkleProof[i])
		}

		spv := SPVProof{
			blockHash:  hashFromString(t, "9d5369290c3d97f47d07a12e1b7f171c9c69ddfc876ecec6a16dfe3e94773a1c"), // dummy hash for testing only
			txId:       tc.txIdHash,
			txIndex:    tc.index,
			merklePath: merklePath,
		}

		actual := spv.MerkleRoot()
		assert.Assert(t, actual.String() == tc.merkleRoot)
	}

	testCases := []testCase{
		{
			// https://docs.chainstack.com/reference/bitcoin-gettxoutproof
			name:     "Success Verify TX ID",
			txIdHash: "4224625b409323d17e8842f935ce3764c3e7203ad0de3d403558881089cb3632",
			merkleProof: []string{
				"4224625b409323d17e8842f935ce3764c3e7203ad0de3d403558881089cb3632",
				"8e2f5bf7ad1839ec83538d0c7a8bc4f0f9b624896760cccc2bfe399746c8277f",
				"2feb97ce13d350cab8c94aa58772ad825486b1b39112b2189da7294fd8740085",
				"60ba6889cf8183fd353bd87248bbeaa46c17483d6a34d4914c6cb709fb5497bc",
				"05f936db24cd79d0629346a4117a10c33efc4b0953311f925e257cc94ff1ad0c",
				"58a27cc4eb3bee432e0ab49284f5f01a2657f300031178f44cd2349d3ba72257",
				"e95ca580bef27a940c638864ac39cdf9fe85b6f79658484fbf8add73585f6dc4",
				"849519e837153d349de2541689bf7255af169e9aac193af985abaa76ef620a64",
				"cd1c4dbcc6d6a62e68054ff760b175c9e80361e6ae90dd0a3be8815b0a8c11cf",
			},
			merkleRoot: "e4dc1f2ab5ac974d6b23690bd4d8ddbde63e9647c09a3d8f6b77fc0bf53544e5",
			index:      0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestSPVVerification(t *testing.T) {
	type testCase struct {
		name          string
		gettxoutproof string
		txHash        string
		spvStatus     SPVStatus
	}

	run := func(t *testing.T, tc testCase) {
		header, err := BlockHeaderFromHex(tc.gettxoutproof[:160])
		assert.NilError(t, err)

		lc := NewBTCLightClientWithData(&chaincfg.MainNetParams, []wire.BlockHeader{header}, 1000)
		spvProof, err := SPVProofFromHex(tc.gettxoutproof, tc.txHash)
		assert.NilError(t, err)
		verifyStatus := lc.VerifySPV(*spvProof)
		assert.Equal(t, verifyStatus, tc.spvStatus)
	}

	testCases := []testCase{
		// https://docs.chainstack.com/reference/bitcoin-gettxoutproof
		// transactionIndex is even
		{
			name:          "Valid SPV Test Case  even transactionIndex",
			gettxoutproof: "00e0002000471175ec71a72541c100f21bb79f9da0e5ca98259a000000000000000000004769eae15b51056127304c5dec6d94c7840f8f922c0b65bc32177cb46ce05de9b8c10866d36203175dae051fdc0a00000d625aa7b5510f7c003624338259d21544e61ccb3666792dde9734b7621d2cf80bb81ffa45657310bdc47ad3b3f5e5346c150d4fc1b98a5446cc560c6f38f7156138761aab058be861e51fe52ea7cf7b4914a1e1b159ecebe46b51db0ec5cfd4c2324ae9c132169d1f133981632895c216a8e3c3d3a9cea545fade4c0ab8b626a2791862728b657abbdb06dedcc3faabee9d72ce6b8252b45fc99d6fe0f79cec401e1a431774d8830b962e5dee97fc96f4f85f84a6e50b986a37b35318537a81f3f8c604554e5b4f5ca4b4437caa3b0723896396532c1985d52f42f915084534c6bedb4ded1238781d23be0173b94ca25d7faff2832ac99fa16b2f9b219ff276062f100d4b7ce774ba405fbad36b65165e2e5aece3e0b9718886d7b24708be5ae72d10911e9301811b19fcb218ce7dfee31729f4ef56a3d8f31670865a039b3678b34fcb47f12bd157a064339c3e91a960c5a14b9e8da9c8ce211a02bb94e7165a1668d8a17663e95adcdacdbc8e8ab793e8796fda9b270ca957e67aa33dc95cff158cb2ff6882064942ef545612a8eceb3c60415d677d170f4351ede1f7a8807504ff1f0000",
			txHash:        "0bf82c1d62b73497de2d796636cb1ce64415d25982332436007c0f51b5a75a62",
			spvStatus:     ValidSPVProof,
		},
		// transactionIndex is odd
		{
			name:          "Valid SPV Test Case  odd transactionIndex",
			gettxoutproof: "000000202834abd71bdd0d3298542af4506918ea168ce002936b040000000000000000001da8757e4d756e848245cacf3e103c1b9f6ed2405c6d818a73172c8ec72856d4db3864606fdf0c17dcc1000ccb0300000cb86343fc64abcdab51e530303a4ee2b420fa6b5a12b435c9c76fe953ca5471ca074a0bfaf4462cef0a5665b89fd7fd5e4f8536630cde6824d09b20400b2f65eed9f744b2dc695b0ea0c4afd06310a21b93ddd7270a781acd0ada1afdd23b5750aa59aac6bcb5a037cbc56b9efbfc159a36142a07d23e81c4b89d3dbbc31be1cefe0bb7b0369ffc3b1d530e234987543a2613bbb8b06c86f993a930dee7b9d87f661ef556adc0174c7f180aa28006ee93ce2291302801ecd045c234c00b186ea35ff1e77eac3f113492e2eb12f38b9df452f5831f55c861865ac8f3c7dd06be2377f859ba1d12dea2ec44987796a27d42d5727250c1e0181d6a251f8272f21b9a2034069a2471de43de655619904d43b4665f6ce38741320998dc97838c32c79f1ada066ddf7a441357d55cc42a8906970bff2d5342be694002476733ff593af26f320c10df7ba9a76355438f462c040b598868dfb67c5e88d6d9a426ec8cdd74337d42df6b29e9fb319410848f3ff7228d00dc539e2962d185348ab9663a112a03ff6e00",
			txHash:        "ee652f0b40209bd02468de0c6336854f5efdd79fb865560aef2c46f4fa0b4a07",
			spvStatus:     ValidSPVProof,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
