package btclightclient

import "errors"

var ErrForkTooOld = errors.New("fork too old")
var ErrBlockNotInChain = errors.New("block not in the chain")
var ErrInvalidHeaderSize = errors.New("invalid header size, must be 80 bytes")
var ErrParentBlockNotInChain = errors.New("parent block not in chain")
