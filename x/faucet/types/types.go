package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Mining is a struct that contains all the metadata of a mint
type Mining struct {
	Minter   sdk.AccAddress `json:"Minter"`
	LastTime int64          `json:"LastTime"`
	Total    sdk.Coin       `json:"Total"`
}

// NewMining returns a new Mining
func NewMining(minter sdk.AccAddress, coin sdk.Coin) Mining {
	return Mining{
		Minter:   minter,
		LastTime: 0,
		Total:    coin,
	}
}

// GetMinter get minter of mining
func (w Mining) GetMinter() sdk.AccAddress {
	return w.Minter
}

// implement fmt.Stringer
func (w Mining) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Minter: %s, Time: %s, Total: %s`, w.Minter, w.LastTime, w.Total))
}

type FaucetKey struct {
	Armor string `json:" armor"`
}

// NewFaucetKey create a instance
func NewFaucetKey(armor string) FaucetKey {
	return FaucetKey{
		Armor: armor,
	}
}

// implement fmt.Stringer
func (f FaucetKey) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Armor: %s`, f.Armor))
}
