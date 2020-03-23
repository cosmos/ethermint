package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"gopkg.in/yaml.v2"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var _ exported.Account = (*Account)(nil)
var _ exported.GenesisAccount = (*Account)(nil)

// ----------------------------------------------------------------------------
// Main Ethermint account
// ----------------------------------------------------------------------------

// // Account implements the auth.Account interface and embeds an
// // auth.BaseAccount type. It is compatible with the auth.AccountKeeper.
// type Account struct {
// 	*auth.BaseAccount

// 	// merkle root of the storage trie
// 	//
// 	// TODO: add back root if needed (marshalling is broken if not initializing)
// 	// Root ethcmn.Hash

// 	CodeHash []byte
// }

// ProtoAccount defines the prototype function for BaseAccount used for an
// AccountKeeper.
func ProtoAccount() exported.Account {
	return &Account{
		BaseAccount: &auth.BaseAccount{},
		CodeHash:    ethcrypto.Keccak256(nil),
	}
}

type ethermintAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        []byte         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	CodeHash      string         `json:"code_hash" yaml:"code_hash"`
}

// MarshalYAML returns the YAML representation of an account.
func (acc Account) MarshalYAML() (interface{}, error) {
	alias := ethermintAccountPretty{
		Address:       acc.Address,
		PubKey:        acc.PubKey,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

// MarshalJSON returns the JSON representation of an Account.
func (acc Account) MarshalJSON() ([]byte, error) {
	alias := ethermintAccountPretty{
		Address:       acc.Address,
		PubKey:        acc.PubKey,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		CodeHash:      ethcmn.Bytes2Hex(acc.CodeHash),
	}

	return json.Marshal(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into an Account.
func (acc *Account) UnmarshalJSON(bz []byte) error {
	acc.BaseAccount = &authtypes.BaseAccount{}
	var alias ethermintAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	if alias.PubKey != nil {
		pubk, err := tmamino.PubKeyFromBytes(alias.PubKey)
		if err != nil {
			return err
		}

		acc.BaseAccount.PubKey = pubk.Bytes()
	}

	acc.BaseAccount.Address = alias.Address
	acc.BaseAccount.AccountNumber = alias.AccountNumber
	acc.BaseAccount.Sequence = alias.Sequence
	acc.CodeHash = ethcmn.Hex2Bytes(alias.CodeHash)

	return nil
}
