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


	value := store.Get(types.LatestBlockKey)

	var latestBlock types.BTCLightBlock

	if err := latestBlock.Unmarshal(value); err != nil {
		return nil, err
	} else {
		// TODO: make logic more simple
		var headerBytes types.BTCHeaderBytes

		headerBytes, _ = types.ByteFromBlockHeader(latestBlock.Header)
		return &types.QueryLatestBlockResponse{Height: int64(latestBlock.Height), HeaderHex: headerBytes.MarshalHex()}, nil
	}
}
