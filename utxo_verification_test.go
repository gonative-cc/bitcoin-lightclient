package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)


func BigToLitle(big []byte) []byte {
	little := make([]byte, len(big))
	for i := 0; i < len(big); i += 4 {
		end := i + 4
		if end > len(big) {
			end = len(big)
		}
		for j := 0; j < end-i; j++ {
			little[i+j] = little[end-1-j]
		}
	}
	return little
	
}

func TestUTXOVerification(t *testing.T) {
	merkleProofStrings := []string{
		"8e2f5bf7ad1839ec83538d0c7a8bc4f0f9b624896760cccc2bfe399746c8277f",
		"2feb97ce13d350cab8c94aa58772ad825486b1b39112b2189da7294fd8740085",
		"60ba6889cf8183fd353bd87248bbeaa46c17483d6a34d4914c6cb709fb5497bc",
		"05f936db24cd79d0629346a4117a10c33efc4b0953311f925e257cc94ff1ad0c",
		"58a27cc4eb3bee432e0ab49284f5f01a2657f300031178f44cd2349d3ba72257",
		"e95ca580bef27a940c638864ac39cdf9fe85b6f79658484fbf8add73585f6dc4",
		"849519e837153d349de2541689bf7255af169e9aac193af985abaa76ef620a64",
		"cd1c4dbcc6d6a62e68054ff760b175c9e80361e6ae90dd0a3be8815b0a8c11cf",
		"e4dc1f2ab5ac974d6b23690bd4d8ddbde63e9647c09a3d8f6b77fc0bf53544e5",
	}
	tx, _ := hex.DecodeString("4224625b409323d17e8842f935ce3764c3e7203ad0de3d403558881089cb3632")
	tx = BigToLitle(tx)
	proof := []byte{}

	
	proof = append(proof, tx...)
	
	for _, element := range merkleProofStrings {
		b, _ := hex.DecodeString(element)
		b = BigToLitle(b)
		proof = append(proof, b...)
	}


	fmt.Println(len(proof))
	ans := VerifyHash256Merkle(proof, 0)
	if !ans {
		t.Fatal("invalid")
	}
}
