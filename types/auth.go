package types

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// CompleteAndBroadcastTx implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func CompleteAndBroadcastTx(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (sdk.TxResponse, error) {
	txBldr, err := authutils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = authutils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		gasEst := authutils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, keys.DefaultKeyPass, msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	return cliCtx.BroadcastTx(txBytes)
}
