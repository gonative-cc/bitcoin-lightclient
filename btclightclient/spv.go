package btclightclient

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type SPVProof struct {
	blockHash  chainhash.Hash
	txId       string
	txIndex    uint
	merklePath []chainhash.Hash
}

type SPVStatus int

const (
	// proof is complete wrong
	InvalidSPVProof SPVStatus = iota
	// proof valid but the block have not "finalized" yet.
	PartialValidSPVProof
	// proof vlaid and block is "finalized".
	ValidSPVProof
)

// / Get SPV proof from gettxoutproof Bitcoin API.
func SPVProofFromHex(txoutProof string, txID string) (*SPVProof, error) {
	blockheader, err := BlockHeaderFromHex(txoutProof[:160])
	if err != nil {
		return nil, err
	}

	pmt, err := PartialMerkleTreeFromHex(txoutProof[160:])
	if err != nil {
		return nil, err; 
	}
	merkleProof, err := pmt.GetProof(txID)	
	if err != nil {
		return nil, err
	}

	return &SPVProof{
		blockHash:  blockheader.BlockHash(),
		txId:       txID,
		txIndex:    uint(merkleProof.transactionIndex),
		merklePath: merkleProof.merklePath,
	}, nil
}

func (spvProof SPVProof) MerkleRoot() chainhash.Hash {
	hashValue := &spvProof.merklePath[0]
	numberSteps := len(spvProof.merklePath)
	transactionIndex := spvProof.txIndex
	for i := 1; i < numberSteps; i++ {
		if transactionIndex % 2 == 0 {
			hashValue = HashNodes(hashValue, &spvProof.merklePath[i])
		} else {
			hashValue = HashNodes(&spvProof.merklePath[i], hashValue);
		}
		transactionIndex /= 2
	}
	return *hashValue
}

func (lc *BTCLightClient) VerifySPV(spvProof SPVProof) SPVStatus {
	lightBlock := lc.btcStore.LightBlockByHash(spvProof.blockHash)

	// In the case light block not belong currect database
	if lightBlock == nil {
		return InvalidSPVProof
	}

	blockMerkleRoot := lightBlock.Header.MerkleRoot
	spvMerkleRoot := spvProof.MerkleRoot()

	if !spvMerkleRoot.IsEqual(&blockMerkleRoot) {
		return InvalidSPVProof
	}

	// in the case the block not finalize
	if lc.btcStore.LatestFinalizedHeight() < int64(lightBlock.Height) {
		return PartialValidSPVProof
	}

	return ValidSPVProof
}
