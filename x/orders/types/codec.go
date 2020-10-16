package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// ModuleCdc references the global orders module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/orders and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterDerivativeMarket{},
		&MsgSuspendDerivativeMarket{},
		&MsgResumeDerivativeMarket{},
		&MsgCreateDerivativeOrder{},
		&MsgFilledDerivativeOrder{},
		&MsgCancelledDerivativeOrder{},
		&MsgRegisterSpotMarket{},
		&MsgSuspendSpotMarket{},
		&MsgResumeSpotMarket{},
		&MsgCreateSpotOrder{},
		&MsgRequestFillSpotOrder{},
		&MsgRequestSoftCancelSpotOrder{},
		&MsgFilledSpotOrder{},
		&MsgCancelledSpotOrder{},
	)
	registry.RegisterInterface("orders/EvmSyncStatus", &EvmSyncStatus{})
}
