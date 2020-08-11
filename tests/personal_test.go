package tests

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/stretchr/testify/require"
)

func TestPersonal_ListAccounts(t *testing.T) {
	rpcRes := call(t, "personal_listAccounts", []string{})

	var res []hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
}

func TestPersonal_Sign(t *testing.T) {
	addr := getAddress(t)
	rpcRes := call(t, "personal_sign", []interface{}{hexutil.Bytes{0x88}, hexutil.Bytes(addr), ""})

	var res hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 65, len(res))
	// TODO: check that signature is same as with geth, requires importing a key
}

func TestPersonal_EcRecover(t *testing.T) {
	addr := hexutil.Bytes(getAddress(t))
	data := hexutil.Bytes{0x88}
	rpcRes := call(t, "personal_sign", []interface{}{data, addr, ""})

	var res hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 65, len(res))

	rpcRes = call(t, "personal_ecRecover", []interface{}{data, res})
	var ecrecoverRes common.Address
	err = json.Unmarshal(rpcRes.Result, &ecrecoverRes)
	require.NoError(t, err)
	require.Equal(t, []byte(addr), ecrecoverRes[:])
}

func TestPersonal_NewAccount(t *testing.T) {
	rpcRes := call(t, "personal_newAccount", []string{""})
	var addr common.Address
	err := json.Unmarshal(rpcRes.Result, &addr)
	require.NoError(t, err)
	t.Log(addr.Hex())

	rpcRes = call(t, "personal_listAccounts", []string{})
	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
}
