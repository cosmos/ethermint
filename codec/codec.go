package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"

	emintcrypto "github.com/cosmos/ethermint/crypto"
)

// MakeCodec registers the necessary types and interfaces for an sdk.App. This
// codec is provided to all the modules the application depends on.
//
// NOTE: This codec will be deprecated in favor of AppCodec once all modules are
// migrated.
func MakeCodec(bm module.BasicManager) *codec.Codec {
	cdc := std.MakeCodec(bm)
	emintcrypto.RegisterCodec(cdc)
	keyring.RegisterCodec(cdc) // temporary. Used to register keyring.Info
	return cdc
}
