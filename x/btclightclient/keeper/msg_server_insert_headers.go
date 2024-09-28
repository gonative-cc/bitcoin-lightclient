package keeper

import (
	"context"

	"bitcoin-lightclient/x/btclightclient/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InsertHeaders(goCtx context.Context, msg *types.MsgInsertHeaders) (*types.MsgInsertHeadersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Keeper.InsertHeader(ctx, []*types.Btcheader{})
	return &types.MsgInsertHeadersResponse{}, nil
}
