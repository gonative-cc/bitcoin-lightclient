package btclightclient

import "errors"

// Light client errors
var ErrForkTooOld = errors.New("fork too old")
var ErrBlockNotInChain = errors.New("block not in the chain")
var ErrInvalidHeaderSize = errors.New("invalid header size, must be 80 bytes")
var ErrParentBlockNotInChain = errors.New("parent block not in chain")
var ErrBlockIsNotForkHead = errors.New("block is not a fork head")

// SPV errors
var ErrValueIsNotMerkleLeaf = errors.New("value doesn't exist in merkle tree")
var ErrMerkleDecodeOutbound = errors.New("out-bound of vHash")
var ErrMerkleDecodeHashNumberInvalid = errors.New("number of hashes reach to limit")
