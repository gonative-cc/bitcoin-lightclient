package keeper

import (
	"context"

	"bitcoin-lightclient/x/btclightclient/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LatestBlock(goCtx context.Context, req *types.QueryLatestBlockRequest) (*types.QueryLatestBlockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	storeAdaptor := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdaptor, []byte{})

	// TODO: set this key as constant in types package
	value := store.Get(types.LatestBlockKey);

	var latestBlock types.BTCLightBlock

	if err := latestBlock.Unmarshal(value); err != nil {
		return nil, err
	} else {
		return &types.QueryLatestBlockResponse{Btclightblock: &latestBlock}, nil
	}
}
