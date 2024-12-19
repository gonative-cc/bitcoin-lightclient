package btclightclient

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	// "fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type Hash256Digest [32]byte

type VerifyStatus int

type SPVProof struct {
	blockHash   *chainhash.Hash
	txId        *chainhash.Hash
	txIndex     uint
	merkleProof []byte
}

type PMerkleTree struct {
	numberTransactions uint
	vBits              []bool
	vHash              []*chainhash.Hash
}


func readPMerkleTree(r io.Reader, pMerkletree *PMerkleTree, buf []byte) error {
	if _, err := io.ReadFull(r, buf[:4]); err != nil {
		return err
	}

	pMerkletree.numberTransactions = uint(binary.LittleEndian.Uint32(buf[:4]))

	var pver uint32; // pversion but btcd don't use this in those function we want.

	// TODO: verify maxAllowBytes
	maxAllowBytes := 65536
	
	if numberHash, err := wire.ReadVarIntBuf(r, pver, buf); err != nil {
		return err
	} else {
		if numberHash * 32 > uint64(maxAllowBytes) {
			return errors.New("number of hash is too big")
		}

		bytes := make([]byte, numberHash * 32);
		_, err := io.ReadFull(r, bytes);

		if err != nil {
			return err
		}

		vHash := make([]*chainhash.Hash, numberHash)
		for i := 0; i < int(numberHash); i++ {
			hash, err := chainhash.NewHash(bytes[i * 32 : (i + 1) * 32]);
			if err != nil {
				return err
			}
			vHash[i] = hash
		}

		
		pMerkletree.vHash = vHash
	}

	
	return nil
	
}



func PMerkleTreeFromBytes(hexStr string) PMerkleTree {
	// 4 bytes for numnerTransactions
	fmt.Println(hexStr);
	proof, _ := hex.DecodeString(hexStr)
	numberTransaction := binary.LittleEndian.Uint32(proof[:4])

	// VHash decode
	
	// vBits decode
	return PMerkleTree{
		numberTransactions: uint(numberTransaction),
		vBits:              make([]bool, 0),
		vHash:              make([]*chainhash.Hash, 0),
	}
}

const (
	// proof is complete wrong
	InValidTXOut VerifyStatus = iota
	// proof valid but the block have not "finalized" yet.
	ParialValidTXOut
	// proof vlaid and block is "finalized".
	ValidTXOut
)

// We copy logic from bitcoin-spv. The main reason is bitcoin-spv is not maintain anymore.
// https://github.com/summa-tx/bitcoin-spv/
// Thank summa-tx for their awesome work

// Hash256MerkleStep concatenates and hashes two inputs for merkle proving
func Hash256MerkleStep(a, b []byte) Hash256Digest {
	c := []byte{}
	c = append(c, a...)
	c = append(c, b...)
	return DoubleHash(c)
}

func DoubleHash(in []byte) Hash256Digest {
	first := sha256.Sum256(in)
	second := sha256.Sum256(first[:])
	return Hash256Digest(second)
}

// follow logic on bitcoin-spv.
// This is check the tx belong to merkle tree hash in BTC header.
func VerifyHash256Merkle(proof []byte, index uint) bool {
	var current Hash256Digest
	idx := index
	proofLength := len(proof)

	if proofLength%32 != 0 {
		return false
	}
	if proofLength == 32 {
		return true
	}
	if proofLength == 64 {
		return false
	}

	root := proof[proofLength-32:]
	cur := proof[:32:32]
	copy(current[:], cur)
	numSteps := (proofLength / 32) - 1

	for i := 1; i < numSteps; i++ {
		start := i * 32
		end := i*32 + 32
		next := proof[start:end:end]
		if idx%2 == 1 {
			current = Hash256MerkleStep(next, current[:])
		} else {
			current = Hash256MerkleStep(current[:], next)
		}
		idx >>= 1
	}

	return bytes.Equal(current[:], root)
}

func (lc *BTCLightClient) VerifySPV(spvProof SPVProof) VerifyStatus {

	lightBlock := lc.btcStore.LightBlockByHash(*spvProof.blockHash)

	// In the case light block not belong currect database
	if lightBlock == nil {
		return InValidTXOut
	}

	proof := []byte{}
	proof = append(proof, spvProof.txId[:]...)
	proof = append(proof, spvProof.merkleProof...)
	merkleRoot := lightBlock.Header.BlockHash()
	proof = append(proof, merkleRoot[:]...)

	validProof := VerifyHash256Merkle(proof, spvProof.txIndex)

	if !validProof {
		return InValidTXOut
	}

	// in the case the block not finalize
	if lc.btcStore.LatestHeight() < int64(lightBlock.Height) {
		return ParialValidTXOut
	}

	return ValidTXOut
}

// / Get Merkle proof for single transaction from gettxoutproof output
func (lc *BTCLightClient) VerifySPVFromHex(hexStr string, txId *chainhash.Hash) error {
	return nil
}

func SPVFromHex(hexStr string) (*SPVProof, error) {
	// get header from gettxoutproof encode
	header, err := BlockHeaderFromHex(hexStr[0:160])

	if err != nil {
		return nil, err
	}

	blockHash := header.BlockHash()

	// txId :=
	return &SPVProof{
		blockHash:   &blockHash,
		txId:        nil,
		txIndex:     0,
		merkleProof: make([]byte, 0),
	}, nil

}
