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

// We verify a partial merkle tree proof. This logic is used in verifytxoutproof.
// First one, we must know how merkle tree store in bitcoin. The tree is prefect
// binary tree. At height `h` this index from 0 to `tree width` at level h.
// At level(height) 0, this is the list of transaction IDs. At level 1, we compute
// node value at positon `p` = hash(node at (h - 1, p * 2), node at (h - 1, p * 2 + 1))
// You can see p * 2 is left node and p * 2 + 1 is right node. We define (h, p)
// is  merkle node at height h, position p on this level. We can virualize the tree
// like below:
//
// level h                              (h, 0)
//                                  /              \
// level h - 1                  (h-1, 0)             (h- 1, 1)
//                             /       \             /         \
// level h - 2              (h - 2, 0)  (h - 2, 1)  (h - 2, 2)  (h - 2, 3)
//                             .................................
// level 0            (0, 0)  (0, 1)  (0, 2)  (0, 3) ... (0, numberTransactions)
// We travel the tree in depth-first order. A vBits[i] is true if node i-th in DFS
// is parent of leaf node which is we want to verify, otherwise this value is false.
// vHash store hash value at node i in DFS order. Follow the vBits and vHash,
// we can rebuild the tree and extract merkle path we want.

type PartialMerkleTree struct {
	numberTransactions uint
	vBits              []bool
	vHash              []*chainhash.Hash
}

// Merkle proof use for rebuild the merkle tree and extract merkle proof for
// single transaction
type MerkleProof struct {
	// merkle node value at this sub tree  
	nodeValue  chainhash.Hash
	// merkle path if the txID we want to build merkle proof in this subtree
	// if not this is empty
	merklePath []chainhash.Hash
	// position in the level 0. We use this for check "left, right" when
	// compute merkle root.
	pos        uint32
}

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
		if numberOfHashes*chainhash.HashSize > maxAllowBytes {
			return nil, errors.New("number of hashes is too big")
		}

		bytes := make([]byte, numberOfHashes*chainhash.HashSize)
		if _, err := io.ReadFull(r, bytes); err != nil {
			return nil, err
		}

		vHash = make([]*chainhash.Hash, numberOfHashes)
		for i := 0; i < int(numberOfHashes); i++ {
			vHash[i], err = chainhash.NewHash(bytes[i*chainhash.HashSize : (i+1)*chainhash.HashSize])
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

	return &partialMerkleTree, nil
}

func ParialMerkleTreeFromHex(hexStr string) (*PartialMerkleTree, error) {
	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(hexBytes)
	return readPartialMerkleTree(reader, hexBytes)
}

func (pmt *PartialMerkleTree) CalcTreeWidth(height uint32) uint {
	return (pmt.numberTransactions + (1 << height) - 1) >> height
}

// MinHeight returns the minimum height of a Merkele tree to fit `pmt.numberTransactions`.
func (pmt *PartialMerkleTree) Height() uint32 {
	var nHeight uint32 = 0
	for pmt.CalcTreeWidth(nHeight) > 1 {
		nHeight++
	}
	return nHeight
}

// Port logic from gettxoutproof from bitcoin-core
// TODO: Make error handler more sense
func (pmt *PartialMerkleTree) computeMerkleProofRecursive(height, pos uint32, nBitUsed, nHashUsed *uint32, txID *chainhash.Hash) (*MerkleProof, error) {
	if int(*nBitUsed) >= len(pmt.vBits) {
		return nil, errors.New("Error")
	}

	fParentOfMatch := pmt.vBits[*nBitUsed]
	*nBitUsed = *nBitUsed + 1

	if height == 0 || !fParentOfMatch {
		if int(*nHashUsed) >= len(pmt.vHash) {
			return nil, errors.New("error")
		}

		hash := pmt.vHash[*nHashUsed]
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
			pos: uint32(pmt.numberTransactions) + 1,
		}, nil
		
	} else {
		left, err := pmt.computeMerkleProofRecursive(height-1, pos*2, nBitUsed, nHashUsed, txID)
		var right *MerkleProof
		if err != nil {
			return nil, err
		}
		if uint(pos*2+1) < pmt.CalcTreeWidth(height-1) {
			right, err = pmt.computeMerkleProofRecursive(height-1, pos*2+1, nBitUsed, nHashUsed, txID)
			if err != nil {
				return nil, err
			}

			if left.nodeValue.IsEqual(&right.nodeValue) {
				return nil, errors.New("error")
			}
		} else {
			right = left
		}

		nodeValue := HashNodes(&left.nodeValue, &right.nodeValue)
		// Compute new proof
		if left.pos != uint32(pmt.numberTransactions) {
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

func (pmt *PartialMerkleTree) ComputeMerkleProof(txID string) (*MerkleProof, error) {
	txIDHash, err := chainhash.NewHashFromStr(txID);
	if err != nil {
		return nil, err
	}
	height := pmt.Height();
	nUsedBit := uint32(0)
	nUsedHash := uint32(0)	
	return pmt.computeMerkleProofRecursive(height, 0, &nUsedBit, &nUsedHash, txIDHash)
}



func HashNodes(l, r *chainhash.Hash) *chainhash.Hash {
	h := make([]byte, 0, chainhash.HashSize*2)
	h = append(h, l.CloneBytes()...)
	h = append(h, r.CloneBytes()...)
	newHash := chainhash.DoubleHashH(h)
	return &newHash
}
