package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the necessary x/evm interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	cdc.RegisterConcrete(MsgEthermint{}, "ethermint/MsgEthermint", nil)
	cdc.RegisterConcrete(TxData{}, "ethermint/TxData", nil)
}

var (
	// AminoCdc defines the amino codec.
	// NOTE: make private or remove once proto migration on SDK is complete.
	AminoCdc = codec.New()

	// ModuleCdc references the global x/evm module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/evm and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(AminoCdc)
)

func init() {
	RegisterCodec(AminoCdc)
	codec.RegisterCrypto(AminoCdc)
	AminoCdc.Seal()
}
