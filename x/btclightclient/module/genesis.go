package btclightclient

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitcoin-lightclient/x/btclightclient/keeper"
	"bitcoin-lightclient/x/btclightclient/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	k.Logger().Debug(genState.Header, genState.Height)
	
	if err := k.InitGenesisBTCBlock(ctx, genState.Height, genState.Header); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	// genesis.Params = k.GetParams(ctx)
	// latestBlock, _ := k.LatestBlock(ctx, &types.QueryLatestBlockRequest{})
	// genesis.Height = uint64(latestBlock.Height)
	// genesis.Header = latestBlock.HeaderHex
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
