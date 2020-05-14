package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/cosmos/ethermint/x/faucet/types"
)

// RegisterRoutes register REST endpoints for the faucet module
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc("/faucet/fund", fundHandlerFn(cliCtx)).Methods("POST")
}

// PostFundRequest defines fund request's body.
type PostFundRequest struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount    sdk.Coins    `json:"amount" yaml:"amount"`
	Recipient string       `json:"receipient" yaml:"receipient"`
}

func fundHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostFundRequest
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var recipient sdk.AccAddress
		if req.Recipient == "" {
			recipient = sender
		} else {
			recipient, err = sdk.AccAddressFromBech32(req.Recipient)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgFund(req.Amount, sender, recipient)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
