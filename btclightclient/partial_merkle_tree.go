package btclightclient

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// We verify a partial merkle tree proof. This logic is used in verifytxoutproof

type PartialMerkleTree struct {
	numberTransactions uint
	vBits              []bool
	vHash              []*chainhash.Hash
}

type MerkleProof struct {
	nodeValue  chainhash.Hash
	merklePath []chainhash.Hash
	pos        uint32
}

// TODO: verify this value
const maxAllowBytes = 65536

// TODO: wrap error
func readPartialMerkleTree(r io.Reader, buf []byte) (*PartialMerkleTree, error) {
	if _, err := io.ReadFull(r, buf[:4]); err != nil {
		return nil, err
	}

	numberTransactions := uint(binary.LittleEndian.Uint32(buf[:4]))

	var pver uint32 // pversion but btcd don't use this in those function we want.

	var vHash []*chainhash.Hash
	if numberOfHashes, err := wire.ReadVarIntBuf(r, pver, buf); err != nil {
		return nil, err
	} else {
		if numberHash*32 > maxAllowBytes {
			return nil, errors.New("number of hashes is too big")
		}

		bytes := make([]byte, numberHash*32)
		if _, err := io.ReadFull(r, bytes); err != nil {
			return nil, err
		}

		vHash = make([]*chainhash.Hash, numberHash)
		for i := 0; i < numberHash; i++ {
			vHash[i], err = chainhash.NewHash(bytes[i*32 : (i+1)*32])
			if err != nil {
				return nil, err
			}
		}
	}

	var vBits []bool
	if vBytes, err := wire.ReadVarBytes(r, pver, uint32(maxAllowBytes), "vBits"); err != nil {
		return nil, err
	} else {

		vBits = make([]bool, len(vBytes)*8)
		i := 0
		for _, b := range vBytes {
			for j := 0; j < 8; j++ {
				vBits[i] = (b & (1 << j)) != 0
				i++
			}
		}
	}

	partialMerkleTree := PartialMerkleTree{
		numberTransactions: numberTransactions,
		vBits:              vBits,
		vHash:              vHash,
	}

	return &parialMerkleTree, nil
}

func ParialMerkleTreeFromHex(hexStr string) (PartialMerkleTree, error) {
	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(hexBytes)
	return readPartialMerkleTree(reader, hexBytes)
}

func (pmk *PartialMerkleTree) CalcTreeWidth(height uint32) uint {
	return (pmk.numberTransactions + (1 << height) - 1) >> height
}

// MinHeight returns the minimum height of a Merkele tree to fit `pmk.numberTransactions`.
func (pmk *PartialMerkleTree) Height() uint32 {
	var nHeight uint32 = 0
	for pmk.CalcTreeWidth(nHeight) > 1 {
		nHeight++
	}
	return nHeight
}

// Port logic from gettxoutproof from bitcoin-core
// TODO: Make error handler more sense
func (pmk *PartialMerkleTree) computeMerkleProofRecursive(height, pos uint32, nBitUsed, nHashUsed *uint32, txID *chainhash.Hash) (*MerkleProof, error) {
	if int(*nBitUsed) >= len(pmk.vBits) {
		return nil, errors.New("Error")
	}

	fParentOfMatch := pmk.vBits[*nBitUsed]
	*nBitUsed = *nBitUsed + 1

	if height == 0 || !fParentOfMatch {
		if int(*nHashUsed) >= len(pmk.vHash) {
			return nil, errors.New("error")
		}

		hash := pmk.vHash[*nHashUsed]
		*nHashUsed++
		if height == 0 && fParentOfMatch {
			return &MerkleProof{
				nodeValue:  *hash,
				merklePath: []chainhash.Hash{*hash},
				pos:        pos,
			}, nil
		}

		return &MerkleProof{
			nodeValue: *hash,
			merklePath: []chainhash.Hash{},
			pos: uint32(pmk.numberTransactions) + 1,
		}, nil
		
	} else {
		left, err := pmk.computeMerkleProofRecursive(height-1, pos*2, nBitUsed, nHashUsed, txID)
		var right *MerkleProof
		if err != nil {
			return nil, err
		}
		if uint(pos*2+1) < pmk.CalcTreeWidth(height-1) {
			right, err = pmk.computeMerkleProofRecursive(height-1, pos*2+1, nBitUsed, nHashUsed, txID)
			if err != nil {
				return nil, err
			}

			if left.nodeValue.IsEqual(&right.nodeValue) {
				return nil, errors.New("error")
			}
		} else {
			right = left
		}

		nodeValue := Hash256MerkleStepHashChain(&left.nodeValue, &right.nodeValue)
		// Compute new proof
		if left.pos != uint32(pmk.numberTransactions) {
			// txID on the left side
			return &MerkleProof {
				nodeValue : *nodeValue,
				pos: left.pos,
				merklePath : append(left.merklePath, right.nodeValue),
			}, nil
		}
		// TxID on right side
		return &MerkleProof {
				nodeValue : *nodeValue,
				pos: right.pos,
				merklePath : append(right.merklePath, left.nodeValue),
		}, nil 
	}
}

func (pmk *PartialMerkleTree) ComputeMerkleProof(txID string) (*MerkleProof, error) {
	txIDHash, err := chainhash.NewHashFromStr(txID);
	if err != nil {
		return nil, err
	}
	height := pmk.Height();
	nUsedBit := uint32(0)
	nUsedHash := uint32(0)	
	return pmk.computeMerkleProofRecursive(height, 0, &nUsedBit, &nUsedHash, txIDHash)
}

// Hash256MerkleStep concatenates and hashes two inputs for merkle proving
func Hash256MerkleStep(a, b []byte) chainhash.Hash {
	c := make([]byte{}, 0, len(a) + len(b))
	c = append(c, a...)
	c = append(c, b...)
	return  chainhash.DoubleHashH(c)
}

// TODO: make it more simple
func HashNodes(a, b *chainhash.Hash) *chainhash.Hash {
	c := make([]byte{}, 0, chainhash.HashSize*2)
	c = append(c, a...)
	c = append(c, b...)
	return &chainhash.DoubleHashH(c)
}
