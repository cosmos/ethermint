package faucet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tendermint/tendermint/libs/log"

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

// Faucet for sending tokens upon request during testnets.
type Faucet struct {
	// TODO: config
	config string
	// faucet keyring for signing transfers
	keyring keyring.Keybase
	// addresses can send requests every <defaultTimeout> duration
	defaultTimeout time.Duration
	// max amount
	capAmount sdk.Int
	// history of users and their funding timeouts
	timeouts map[string]time.Time

	// Ethermint codecs
	amino *sdkcdc.Codec
	cdc   *codec.Codec

	logger log.Logger
}

// New creates a new Faucet instance
func New(keyname string, timeout time.Duration, maxAmount sdk.Int) (Faucet, error) {
	keyring, err := keyring.NewKeyring("ethermint", "test", keysDir(homePath, src.ChainID), nil)
	if err != nil {
		return Faucet{}, err
	}

	amino := codec.MakeCodec(app.ModuleBasics)
	appCodec := codec.NewAppCodec(amino)

	return Faucet{
		keyring: keyring,
		defaultTimeout: timeout,
		capAmount:      maxAmount,
		timeouts:       make(map[string]time.Time),
		amino:          amino,
		cdc:            appCodec,
		logger:         log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}, nil
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
	lastRequest, ok := f.timeouts[address]
	if !ok {
		f.timeouts[address] = time.Now().UTC()
		return nil
	}

	sinceLastRequest := time.Since(lastRequest)
	if f.defaultTimeout > sinceLastRequest {
		wait := f.defaultTimeout - sinceLastRequest
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", address, f.timeout.String(), wait.String())
	}

	// user able to send funds since they have waited for period
	f.timeouts[address] = time.Now().UTC()
	return nil
}

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

	acc, err := auth.NewAccountRetriever(f.cdc, f).GetAccount(info.GetAddress())
	if err != nil {
		return nil, err
	}

	memo := "faucet transfer"

	return auth.NewTxBuilder(
		auth.DefaultTxEncoder(f.amino),
		acc.GetAccountNumber(),
		acc.GetSequence(), f.Gas, f.GasAdjustment, false, f.ChainID,
		memo, f.getGasPrices()).WithKeybase(f.keyring).
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
	_, err := r.GetAddress()
	return err
}
