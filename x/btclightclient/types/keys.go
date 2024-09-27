package types

const (
	// ModuleName defines the module name
	ModuleName = "btclightclient"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_btclightclient"
)

var (
	ParamsKey = []byte("p_btclightclient")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
