package btclightclient

import (
	"testing"

	"gotest.tools/assert"
)

func TestPartialMerkleTree(t *testing.T) {
	type testCase struct {
		name                 string
		txoutproof           string
		txID                 string
		expectedErrorMessage string
		root                 string
	}

	run := func(t *testing.T, tc testCase) {
		pmt, err := PartialMerkleTreeFromHex(tc.txoutproof[160:])
		if err != nil {
			assert.Error(t, err, tc.expectedErrorMessage)
			return
		}

		merkleProof, err := pmt.GetProof(tc.txID)

		if err != nil {
			assert.Error(t, err, tc.expectedErrorMessage)
			return
		}

		if len(tc.expectedErrorMessage) > 0 {
			t.Fatalf("Should return error")
		}

		assert.Equal(t, merkleProof.merkleRoot.String(), tc.root)
	}

	testCases := []testCase{
		{
			name: "Happy test case",
			// data from chainstack example: https://docs.chainstack.com/reference/bitcoin-gettxoutproof
			txoutproof:           "00e0002000471175ec71a72541c100f21bb79f9da0e5ca98259a000000000000000000004769eae15b51056127304c5dec6d94c7840f8f922c0b65bc32177cb46ce05de9b8c10866d36203175dae051fdc0a00000d625aa7b5510f7c003624338259d21544e61ccb3666792dde9734b7621d2cf80bb81ffa45657310bdc47ad3b3f5e5346c150d4fc1b98a5446cc560c6f38f7156138761aab058be861e51fe52ea7cf7b4914a1e1b159ecebe46b51db0ec5cfd4c2324ae9c132169d1f133981632895c216a8e3c3d3a9cea545fade4c0ab8b626a2791862728b657abbdb06dedcc3faabee9d72ce6b8252b45fc99d6fe0f79cec401e1a431774d8830b962e5dee97fc96f4f85f84a6e50b986a37b35318537a81f3f8c604554e5b4f5ca4b4437caa3b0723896396532c1985d52f42f915084534c6bedb4ded1238781d23be0173b94ca25d7faff2832ac99fa16b2f9b219ff276062f100d4b7ce774ba405fbad36b65165e2e5aece3e0b9718886d7b24708be5ae72d10911e9301811b19fcb218ce7dfee31729f4ef56a3d8f31670865a039b3678b34fcb47f12bd157a064339c3e91a960c5a14b9e8da9c8ce211a02bb94e7165a1668d8a17663e95adcdacdbc8e8ab793e8796fda9b270ca957e67aa33dc95cff158cb2ff6882064942ef545612a8eceb3c60415d677d170f4351ede1f7a8807504ff1f0000",
			txID:                 "0bf82c1d62b73497de2d796636cb1ce64415d25982332436007c0f51b5a75a62",
			expectedErrorMessage: "",
			root:                 "e95de06cb47c1732bc650b2c928f0f84c7946dec5d4c30276105515be1ea6947",
		},
		{
			name: "Error b/c txID not exist in merkle tree",
			// data from chainstack example: https://docs.chainstack.com/reference/bitcoin-gettxoutproof
			txoutproof:           "00e0002000471175ec71a72541c100f21bb79f9da0e5ca98259a000000000000000000004769eae15b51056127304c5dec6d94c7840f8f922c0b65bc32177cb46ce05de9b8c10866d36203175dae051fdc0a00000d625aa7b5510f7c003624338259d21544e61ccb3666792dde9734b7621d2cf80bb81ffa45657310bdc47ad3b3f5e5346c150d4fc1b98a5446cc560c6f38f7156138761aab058be861e51fe52ea7cf7b4914a1e1b159ecebe46b51db0ec5cfd4c2324ae9c132169d1f133981632895c216a8e3c3d3a9cea545fade4c0ab8b626a2791862728b657abbdb06dedcc3faabee9d72ce6b8252b45fc99d6fe0f79cec401e1a431774d8830b962e5dee97fc96f4f85f84a6e50b986a37b35318537a81f3f8c604554e5b4f5ca4b4437caa3b0723896396532c1985d52f42f915084534c6bedb4ded1238781d23be0173b94ca25d7faff2832ac99fa16b2f9b219ff276062f100d4b7ce774ba405fbad36b65165e2e5aece3e0b9718886d7b24708be5ae72d10911e9301811b19fcb218ce7dfee31729f4ef56a3d8f31670865a039b3678b34fcb47f12bd157a064339c3e91a960c5a14b9e8da9c8ce211a02bb94e7165a1668d8a17663e95adcdacdbc8e8ab793e8796fda9b270ca957e67aa33dc95cff158cb2ff6882064942ef545612a8eceb3c60415d677d170f4351ede1f7a8807504ff1f0000",
			txID:                 "e95de06cb47c1732bc650b2c928f0f84c7946dec5d4c30276105515be1ea6947",
			expectedErrorMessage: "node value doesn't exist in merkle tree",
			root:                 "",
		},
		{
			name:                 "Can't decode txoutproof",
			txoutproof:           "00e0002000471175ec71a72541c100f21bb79f9da0e5ca98259a000000000000000000004769eae15b51056127304c5edec6d94c7840f8f922c0b65bc32177cb46ce05de9b8c10866d36203175dae051fdc0a00000d625aa7b5510f7c003624338259d21544e61ccb3666792dde9734b7621d2cf80bb81ffa45657310bdc47ad3b3f5e5346c150d4fc1b98a5446cc560c6f38f7156138761aab058be861e51fe52ea7cf7b4914a1e1b159ecebe46b51db0ec5cfd4c2324ae9c132169d1f133981632895c216a8e3c3d3a9cea545fade4c0ab8b626a2791862728b657abbdb06dedcc3faabee9d72ce6b8252b45fc99d6fe0f79cec401e1a431774d8830b962e5dee97fc96f4f85f84a6e50b986a37b35318537a81f3f8c604554e5b4f5ca4b4437caa3b0723896396532c1985d52f42f915084534c6bedb4ded1238781d23be0173b94ca25d7faff2832ac99fa16b2f9b219ff276062f100d4b7ce774ba405fbad36b65165e2e5aece3e0b9718886d7b24708be5ae72d10911e9301811b19fcb218ce7dfee31729f4ef56a3d8f31670865a039b3678b34fcb47f12bd157a064339c3e91a960c5a14b9e8da9c8ce211a02bb94e7165a1668d8a17663e95adcdacdbc8e8ab793e8796fda9b270ca957e67aa33dc95cff158cb2ff6882064942ef545612a8eceb3c60415d677d170f4351ede1f7a8807504ff1f0000e",
			txID:                 "",
			expectedErrorMessage: "out-bound of vHash",
			root:                 "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}

}
