package btclightclient

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
		index       uint
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

			index:    0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
