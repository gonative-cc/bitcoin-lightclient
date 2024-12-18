package btclightclient

import (
	"encoding/hex"
	"fmt"
	"slices"
	"testing"

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


func TestSPVFromHex(t *testing.T) {
	hexStr := "00000030516567e505288fe41b2fc6be9b96318c406418c7d338168fe75a26111490eb2fec401c3902aa39842e53a0c641af518957ec3aa5984a44d32e2a9f7fee2fa67a3f5b6167ffff7f20040000000100000001ec401c3902aa39842e53a0c641af518957ec3aa5984a44d32e2a9f7fee2fa67a0101";
	spv, err := SPVFromHex(hexStr)

	assert.NilError(t, err)
	fmt.Println(spv.blockHash)
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
