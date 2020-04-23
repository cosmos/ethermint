package types

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/cosmos/ethermint/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

var zeroAddrBytes = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type (
	// GenesisState defines the application's genesis state. It contains all the
	// information required and accounts to initialize the blockchain.
	GenesisState struct {
		Accounts []GenesisAccount `json:"accounts"`
	}

	// GenesisAccount defines an account to be initialized in the genesis state.
	GenesisAccount struct {
		Address ethcmn.Address `json:"address"`
		Balance *big.Int       `json:"balance"`
		Code    []byte         `json:"code,omitempty"`
		Storage types.Storage  `json:"storage,omitempty"`
	}
)

// DefaultGenesisState sets default evm genesis config
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: []GenesisAccount{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, acc := range gs.Accounts {
		if bytes.Equal(acc.Address.Bytes(), zeroAddrBytes) {
			return errors.New("invalid GenesisAccount: address cannot be empty")
		}
		if acc.Balance == nil {
			return errors.New("invalid GenesisAccount: balance cannot be empty")
		}
	}
	return nil
}
