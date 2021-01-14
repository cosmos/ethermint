package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/cosmos/ethermint/x/evm/types"
)

type ValidateMsgHandler func(ctx sdk.Context, msgs []sdk.Msg) error

type ValidateMsgHandlerDecorator struct {
}

func NewValidateMsgHandlerDecorator() ValidateMsgHandlerDecorator {
	return ValidateMsgHandlerDecorator{}
}

func (vmhd ValidateMsgHandlerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// *ABORT* the tx in case of failing to validate it in checkTx mode
	if ctx.IsCheckTx() {
		wrongMsgErr := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest,
			"It is not allowed that a transaction with more than one message and contains evm message")
		msgs := tx.GetMsgs()
		msgNum := len(msgs)
		if msgNum == 0 {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "The msg array from tx cannot be empty")
		} else if msgNum == 1 {
			return next(ctx, tx, simulate)
		}

		for _, msg := range msgs {
			switch msg.(type) {
			case evmtypes.MsgEthermint:
				return ctx, wrongMsgErr
			}
		}
	}

	return next(ctx, tx, simulate)
}
