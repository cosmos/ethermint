package types

import (
	"fmt"
	"time"
)

// GenesisState defines the application's genesis state. It contains all the
// information required and accounts to initialize the blockchain.
type GenesisState struct {
	// enable faucet funding
	EnableFaucet bool `json:"enable_faucet" yaml:"enable_faucet"`
	// addresses can send requests every <Timeout> duration
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// max total amount to be funded by the faucet
	FaucetCap int64 `json:"faucet_cap" yaml:"faucet_cap"`
	// max amount per request (i.e sum of all requested coin amounts).
	MaxAmountPerRequest int64 `json:"max_amount_per_request" yaml:"max_amount_per_request"`
}

// ValidateGenesis validates faucet genesis files
func (gs GenesisState) Validate error {
	if gs.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative: %s", gs.Timeout)
	}
	if gs.FaucetCap < 0 {
		return fmt.Errorf("faucet cap cannot be negative: %d", gs.FaucetCap)
	}
	if gs.MaxAmountPerRequest < 0 {
		return fmt.Errorf("max amount per request cannot be negative: %d", gs.MaxAmountPerRequest)
	}
	return nil
}

// DefaultGenesisState sets default evm genesis config
func DefaultGenesisState() GenesisState {
	return GenesisState{
		EnableFaucet: false,
		Timeout:      time.Hour,
		FaucetCap:    10 ^ 9, // 1B max amount to be funded by the faucet
		MaxAmountPerRequest: 1000,
	}
}
