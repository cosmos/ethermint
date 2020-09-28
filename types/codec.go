package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "ethermint/EthAccount"
)

// RegisterLegacyAminoCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterLegacyAminoCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)
}
