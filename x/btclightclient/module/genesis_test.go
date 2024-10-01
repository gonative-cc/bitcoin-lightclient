package btclightclient_test

import (
	"fmt"
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
		Height: 863651,
		Header: "0000002ef43043c56fb82eb74f8f0e12b13ae548f222c300237802000000000000000000142e689b4084b19ecf9fcbc46f590b206422e35c61103dd3c55ad4df7d40d6fc50fdfb66142f0317971b98b5",
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BtclightclientKeeper(t)
	btclightclient.InitGenesis(ctx, k, genesisState)
	got := btclightclient.ExportGenesis(ctx, k)

	fmt.Println(got)
	require.NotNil(t, got)
	
	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
