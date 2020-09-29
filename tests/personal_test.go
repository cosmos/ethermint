package tests

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/require"
)

func TestPersonal_ListAccounts(t *testing.T) {
	rpcRes := call(t, "personal_listAccounts", []string{})

	var res []hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
}

func TestPersonal_NewAccount(t *testing.T) {
	rpcRes := call(t, "personal_newAccount", []string{"password"})
	var addr common.Address
	err := json.Unmarshal(rpcRes.Result, &addr)
	require.NoError(t, err)

	rpcRes = call(t, "personal_listAccounts", []string{})
	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
}

func TestPersonal_Sign(t *testing.T) {
	rpcRes := call(t, "personal_sign", []interface{}{hexutil.Bytes{0x88}, hexutil.Bytes(from), ""})

	var res hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 65, len(res))
	// TODO: check that signature is same as with geth, requires importing a key
}

func TestPersonal_ImportRawKey(t *testing.T) {
	privkey, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	// parse priv key to hex
	hexPriv := common.Bytes2Hex(ethcrypto.FromECDSA(privkey))
	rpcRes := call(t, "personal_importRawKey", []string{hexPriv, "password"})

	var res hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)

	addr := ethcrypto.PubkeyToAddress(privkey.PublicKey)
	resAddr := common.BytesToAddress(res)

	require.Equal(t, addr.String(), resAddr.String())
}

func TestPersonal_EcRecover(t *testing.T) {
	data := hexutil.Bytes{0x88}
	rpcRes := call(t, "personal_sign", []interface{}{data, hexutil.Bytes(from), ""})

	var res hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 65, len(res))

	rpcRes = call(t, "personal_ecRecover", []interface{}{data, res})
	var ecrecoverRes common.Address
	err = json.Unmarshal(rpcRes.Result, &ecrecoverRes)
	require.NoError(t, err)
	require.Equal(t, from, ecrecoverRes[:])
}

func TestPersonal_UnlockAccount(t *testing.T) {
	pswd := "nootwashere"
	rpcRes := call(t, "personal_newAccount", []string{pswd})
	var addr common.Address
	err := json.Unmarshal(rpcRes.Result, &addr)
	require.NoError(t, err)

	// try to sign, should be locked
	_, err = callWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, ""})
	require.Error(t, err)

	rpcRes = call(t, "personal_unlockAccount", []interface{}{addr, ""})
	var unlocked bool
	err = json.Unmarshal(rpcRes.Result, &unlocked)
	require.NoError(t, err)
	require.True(t, unlocked)

	// try to sign, should work now
	rpcRes = call(t, "personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, pswd})
	var res hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 65, len(res))
}

func TestPersonal_LockAccount(t *testing.T) {
	pswd := "nootwashere"
	rpcRes := call(t, "personal_newAccount", []string{pswd})
	var addr common.Address
	err := json.Unmarshal(rpcRes.Result, &addr)
	require.NoError(t, err)

	rpcRes = call(t, "personal_unlockAccount", []interface{}{addr, ""})
	var unlocked bool
	err = json.Unmarshal(rpcRes.Result, &unlocked)
	require.NoError(t, err)
	require.True(t, unlocked)

	rpcRes = call(t, "personal_lockAccount", []interface{}{addr})
	var locked bool
	err = json.Unmarshal(rpcRes.Result, &locked)
	require.NoError(t, err)
	require.True(t, locked)

	// try to sign, should be locked
	_, err = callWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, ""})
	require.Error(t, err)
}
