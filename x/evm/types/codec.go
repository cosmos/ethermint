package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the evm module's codec
var ModuleCdc = codec.NewLegacyAminoLegacyAmino()

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// evm module
func RegisterLegacyAminoCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	cdc.RegisterConcrete(MsgEthermint{}, "ethermint/MsgEthermint", nil)
	cdc.RegisterConcrete(TxData{}, "ethermint/TxData", nil)
	cdc.RegisterConcrete(ChainConfig{}, "ethermint/ChainConfig", nil)
}

func init() {
	RegisterLegacyAminoCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
