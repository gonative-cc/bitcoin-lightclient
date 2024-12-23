package btclightclient

import (
	"bytes"
	"encoding/hex"
)

/// Get SPV proof from gettxoutproof Bitcoin API.
func SPVProofFromHex(proofHex string, txID string) (*SPVProof, error){
	// get block header
	blockheader, err := BlockHeaderFromHex(proofHex[:160]);
	if err != nil {
		return nil, err;
	}

	// get merkle proof for txID
	merkleProofBytes, _ := hex.DecodeString(proofHex[160:]);
	reader := bytes.NewReader(merkleProofBytes)
	pmk, err := readPartialMerkleTree(reader, merkleProofBytes)
	if err != nil {
		return nil, err
 	}
	
	merkleProof, err := pmk.ComputeMerkleProof(txID);
	if err != nil {
		return nil, err
	}

	return &SPVProof {
		blockHash: blockheader.BlockHash(),
			txId: txID,
			txIndex: uint(merkleProof.pos),
			merkleProof: *merkleProof.merklePath[:],
	}, nil
}
