package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var typesCodec = codec.New()

func init() {
	RegisterCodec(typesCodec)
}

const (
	// Amino encoding name
	EthermintAccountName = "ethermint/EthAccount"
)

// RegisterCodec registers all the necessary types with amino for the given
// codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthermintAccountName, nil)
}
