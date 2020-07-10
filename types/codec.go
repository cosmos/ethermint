package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "ethermint/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)
}

var (
	amino = codec.New()

	// ModuleCdc references the global x/auth module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/auth and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino)
)

func init() {
	RegisterCodec(amino)
	codec.RegisterCrypto(amino)
}
