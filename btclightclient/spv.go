package btclightclient

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// SPV proof. We use this for verify transaction inclusives in block.
// We are verify for single transaction in this version.
type SPVProof struct {
	BlockHash  chainhash.Hash
	TxId       string // 32bytes hash value in string hex format
	TxIndex    uint32 // index of transaction in block
	MerklePath []chainhash.Hash
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
		BlockHash:  blockheader.BlockHash(),
		TxId:       txID,
		TxIndex:    merkleProof.transactionIndex,
		MerklePath: merkleProof.merklePath,
	}, nil
}

func (spvProof SPVProof) MerkleRoot() chainhash.Hash {
	hashValue := &spvProof.MerklePath[0]
	numberSteps := len(spvProof.MerklePath)
	transactionIndex := spvProof.TxIndex
	for i := 1; i < numberSteps; i++ {
		if transactionIndex%2 == 0 {
			hashValue = HashNodes(hashValue, &spvProof.MerklePath[i])
		} else {
			hashValue = HashNodes(&spvProof.MerklePath[i], hashValue)
		}
		transactionIndex /= 2
	}
	return *hashValue
}

func (lc *BTCLightClient) VerifySPV(spvProof SPVProof) SPVStatus {
	lightBlock := lc.btcStore.LightBlockByHash(spvProof.BlockHash)

	// light block not belong currect database
	if lightBlock == nil {
		return InvalidSPVProof
	}

	if len(spvProof.MerklePath) == 0 {
		return InvalidSPVProof
	}

	if spvProof.TxId != spvProof.MerklePath[0].String() {
		return InvalidSPVProof
	}

	blockMerkleRoot := lightBlock.Header.MerkleRoot
	spvMerkleRoot := spvProof.MerkleRoot()

	if !spvMerkleRoot.IsEqual(&blockMerkleRoot) {
		return InvalidSPVProof
	}

	// the block not finalize
	if lc.btcStore.LatestFinalizedHeight() < int64(lightBlock.Height) {
		return PartialValidSPVProof
	}

	return ValidSPVProof
}

func (lc *BTCLightClient) VerifySPVs(spvProofs []SPVProof) []SPVStatus {
	result := make([]SPVStatus, len(spvProofs))
	for i, spv := range spvProofs {
		result[i] = lc.VerifySPV(spv)
	}
	return result
}
