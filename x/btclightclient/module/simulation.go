package btclightclient

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"bitcoin-lightclient/testutil/sample"
	btclightclientsimulation "bitcoin-lightclient/x/btclightclient/simulation"
	"bitcoin-lightclient/x/btclightclient/types"
)

// avoid unused import issue
var (
	_ = btclightclientsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgInsertHeaders = "op_weight_msg_insert_headers"
	// TODO: Determine the simulation weight value
	defaultWeightMsgInsertHeaders int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	btclightclientGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&btclightclientGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgInsertHeaders int
	simState.AppParams.GetOrGenerate(opWeightMsgInsertHeaders, &weightMsgInsertHeaders, nil,
		func(_ *rand.Rand) {
			weightMsgInsertHeaders = defaultWeightMsgInsertHeaders
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInsertHeaders,
		btclightclientsimulation.SimulateMsgInsertHeaders(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgInsertHeaders,
			defaultWeightMsgInsertHeaders,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				btclightclientsimulation.SimulateMsgInsertHeaders(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
