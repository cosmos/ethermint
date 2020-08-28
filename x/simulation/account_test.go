package simulation_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdksim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/ethermint/x/simulation"
)

func TestRandomAccounts(t *testing.T) {
	t.Parallel()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"0-accounts", 0, 0},
		{"0-accounts", 1, 1},
		{"0-accounts", 1_000, 1_000},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := simulation.RandomAccounts(tt.n)
			require.Equal(t, tt.want, len(got))
			if tt.n == 0 {
				return
			}
			acc, i := sdksim.RandomAcc(r, got)
			require.True(t, acc.Equals(got[i]))
			accFound, found := sdksim.FindAccount(got, acc.Address)
			require.True(t, found)
			require.True(t, accFound.Equals(acc))
		})
	}
}

func TestFindAccountEmptySlice(t *testing.T) {
	t.Parallel()
	accs := simulation.RandomAccounts(1)
	require.Equal(t, 1, len(accs))
	acc, found := sdksim.FindAccount(nil, accs[0].Address)
	require.False(t, found)
	require.Nil(t, acc.Address)
	require.Nil(t, acc.PrivKey)
	require.Nil(t, acc.PubKey)
}
