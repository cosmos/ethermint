package faucet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	sdkcdc "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	ethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/codec"
)

// Faucet represents the necessary data for connecting to and indentifying a Faucet and its counterparites
type Faucet struct {
	keyring keyring.Keybase
	// addresses can send requests every <timeout> duration
	timeout time.Duration
	// stores faucet addresses that have been used reciently
	faucetAddrs map[string]time.Time

	// ethermint codecs
	amino *sdkcdc.Codec
	cdc   *codec.Codec
}

// New creates a new Faucet instance
func New(timeout time.Duration) Faucet {
	cdc := codec.MakeCodec(app.ModuleBasics)
	appCodec := codec.NewAppCodec(cdc)

	return Faucet{
		timeout:     timeout,
		faucetAddrs: make(map[string]time.Time),
		amino:       cdc,
		cdc:         appCodec,
	}
}

// Handler listens for addresses
func (f Faucet) Handler(from sdk.AccAddress, amount sdk.Coin) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %s", err))
			return
		}

		defer r.Body.Close()

		if err := req.Validate(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := f.rateLimit(req.Address); err != nil {
			rest.WriteErrorResponse(w, http.StatusTooManyRequests, err.Error())
			return
		}

		// error validated on req.Validate
		to, _ := req.GetAddress()

		if status, err := f.transfer(from, to, amount); err != nil {
			rest.WriteErrorResponse(w, status, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// _, _ = w.Write(req)
	}
}

// TODO:
// - introduce a counter for the number of requests within <timeout>
// - cap amount sent to address
func (f Faucet) rateLimit(address string) error {
	// first time requester, can send request
	lastRequest, ok := f.faucetAddrs[address]
	if !ok {
		f.faucetAddrs[address] = time.Now().UTC()
		return nil
	}

	sinceLastRequest := time.Since(lastRequest)
	if f.timeout > sinceLastRequest {
		wait := f.timeout - sinceLastRequest
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", address, f.timeout.String(), wait.String())
	}

	// user able to send funds since they have waited for period
	f.faucetAddrs[address] = time.Now().UTC()
	return nil
}

// TODO: prob use an sdk.Int instead of coins.
func (f Faucet) transfer(from, to sdk.AccAddress, amount sdk.Coin) (int, error) {
	_, err := f.keyring.GetByAddress(from)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_ = bank.NewMsgSend(from, to, sdk.NewCoins(amount))
	// res, err := f.SendMsgWithKey(msg, info.GetName())
	// if err != nil {
	// 	return http.StatusBadRequest, err
	// }

	return http.StatusOK, nil
}

// buildAndSignTxWithKey allows the user to specify which relayer key will sign the message
func (f Faucet) buildAndSignTxWithKey(msgs []sdk.Msg, keyName string) ([]byte, error) {
	// Fetch account and sequence numbers for the account
	info, err := f.keyring.Get(keyName)
	if err != nil {
		return nil, err
	}

	// acc, err := auth.NewAccountRetriever(f.Cdc, f).GetAccount(info.GetAddress())
	// if err != nil {
	// 	return nil, err
	// }

	memo := "faucet transfer"

	return auth.NewTxBuilder(
		auth.DefaultTxEncoder(f.amino),
		acc.GetAccountNumber(),
		acc.GetSequence(), f.Gas, f.GasAdjustment, false, f.ChainID,
		memo, f.getGasPrices()).WithKeybase(f.Keybase).
		BuildAndSign(info.GetName(), ckeys.DefaultKeyPass, msgs)
}

// Request represents a request to the facuet
type Request struct {
	Address string `json:"address"` // cosmos or ethereum address of the requester
}

// GetAddress parses the string address from the request to an sdk.Address
func (r Request) GetAddress() (sdk.AccAddress, error) {
	if strings.HasPrefix(r.Address, "0x") {
		ethAddress := ethcmn.HexToAddress(r.Address)
		return sdk.AccAddress(ethAddress.Bytes()), nil
	}

	return sdk.AccAddressFromBech32(r.Address)
}

// Validate validates a request fields
func (r Request) Validate() error {
	// if strings.TrimSpace(r.ChainID) == "" {
	// 	return errors.New("chain-id cannot be blank")
	// }

	_, err := r.GetAddress()
	return err
}
