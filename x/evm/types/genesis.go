package types

import (
	"fmt"
	"math/big"

	"github.com/cosmos/ethermint/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

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

// ValidateGenesis validates evm genesis config
func ValidateGenesis(data GenesisState) error {
	for _, acct := range data.Accounts {
		if len(acct.Address.Bytes()) == 0 {
			return fmt.Errorf("invalid GenesisAccount Error: Missing Address")
		}
		if acct.Balance == nil {
			return fmt.Errorf("invalid GenesisAccount Error: Missing Balance")
		}
	}
	return nil
}

// DefaultGenesisState sets default evm genesis config
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: []GenesisAccount{},
	}
}
