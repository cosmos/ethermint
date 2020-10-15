package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrOrderInvalid    = sdkerrors.Register(ModuleName, 1, "failed to validate order")
	ErrOrderNotFound   = sdkerrors.Register(ModuleName, 2, "no active order found for the hash")
	ErrPairSuspended   = sdkerrors.Register(ModuleName, 3, "trade pair suspended")
	ErrPairNotFound    = sdkerrors.Register(ModuleName, 4, "trade pair not found")
	ErrPairExists      = sdkerrors.Register(ModuleName, 5, "trade pair exists")
	ErrPairMismatch    = sdkerrors.Register(ModuleName, 6, "trade pair mismatch")
	ErrBadField        = sdkerrors.Register(ModuleName, 7, "struct field error")
	ErrMarketNotFound  = sdkerrors.Register(ModuleName, 8, "derivative market not found")
	ErrMarketInvalid   = sdkerrors.Register(ModuleName, 9, "failed to validate derivative market")
	ErrMarketExists    = sdkerrors.Register(ModuleName, 10, "market exists")
	ErrMarketSuspended = sdkerrors.Register(ModuleName, 11, "market suspended")
	ErrBadUpdateEvent  = sdkerrors.Register(ModuleName, 12, "order update event not confirmed")
	ErrUpdateSameValue = sdkerrors.Register(ModuleName, 13, "cannot update the record's field with the same value")
)
