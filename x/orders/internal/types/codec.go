package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterDerivativeMarket{}, "orders/MsgRegisterDerivativeMarket", nil)
	cdc.RegisterConcrete(MsgSuspendDerivativeMarket{}, "orders/MsgSuspendDerivativeMarket", nil)
	cdc.RegisterConcrete(MsgResumeDerivativeMarket{}, "orders/MsgResumeDerivativeMarket", nil)
	cdc.RegisterConcrete(MsgCreateDerivativeOrder{}, "orders/MsgCreateDerivativeOrder", nil)
	cdc.RegisterConcrete(MsgFilledDerivativeOrder{}, "orders/MsgFilledDerivativeOrder", nil)
	cdc.RegisterConcrete(MsgCancelledDerivativeOrder{}, "orders/MsgCancelledDerivativeOrder", nil)

	cdc.RegisterConcrete(MsgRegisterSpotMarket{}, "orders/MsgRegisterSpotMarket", nil)
	cdc.RegisterConcrete(MsgSuspendSpotMarket{}, "orders/MsgSuspendSpotMarket", nil)
	cdc.RegisterConcrete(MsgResumeSpotMarket{}, "orders/MsgResumeSpotMarket", nil)
	cdc.RegisterConcrete(MsgCreateSpotOrder{}, "orders/MsgCreateSpotOrder", nil)
	cdc.RegisterConcrete(MsgRequestFillSpotOrder{}, "orders/MsgRequestFillSpotOrder", nil)
	cdc.RegisterConcrete(MsgRequestSoftCancelSpotOrder{}, "orders/MsgRequestSoftCancelSpotOrder", nil)
	cdc.RegisterConcrete(MsgFilledSpotOrder{}, "orders/MsgFilledSpotOrder", nil)
	cdc.RegisterConcrete(MsgCancelledSpotOrder{}, "orders/MsgCancelledSpotOrder", nil)

	cdc.RegisterConcrete(EvmSyncStatus{}, "orders/EvmSyncStatus", nil)
}

// Enforce the msg types at compile time
var (
	_ sdk.Msg = MsgRegisterDerivativeMarket{}
	_ sdk.Msg = MsgSuspendDerivativeMarket{}
	_ sdk.Msg = MsgResumeDerivativeMarket{}
	_ sdk.Msg = MsgCreateDerivativeOrder{}
	_ sdk.Msg = MsgFilledDerivativeOrder{}
	_ sdk.Msg = MsgCancelledDerivativeOrder{}

	_ sdk.Msg = MsgRegisterSpotMarket{}
	_ sdk.Msg = MsgSuspendSpotMarket{}
	_ sdk.Msg = MsgResumeSpotMarket{}
	_ sdk.Msg = MsgCreateSpotOrder{}
	_ sdk.Msg = MsgRequestFillSpotOrder{}
	_ sdk.Msg = MsgRequestSoftCancelSpotOrder{}
	_ sdk.Msg = MsgFilledSpotOrder{}
	_ sdk.Msg = MsgCancelledSpotOrder{}
)
