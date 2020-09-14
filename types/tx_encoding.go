package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx sdk.Tx

		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		// sdk.Tx is an interface. The concrete message types
		// are registered by MakeTxCodec
		// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		return tx, nil
	}
}
