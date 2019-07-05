package types

import "github.com/cosmos/cosmos-sdk/codec"

// TODO: Is this the MAIN codec? or should it only be applied to eth specific messages (originally `MsgCodec`)
var ModuleCdc = codec.New()

func init() {
	cdc := codec.New()

	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	ModuleCdc = cdc.Seal()
}

// RegisterCodec registers concrete types and interfaces on the given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthereumTxMsg{}, "ethermint/MsgEthereumTx", nil)
}
