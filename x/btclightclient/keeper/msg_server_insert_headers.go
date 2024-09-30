package keeper

import (
	"context"
	// "errors"
	// errorsmod "cosmossdk.io/errors"
	"bitcoin-lightclient/x/btclightclient/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/btcsuite/btcd/wire"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InsertHeaders(goCtx context.Context, msg *types.MsgInsertHeaders) (*types.MsgInsertHeadersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	headers := msg.Headers

	blockHeaders := make([]*wire.BlockHeader, len(headers))

	// transform string to wire.BlockHeader
	for i, header := range headers {
		blockHeader, err := types.NewBlockHeaderFromBytes([]byte(header))
		// TODO: reorg this code more readable
		if err != nil {
			ctx.Logger().With("module", "x/btclightclient").Error("This is error when cover", "err", err)
			return nil, sdkerrors.ErrTxDecode
		} else {
			blockHeaders[i] = blockHeader
		}
	}

	k.Keeper.InsertHeader(ctx, blockHeaders)
	return &types.MsgInsertHeadersResponse{}, nil
}
