package types

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
		Height: 863651,
		Header: `0000002ef43043c56fb82eb74f8f0e12b13ae548f222c300237802000000000000000000142e689b4084b19ecf9fcbc46f590b206422e35c61103dd3c55ad4df7d40d6fc50fdfb66142f0317971b98b5`,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
