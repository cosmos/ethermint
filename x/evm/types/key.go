package types

const (
	// ModuleName string name of module
	ModuleName = "evm"

	// StoreKey key for ethereum storage data (StateDB)
	StoreKey = ModuleName
	// CodeKey key for ethereum code data
	CodeKey = ModuleName + "code"
	// BlockKey key
	BlockKey = ModuleName + "block"

	// RouterKey uses module name for routing
	RouterKey = ModuleName
)
