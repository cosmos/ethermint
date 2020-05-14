package types

const (
	// ModuleName is the name of the module
	ModuleName = "faucet"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey uses module name for routing
	RouterKey = ModuleName
)

var (
	EnableFaucetKey  = []byte{0x01}
	TimeoutKey       = []byte{0x02}
	CapKey           = []byte{0x03}
	MaxPerRequestKey = []byte{0x04}
	FundedKey        = []byte{0x05}
)
