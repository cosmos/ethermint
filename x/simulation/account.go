package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/cosmos/ethermint/crypto"
)

// RandomAccounts generates n random accounts
func RandomAccounts(n int) []simulation.Account {
	accs := make([]simulation.Account, n)
	for i := 0; i < n; i++ {
		var err error
		accs[i].PrivKey, err = crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())
	}

	return accs
}
