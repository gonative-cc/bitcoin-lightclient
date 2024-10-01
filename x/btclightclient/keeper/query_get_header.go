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

func (k Keeper) GetHeader(goCtx context.Context, req *types.QueryGetHeaderRequest) (*types.QueryGetHeaderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	storeAdaptor := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdaptor, []byte{})

	height := req.Height
	key, _ := types.HeaderKey(height)
	value := store.Get(key)

	var lightblock types.BTCLightBlock
	
	if err := lightblock.Unmarshal(value); err == nil {
		var headerBytes types.BTCHeaderBytes
		headerBytes, _ = types.ByteFromBlockHeader(lightblock.Header)
		
		return &types.QueryGetHeaderResponse{Height: int64(lightblock.Height), HeaderHex: headerBytes.MarshalHex()}, nil
	} else {
		return nil, err
	}

}
