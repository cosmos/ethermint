package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	// EthermintAccountName is the amino encoding name for EthAccount
	EthermintAccountName = "ethermint/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthermintAccountName, nil)
}

var amino = codec.New()

func init() {
	RegisterCodec(amino)
}
