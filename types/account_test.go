package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	emintcrypto "github.com/cosmos/ethermint/crypto"
)

func init() {
	tmamino.RegisterKeyType(emintcrypto.PubKeySecp256k1{}, emintcrypto.PubKeyAminoName)
	tmamino.RegisterKeyType(emintcrypto.PrivKeySecp256k1{}, emintcrypto.PrivKeyAminoName)
}

func TestEthermintAccountJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	balance := sdk.NewCoins(sdk.NewCoin(DenomDefault, sdk.OneInt()))
	baseAcc := auth.NewBaseAccount(addr, balance, pubkey, 10, 50)
	ethAcc := EthAccount{BaseAccount: baseAcc, CodeHash: []byte{1, 2}}

	bz, err := json.Marshal(ethAcc)
	require.NoError(t, err)

	bz1, err := ethAcc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a EthAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, ethAcc.String(), a.String())
	require.Equal(t, ethAcc.PubKey, a.PubKey)
}

func TestEthermintPubKeyJSON(t *testing.T) {
	privkey, err := emintcrypto.GenerateKey()
	require.NoError(t, err)
	bz := privkey.PubKey().Bytes()

	pubk, err := tmamino.PubKeyFromBytes(bz)
	require.NoError(t, err)
	require.Equal(t, pubk, privkey.PubKey())
}

func TestSecpPubKeyJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	bz := pubkey.Bytes()

	pubk, err := tmamino.PubKeyFromBytes(bz)
	require.NoError(t, err)
	require.Equal(t, pubk, pubkey)
}
