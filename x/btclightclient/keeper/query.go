package keeper

import (
	"bitcoin-lightclient/x/btclightclient/types"
)

var _ types.QueryServer = Keeper{}
