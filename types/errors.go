package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// RootCodespace is the codespace for all errors defined in this package
	RootCodespace = "ethermint"
)

// NOTE: We can't use 1 since that error code is reserved for internal errors.

var (
	// ErrInvalidValue returns an error resulting from an invalid value.
	ErrInvalidValue = sdkerrors.Register(RootCodespace, 2, "invalid value")

	// ErrInvalidChainID returns an error resulting from an invalid chain ID.
	ErrInvalidChainID = sdkerrors.Register(RootCodespace, 3, "invalid chain ID")

	// ErrVMExecution returns an error resulting from an error in EVM execution.
	ErrVMExecution = sdkerrors.Register(RootCodespace, 4, "error while executing evm transaction")

	// ErrMarshalBigInt returns an error resulting from marshaling a big.Int from a string.
	ErrMarshalBigInt = sdkerrors.Register(RootCodespace, 5, "cannot unmarshal big.Int from string")

	// ErrUnmarshalBigInt returns an error resulting from unmarshaling a big.Int from a string.
	ErrUnmarshalBigInt = sdkerrors.Register(RootCodespace, 6, "cannot unmarshal big.Int from string")
)
