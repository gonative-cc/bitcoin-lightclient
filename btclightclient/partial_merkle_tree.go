package btclightclient

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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
// level 0            (0, 0)  (0, 1)  (0, 2)  (0, 3) ... (0, numberTransactions - 1)
// We travel the tree in depth-first order. A vBits[i] is true if node i-th in DFS
// is parent of leaf node which is we want to verify, otherwise this value is false.
// vHash store hash value at node i in DFS order. Follow the vBits and vHash,
// we can rebuild the tree and extract merkle path we want.

type PartialMerkleTreeData struct {
	numberTransactions uint
	vBits              []bool
	vHash              []*chainhash.Hash
}


type merkleNodes map[uint32]chainhash.Hash

type PartialMerkleTree struct {
	nodesAtHeight []merkleNodes
}

func (mk PartialMerkleTree) getLeafNodeIndex(txID *chainhash.Hash) (uint32, error) {
	// TODO(vu): Should we use reverse map to find position of merkle leaf?
	h := len(mk.nodesAtHeight)
	for leafIndex, leafValue := range mk.nodesAtHeight[h - 1] {
		if leafValue.IsEqual(txID) {
			return leafIndex, nil
		}
	}
	return 0, errors.New("Node value doesn't exist in merkle tree")
	
}
func (mk PartialMerkleTree) GetProof(txID *chainhash.Hash) (*MerkleProof, error) {
	
	txIndex, err := mk.getLeafNodeIndex(txID);
	if err != nil {
		return nil, err 
	}

	merklePath := []chainhash.Hash{}
	h := len(mk.nodesAtHeight);
	for i := h - 1; i > 0; i-- {
		if txIndex % 2 == 0 {
			merklePath = append(merklePath, mk.nodesAtHeight[i][txIndex + 1]);
		} else {
			merklePath = append(merklePath, mk.nodesAtHeight[i][txIndex - 1]);
		}
		txIndex = txIndex / 2
	}

	// TODO: should we return 
	return &MerkleProof{
		merkleRoot: mk.nodesAtHeight[0][0],
		merklePath: merklePath,
		pos: txIndex,
	}, nil 
}

// Merkle proof use for rebuild the merkle tree and extract merkle proof for
// single transaction
type MerkleProof struct {
	// merkle node value at this sub tree
	merkleRoot chainhash.Hash
	// merkle path if the txID we want to build merkle proof in this subtree
	// if not this is empty
	merklePath []chainhash.Hash
	// position in the level 0. We use this for check "left, right" when
	// compute merkle root.
	pos uint32
}

const maxAllowBytes = 65536


func merkleProofAtLeaf(leafHash chainhash.Hash, position uint32) *MerkleProof {
	return &MerkleProof{
		merkleRoot: leafHash,
		merklePath: []chainhash.Hash{leafHash},
		pos:        position,
	}
}

func emptyMekleProof() *MerkleProof {
	return &MerkleProof {
		
	}
}

func readPartialMerkleTreeData(r io.Reader, buf []byte) (PartialMerkleTreeData, error) {
	var pmt PartialMerkleTreeData

	if _, err := io.ReadFull(r, buf[:4]); err != nil {
		return pmt, err
	}

	numberTransactions := uint(binary.LittleEndian.Uint32(buf[:4]))
	var pver uint32 // pversion but btcd don't use this in those function we want.
	var vHash []*chainhash.Hash

	numberOfHashes, err := wire.ReadVarIntBuf(r, pver, buf)
	if err != nil {
		return pmt, err
	}
	if numberOfHashes*chainhash.HashSize > maxAllowBytes {
		return pmt, errors.New("number of hashes is too big")
	}

	bytes := make([]byte, numberOfHashes*chainhash.HashSize)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return pmt, err
	}
	vHash = make([]*chainhash.Hash, numberOfHashes)
	for i := 0; i < int(numberOfHashes); i++ {
		vHash[i], err = chainhash.NewHash(bytes[i*chainhash.HashSize : (i+1)*chainhash.HashSize])
		if err != nil {
			return pmt, err
		}
	}

	var vBits []bool
	vBytes, err := wire.ReadVarBytes(r, pver, uint32(maxAllowBytes), "vBits")
	if err != nil {
		return pmt, err
	}
	vBits = make([]bool, len(vBytes)*8)
	i := 0
	for _, b := range vBytes {
		for j := 0; j < 8; j++ {
			vBits[i] = (b & (1 << j)) != 0
			i++
		}
	}

	pmt.numberTransactions = numberTransactions
	pmt.vBits = vBits
	pmt.vHash = vHash
	return pmt, nil
}

func ParialMerkleTreeFromHex(merkleTreeEncoded string) (PartialMerkleTreeData, error) {

	b, err := hex.DecodeString(merkleTreeEncoded)
	if err != nil {
		return PartialMerkleTreeData{}, err
	}

	r := bytes.NewReader(b)
	return readPartialMerkleTreeData(r, b)
}

func (pmt *PartialMerkleTreeData) CalcTreeWidth(height uint32) uint {
	return (pmt.numberTransactions + (1 << height) - 1) >> height
}

// MinHeight returns the minimum height of a Merkele tree to fit `pmt.numberTransactions`.
func (pmt *PartialMerkleTreeData) Height() uint32 {
	var nHeight uint32 = 0
	for pmt.CalcTreeWidth(nHeight) > 1 {
		nHeight++
	}
	return nHeight
}


// Port logic from gettxoutproof from bitcoin-core
func (pmt *PartialMerkleTreeData) buildMerkleTreeRecursive(height, pos uint32, nBitUsed, nHashUsed *uint32, merkleTree *PartialMerkleTree) (*chainhash.Hash, error) {
	if int(*nBitUsed) >= len(pmt.vBits) {
		return nil, fmt.Errorf("Out-bound of vBits")
	}

	fParentOfMatch := pmt.vBits[*nBitUsed]
	*nBitUsed = *nBitUsed + 1

	// handle leaf  
	if height == 0 || !fParentOfMatch {
		if int(*nHashUsed) >= len(pmt.vHash) {
			return nil, fmt.Errorf("Out-bound of vHash")
		}
		hash := pmt.vHash[*nHashUsed]
		*nHashUsed++
		merkleTree.nodesAtHeight[height][pos] = *hash
		return hash, nil
	} else {
		left, err := pmt.buildMerkleTreeRecursive(height-1, pos*2, nBitUsed, nHashUsed, merkleTree)
		if err != nil {
			return nil, err
		}

		var right *chainhash.Hash
		if uint(pos*2+1) < pmt.CalcTreeWidth(height-1) {
			right, err = pmt.buildMerkleTreeRecursive(height-1, pos*2+1, nBitUsed, nHashUsed, merkleTree)
			if err != nil {
				return nil, err
			}

			if left.IsEqual(right) {
				return nil, fmt.Errorf("In the case tree width is old, the last hash must be duplicate")
			}
		} else {
			right = left
		}

		nodeValue := HashNodes(left, right)
		return nodeValue, nil
	}
}


func HashNodes(l, r *chainhash.Hash) *chainhash.Hash {
	h := make([]byte, 0, chainhash.HashSize*2)
	h = append(h, l.CloneBytes()...)
	h = append(h, r.CloneBytes()...)
	newHash := chainhash.DoubleHashH(h)
	return &newHash
}
