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
type partialMerkleTreeData struct {
	numberTransactions uint
	vBits              []bool
	vHash              []*chainhash.Hash
}

// Merkle proof use for rebuild the merkle tree and extract merkle proof for
// single transaction
type MerkleProof struct {
	// merkle node value at this sub tree
	merkleRoot chainhash.Hash
	// merkle path if the txID we want to build merkle proof in this subtree
	// if not this is empty
	merklePath []chainhash.Hash
	// transaction index. We use this for check "left, right" when
	// compute merkle root.
	transactionIndex uint32
}

const maxAllowBytes = 65536

// Read data for parse merkle tree. Follow encode/decode format:
// *  - uint32     total_transactions (4 bytes)
// *  - varint     number of hashes   (1-3 bytes)
// *  - uint256[]  hashes in depth-first order (<= 32*N bytes)
// *  - varint     number of bytes of flag bits (1-3 bytes)
// *  - byte[]     flag bits, packed per 8 in a byte, least significant bit first (<= 2*N-1 bits)
// This is reference from bitcoin-code.
func readPartialMerkleTreeData(buf []byte) (partialMerkleTreeData, error) {
	var pmt partialMerkleTreeData

	r := bytes.NewReader(buf)

	// read number transaction
	if _, err := io.ReadFull(r, buf[:4]); err != nil {
		return pmt, err
	}
	numberTransactions := uint(binary.LittleEndian.Uint32(buf[:4]))

	// read vhash
	var pver uint32 // pversion but btcd don't use this in those function we want.
	var vHash []*chainhash.Hash
	numberOfHashes, err := wire.ReadVarInt(r, pver)
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

	// Read vBits
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

func parialMerkleTreeDataFromHex(merkleTreeEncoded string) (partialMerkleTreeData, error) {
	b, err := hex.DecodeString(merkleTreeEncoded)
	if err != nil {
		return partialMerkleTreeData{}, err
	}

	return readPartialMerkleTreeData(b)
}

func (pmtd *partialMerkleTreeData) calcTreeWidth(height uint32) uint {
	return (pmtd.numberTransactions + (1 << height) - 1) >> height
}

// returns the minimum height of a Merkele tree to fit `pmt.numberTransactions`.
func (pmtd *partialMerkleTreeData) height() uint32 {
	var nHeight uint32 = 0
	for pmtd.calcTreeWidth(nHeight) > 1 {
		nHeight++
	}
	return nHeight
}

func (pmtd *partialMerkleTreeData) buildMerkleTreeRecursive(height, pos uint32, nBitUsed, nHashUsed *uint32, merkleTree *PartialMerkleTree) (*chainhash.Hash, error) {
	if int(*nBitUsed) >= len(pmtd.vBits) {
		return nil, fmt.Errorf("Out-bound of vBits")
	}

	fParentOfMatch := pmtd.vBits[*nBitUsed]
	*nBitUsed = *nBitUsed + 1

	// handle leaf
	if height == 0 || !fParentOfMatch {
		if int(*nHashUsed) >= len(pmtd.vHash) {
			fmt.Println(len(pmtd.vHash))
			return nil, fmt.Errorf("Out-bound of vHash")
		}
		hash := pmtd.vHash[*nHashUsed]
		*nHashUsed++
		merkleTree.nodesAtHeight[height][pos] = *hash
		return hash, nil
	} else {
		left, err := pmtd.buildMerkleTreeRecursive(height-1, pos*2, nBitUsed, nHashUsed, merkleTree)
		if err != nil {
			return nil, err
		}
		var right *chainhash.Hash
		if uint(pos*2+1) < pmtd.calcTreeWidth(height-1) {
			right, err = pmtd.buildMerkleTreeRecursive(height-1, pos*2+1, nBitUsed, nHashUsed, merkleTree)
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
		merkleTree.nodesAtHeight[height][pos] = *nodeValue
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

type merkleNodes map[uint32]chainhash.Hash
type PartialMerkleTree struct {
	// nodes at level or height.
	// check
	nodesAtHeight []merkleNodes
}

func (mk PartialMerkleTree) getLeafNodeIndex(txID *chainhash.Hash) (uint32, error) {
	// TODO(vu): Should we use reverse map to find position of merkle leaf?
	for leafIndex, leafValue := range mk.nodesAtHeight[0] {
		if leafValue.IsEqual(txID) {
			return leafIndex, nil
		}
	}
	return 0, errors.New("Node value doesn't exist in merkle tree")

}

// Return merkle proof of txID.
// Return error this txID doesn't exist in merkle tree
func (mk PartialMerkleTree) GetProof(txID string) (*MerkleProof, error) {
	txHash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, err
	}

	transactionIndex, err := mk.getLeafNodeIndex(txHash)
	if err != nil {
		return nil, err
	}

	merklePath := []chainhash.Hash{*txHash}
	h := len(mk.nodesAtHeight)

	// We have merkle tree below:
	// level h                              (h, 0)
	//                                  /              \
	// level h - 1                  (h-1, 0)             (h- 1, 1)
	//                             /       \             /         \
	// level h - 2              (h - 2, 0)  (h - 2, 1)  (h - 2, 2)  (h - 2, 3)
	//                             .................................
	// level 0            (0, 0)  (0, 1)  (0, 2)  (0, transactionIndex) ... (0, numberTransactions - 1)
	// Easy obvious:
	// - if position is event, we need the sibling at position + 1 to compute the parent's hash.
	// - if position is odd, we need the sibling at position - 1 to compute the parent's hash.
	position := transactionIndex
	// The node at level 0 is merkle root :D 
	for i := 0; i < h-1; i++ {
		var siblingHash chainhash.Hash
		if position%2 == 0 {
			siblingHash = mk.nodesAtHeight[i][position+1]
		} else {
			siblingHash = mk.nodesAtHeight[i][position-1]
		}
		
		merklePath = append(merklePath, siblingHash)
		position = position / 2
	}

	return &MerkleProof{
		merkleRoot:       mk.nodesAtHeight[h-1][0],
		merklePath:       merklePath,
		transactionIndex: transactionIndex,
	}, nil
}

// Parse hex data(gextxoutproof) to PartialMerkleTree
func PartialMerkleTreeFromHex(mtData string) (PartialMerkleTree, error) {
	var pmt PartialMerkleTree

	// decode information for build partial merkle tree
	pmtInfo, err := parialMerkleTreeDataFromHex(mtData)
	if err != nil {
		return pmt, err
	}

	nBitUsed := uint32(0)
	nHashUsed := uint32(0)
	height := pmtInfo.height()
	pmt.nodesAtHeight = make([]merkleNodes, height+1)
	for i := 0; i <= int(height); i++ {
		pmt.nodesAtHeight[i] = make(map[uint32]chainhash.Hash)
	}
	if _, err := pmtInfo.buildMerkleTreeRecursive(height, 0, &nBitUsed, &nHashUsed, &pmt); err != nil {
		return pmt, err
	}
	return pmt, nil
}