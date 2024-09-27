package btclightclient_test

import (
	"testing"

	keepertest "bitcoin-lightclient/testutil/keeper"
	"bitcoin-lightclient/testutil/nullify"
	btclightclient "bitcoin-lightclient/x/btclightclient/module"
	"bitcoin-lightclient/x/btclightclient/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BtclightclientKeeper(t)
	btclightclient.InitGenesis(ctx, k, genesisState)
	got := btclightclient.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
