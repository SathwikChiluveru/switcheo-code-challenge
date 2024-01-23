package types

const (
	// ModuleName defines the module name
	ModuleName = "myblockchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_myblockchain"
)

var (
	ParamsKey = []byte("p_myblockchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
