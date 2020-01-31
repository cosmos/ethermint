package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var typesCodec = codec.New()

func init() {
	RegisterCodec(typesCodec)
}

var (
	appName = "emint"
	// Amino encoding name
	EthermintAccountName = appName + "/Account"
)

// RegisterCodec registers all the necessary types with amino for the given
// codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&Account{}, EthermintAccountName, nil)
}
