package btclightclient

import (
	"errors"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// SPV proof. We use this for verify transaction inclusives in block.
// We are verify for single transaction in this version.
type SPVProof struct {
	blockHash  chainhash.Hash
	txId       string // 32bytes hash value in string hex format
	txIndex    uint32 // index of transaction in block
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

// Get SPV proof from gettxoutproof Bitcoin API.
func SPVProofFromHex(txoutProof string, txID string) (*SPVProof, error) {
	blockheader, err := BlockHeaderFromHex(txoutProof[:160])
	if err != nil {
		return nil, err
	}

	pmt, err := PartialMerkleTreeFromHex(txoutProof[160:])
	if err != nil {
		return nil, err
	}
	merkleProof, err := pmt.GetProof(txID)
	if err != nil {
		return nil, err
	}

	return &SPVProof{
		blockHash:  blockheader.BlockHash(),
		txId:       txID,
		txIndex:    merkleProof.transactionIndex,
		merklePath: merkleProof.merklePath,
	}, nil
}

func (spvProof SPVProof) MerkleRoot() (chainhash.Hash, error) {
	var currentHash chainhash.Hash
	if len(spvProof.merklePath) == 0 {
		return currentHash, errors.New("merkle path can't empty")
	}

	currentHash = spvProof.merklePath[0]
	numberSteps := len(spvProof.merklePath)
	transactionIndex := spvProof.txIndex
	for i := 1; i < numberSteps; i++ {
		if transactionIndex%2 == 0 {
			currentHash = *HashNodes(&currentHash, &spvProof.merklePath[i])
		} else {
			currentHash = *HashNodes(&spvProof.merklePath[i], &currentHash)
		}
		transactionIndex /= 2
	}
	return currentHash, nil
}

func (lc *BTCLightClient) VerifySPV(spvProof SPVProof) SPVStatus {
	lightBlock := lc.btcStore.LightBlockByHash(spvProof.blockHash)

	// In the case light block not belong currect database
	if lightBlock == nil {
		return InvalidSPVProof
	}

	if spvProof.txId != spvProof.merklePath[0].String() {
		return InvalidSPVProof
	}

	blockMerkleRoot := lightBlock.Header.MerkleRoot
	spvMerkleRoot, err := spvProof.MerkleRoot()
	if err != nil {
		return InvalidSPVProof
	}
	if !spvMerkleRoot.IsEqual(&blockMerkleRoot) {
		return InvalidSPVProof
	}

	// in the case the block not finalize
	if lc.btcStore.LatestFinalizedHeight() < int64(lightBlock.Height) {
		return PartialValidSPVProof
	}

	return ValidSPVProof
}
