package btclightclient

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"slices"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"gotest.tools/assert"
)

func encodeTxID(t *testing.T, txID string) []byte {
	b, err := hex.DecodeString(txID)
	assert.NilError(t, err)
	slices.Reverse(b)
	return b
}


// ExtractTxIndexLE extracts the LE tx input index from the input in a tx
// Returns the tx index as a little endian []byte
func ExtractTxIndexLE(input []byte) []byte {
	return input[32:36:36]
}

// ReverseEndianness takes in a byte slice and returns a
// reversed endian byte slice.
func ReverseEndianness(b []byte) []byte {
	out := make([]byte, len(b), len(b))
	copy(out, b)

	for i := len(out)/2 - 1; i >= 0; i-- {
		opp := len(out) - 1 - i
		out[i], out[opp] = out[opp], out[i]
	}

	return out
}


// BytesToUint converts 1, 2, 3, or 4-byte numbers to uints
func BytesToUint(b []byte) uint {
	total := uint(0)
	length := uint(len(b))

	for i := uint(0); i < length; i++ {
		total += uint(b[i]) << ((length - i - 1) * 8)
	}

	return total
}


// ExtractTxIndex extracts the tx input index from the input in a tx
func ExtractTxIndex(input []byte) uint {
	return BytesToUint(ReverseEndianness(ExtractTxIndexLE(input)))
}


func TestVerifyTxOutProof(t *testing.T) {
	
	// data from chainstack example: https://docs.chainstack.com/reference/bitcoin-gettxoutproof
	hexStr := "00e0002000471175ec71a72541c100f21bb79f9da0e5ca98259a000000000000000000004769eae15b51056127304c5dec6d94c7840f8f922c0b65bc32177cb46ce05de9b8c10866d36203175dae051fdc0a00000d625aa7b5510f7c003624338259d21544e61ccb3666792dde9734b7621d2cf80bb81ffa45657310bdc47ad3b3f5e5346c150d4fc1b98a5446cc560c6f38f7156138761aab058be861e51fe52ea7cf7b4914a1e1b159ecebe46b51db0ec5cfd4c2324ae9c132169d1f133981632895c216a8e3c3d3a9cea545fade4c0ab8b626a2791862728b657abbdb06dedcc3faabee9d72ce6b8252b45fc99d6fe0f79cec401e1a431774d8830b962e5dee97fc96f4f85f84a6e50b986a37b35318537a81f3f8c604554e5b4f5ca4b4437caa3b0723896396532c1985d52f42f915084534c6bedb4ded1238781d23be0173b94ca25d7faff2832ac99fa16b2f9b219ff276062f100d4b7ce774ba405fbad36b65165e2e5aece3e0b9718886d7b24708be5ae72d10911e9301811b19fcb218ce7dfee31729f4ef56a3d8f31670865a039b3678b34fcb47f12bd157a064339c3e91a960c5a14b9e8da9c8ce211a02bb94e7165a1668d8a17663e95adcdacdbc8e8ab793e8796fda9b270ca957e67aa33dc95cff158cb2ff6882064942ef545612a8eceb3c60415d677d170f4351ede1f7a8807504ff1f0000";
	proof, _ := hex.DecodeString(hexStr[160:]);
	fmt.Println(len(hexStr));

	var pmk PMerkleTree
	reader := bytes.NewReader(proof)

	err := readPMerkleTree(reader, &pmk, proof)

	assert.NilError(t, err)


	height := pmk.Height();
	nUsedBit := uint32(0);
	nUsedHash := uint32(0);
	vMatch := make([]*chainhash.Hash, 0);
	vnIndex := make([]uint32, 0);
	root, err := pmk.ComputerRootPMerkleTree(height, 0, &nUsedBit, &nUsedHash, &vMatch, &vnIndex)

	assert.NilError(t, err)
	assert.Equal(t, root.String(), "e95de06cb47c1732bc650b2c928f0f84c7946dec5d4c30276105515be1ea6947");
}


func TestUTXOVerification(t *testing.T) {

	type testCase struct {
		name        string
		txIdHash    string
		merkleRoot  string
		merkleProof []string
		index       uint
		expected    bool
	}

	run := func(t *testing.T, tc testCase) {
		proof := encodeTxID(t, tc.txIdHash)
		for _, n := range tc.merkleProof {
			proof = append(proof, encodeTxID(t, n)...)
		}
		proof = append(proof, encodeTxID(t, tc.merkleRoot)...)

		actual := VerifyHash256Merkle(proof, tc.index)
		assert.Assert(t, actual == tc.expected)
	}

	testCases := []testCase{
		{
			name:     "Success Verify TX ID",
			txIdHash: "4224625b409323d17e8842f935ce3764c3e7203ad0de3d403558881089cb3632",
			merkleProof: []string{
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
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
