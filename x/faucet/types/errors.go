package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// ErrWithdrawTooOften withdraw too often
	ErrWithdrawTooOften = sdkerrors.Register(ModuleName, 2, "each address can withdraw only once")
	ErrFaucetKeyEmpty   = sdkerrors.Register(ModuleName, 3, "armor should not be empty")
	ErrFaucetKeyExisted = sdkerrors.Register(ModuleName, 4, "faucet key existed")
)
