package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"

	emintcrypto "github.com/cosmos/ethermint/crypto"
	ethermint "github.com/cosmos/ethermint/types"
)

// MakeCodecs constructs the *std.Codec and *codec.Codec instances used by the
// Ethermint app. It is useful for tests and clients who do not want to construct
// the full application.
func MakeCodecs(bm module.BasicManager) (*std.Codec, *codec.Codec) {
	cdc := std.MakeCodec(bm)
	emintcrypto.RegisterCodec(cdc)
	keyring.RegisterCodec(cdc) // temporary. Used to register keyring.Info
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	RegisterInterfaces(interfaceRegistry)
	bm.RegisterInterfaceModules(interfaceRegistry)
	appCodec := std.NewAppCodec(cdc, interfaceRegistry)
	return appCodec, cdc
}

// RegisterInterfaces registers Interfaces from sdk/types and vesting
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	std.RegisterInterfaces(registry)
	ethermint.RegisterInterfaces(registry)
}
