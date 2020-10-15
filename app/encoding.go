// +build test_amino
package app

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/cosmos/ethermint/codec"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	codec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
