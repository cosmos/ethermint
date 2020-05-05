package faucet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/codec"
	"github.com/cosmos/ethermint/x/evm"
)

// Faucet represents the necessary data for connecting to and indentifying a Faucet and its counterparites
type Faucet struct {
	// addresses can send requests every <timeout> duration 
	timeout time.Duration
	// stores faucet addresses that have been used reciently
	faucetAddrs map[string]time.Time

	// ethermint std codec
	cdc  codec.Codec 
}

// New creates a new Faucet instance
func New(timeout time.Duration) Faucet {
	cdc := codec.MakeCodec(app.ModuleBasics)
	appCodec := codec.NewAppCodec(cdc)

	return Faucet {
		timeout: timeout,
		faucetAddrs: make(map[string]time.Time)
		cdc: appCodec,
	}
}

// Handler listens for addresses
func (f Faucet) Handler(from sdk.AccAddress, amount sdk.Coin) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload: %w", err)
			return
		}
		
		defer r.Body.Close()

		if err := req.Validate(); err != nil {
			return err
		}

		if err := f.rateLimit(fr.Address); err != nil {
			rest.WriteErrorResponse(w, w, http.StatusTooManyRequests, err)
			return
		}

		// error validated on req.Validate
		to, _ := req.GetAddress()

		if status, err := f.transfer(fromKey, to, amount); err != nil {
			rest.WriteErrorResponse(w, status, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(r.Body)
	}
}

// TODO: 
// - introduce a counter for the number of requests within <timeout>
// - cap amount sent to address
func (f Faucet) rateLimit(address string) error {
	// first time requester, can send request 
	if lastRequest, ok := f.faucetAddrs[addr]; !ok {
		f.faucetAddrs[addr] = time.Now().UTC()
		return nil
	}

	sinceLastRequest := time.Since(lastRequest)
	if f.timeout > sinceLastRequest {
		wait := f.timeout - sinceLastRequest
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", addr, f.timeout.String(), wait.String())
	}

	// user able to send funds since they have waited for period
	f.faucetAddrs[addr] = time.Now().UTC()
	return nil
}

// TODO: prob use an sdk.Int instead of coins.
func (f Faucet) transfer(from, to sdk.AccAddress, amount sdk.Coins) (int, error) {
	_, err := f.Keybase.KeyByAddress(from)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_ := bank.NewMsgSend(from, to, sdk.NewCoins(amount))
	// res, err := f.SendMsgWithKey(msg, info.GetName())
	// if err != nil {
	// 	return http.StatusBadRequest, err
	// }

	return http.StatusOK, nil
}


// buildAndSignTxWithKey allows the user to specify which relayer key will sign the message
func (f Faucet) buildAndSignTxWithKey(msgs []sdk.Msg, keyName string) ([]byte, error) {
	// Fetch account and sequence numbers for the account
	info, err := f.Keybase.Key(keyName)
	if err != nil {
		return nil, err
	}

	// done := f.UseSDKContext()
	// defer done()

	acc, err := auth.NewAccountRetriever(f.Cdc, f).GetAccount(info.GetAddress())
	if err != nil {
		return nil, err
	}

	return auth.NewTxBuilder(
		auth.DefaultTxEncoder(f.Amino.Codec),
		acc.GetAccountNumber(),
		acc.GetSequence(), f.Gas, f.GasAdjustment, false, f.ChainID,
		f.Memo, sdk.NewCoins(), f.getGasPrices()).WithKeybase(f.Keybase).
		BuildAndSign(info.GetName(), ckeys.DefaultKeyPass, datagram)
}

// Request represents a request to the facuet
type Request struct {
	ChainID string `json:"chain-id"` // prevents sending funds if chain is not a testet
	Address string `json:"address"` // address of the requester
}

// GetAddress parses the string address from the request to an sdk.Address
func (r Request) GetAddress() (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(fr.Address)
}

func (r Request) Validate() error {
	if strings.TrimSpace(r.ChainID) == "" {
		return errors.New("chain-id cannot be blank")
	}

	_, err := r.GetAddress()
	return err
}
