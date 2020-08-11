package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeyring(t *testing.T) {
	dir, cleanup := tests.NewTestCaseDir(t)
	mockIn := strings.NewReader("")
	t.Cleanup(cleanup)

	kr, err := keyring.NewKeyring("ethermint", keyring.BackendTest, dir, mockIn, EthSecp256k1Options()...)
	require.NoError(t, err)

	// fail in retrieving key
	info, err := kr.Get("foo")
	require.Error(t, err)
	require.Nil(t, info)

	mockIn.Reset("password\npassword\n")
	info, mnemonic, err := kr.CreateMnemonic("foo", keyring.English, sdk.FullFundraiserPath, EthSecp256k1)
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)
	require.Equal(t, "foo", info.GetName())
	require.Equal(t, "local", info.GetType().String())
	require.Equal(t, EthSecp256k1, info.GetAlgo())

	params := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	hdPath := params.String()

	bz, err := DeriveKey(mnemonic, keyring.DefaultBIP39Passphrase, hdPath, EthSecp256k1)
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	bz, err = DeriveKey(mnemonic, keyring.DefaultBIP39Passphrase, hdPath, keyring.Secp256k1)
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	bz, err = DeriveKey(mnemonic, keyring.DefaultBIP39Passphrase, hdPath, keyring.SigningAlgo(""))
	require.Error(t, err)
	require.Empty(t, bz)
}
