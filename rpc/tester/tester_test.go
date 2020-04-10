// This is a test utility for Ethermint's Web3 JSON-RPC services.
//
// To run these tests please first ensure you have the emintd running
// and have started the RPC service with `emintcl rest-server`.
//
// You can configure the desired port (or host) below.

package tester

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/cosmos/ethermint/version"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const (
	host          = "localhost"
	port          = 8545
	addrA         = "0xc94770007dda54cF92009BFF0dE90c06F603a09f"
	addrAStoreKey = 0
)

var addr = fmt.Sprintf("http://%s:%d", host, port)

type Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Error  *RPCError       `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
}

func createRequest(method string, params interface{}) Request {
	return Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func call(t *testing.T, method string, params interface{}) (*Response, error) {
	req, err := json.Marshal(createRequest(method, params))
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", req)

	/* #nosec */
	res, err := http.Post(addr, "application/json", bytes.NewBuffer(req))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)

	decoder := json.NewDecoder(res.Body)
	var rpcRes *Response
	err = decoder.Decode(&rpcRes)
	if err != nil {
		t.Fatal(err)
	}

	if rpcRes.Error != nil {
		t.Fatal(errors.New(rpcRes.Error.Message))
	}

	err = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	return rpcRes, nil
}

func TestEth_protocolVersion(t *testing.T) {
	expectedRes := hexutil.Uint(version.ProtocolVersion)

	rpcRes, err := call(t, "eth_protocolVersion", []string{})
	require.NoError(t, err)

	var res hexutil.Uint
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)

	t.Logf("Got protocol version: %s\n", res.String())
	require.Equal(t, expectedRes, res, "expected: %s got: %s\n", expectedRes.String(), rpcRes.Result)
}

func TestEth_blockNumber(t *testing.T) {
	rpcRes, err := call(t, "eth_blockNumber", []string{})
	require.NoError(t, err)

	var res hexutil.Uint64
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)

	t.Logf("Got block number: %s\n", res.String())
}

func TestEth_GetBalance(t *testing.T) {
	rpcRes, err := call(t, "eth_getBalance", []string{addrA, "0x0"})
	require.NoError(t, err)

	var res hexutil.Big
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)

	t.Logf("Got balance %s for %s\n", res.String(), addrA)

	// 0 if x == y; where x is res, y is 0
	if res.ToInt().Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected balance: %d, got: %s", 0, res.String())
	}
}

func TestEth_GetStorageAt(t *testing.T) {
	expectedRes := hexutil.Bytes{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	rpcRes, err := call(t, "eth_getStorageAt", []string{addrA, string(addrAStoreKey), "0x0"})
	require.NoError(t, err)

	var storage hexutil.Bytes
	err = storage.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)

	t.Logf("Got value [%X] for %s with key %X\n", storage, addrA, addrAStoreKey)

	require.True(t, bytes.Equal(storage, expectedRes), "expected: %d (%d bytes) got: %d (%d bytes)", expectedRes, len(expectedRes), storage, len(storage))
}

func TestEth_GetCode(t *testing.T) {
	expectedRes := hexutil.Bytes{}
	rpcRes, err := call(t, "eth_getCode", []string{addrA, "0x0"})
	require.NoError(t, err)

	var code hexutil.Bytes
	err = code.UnmarshalJSON(rpcRes.Result)

	require.NoError(t, err)

	t.Logf("Got code [%X] for %s\n", code, addrA)
	require.True(t, bytes.Equal(expectedRes, code), "expected: %X got: %X", expectedRes, code)
}

func getAddress(t *testing.T) []byte {
	rpcRes, err := call(t, "eth_accounts", []string{})
	require.NoError(t, err)

	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)

	t.Logf("Account: %s", res[0])
	return res[0]
}

func TestEth_SendTransaction(t *testing.T) {
	from := getAddress(t)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["data"] = "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029"

	rpcRes, err := call(t, "eth_sendTransaction", param)
	require.NoError(t, err)

	var hash hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &hash)
	require.NoError(t, err)

	t.Log(hash)
}

func TestEth_NewFilter(t *testing.T) {
	param := make([]map[string][]string, 1)
	param[0] = make(map[string][]string)
	param[0]["topics"] = []string{"0x0000000000000000000000000000000000000000000000000000000012341234"}
	rpcRes, err := call(t, "eth_newFilter", param)
	require.NoError(t, err)

	var ID hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &ID)
	require.NoError(t, err)
}

func TestEth_NewBlockFilter(t *testing.T) {
	rpcRes, err := call(t, "eth_newBlockFilter", []string{})
	require.NoError(t, err)

	var ID hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &ID)
	require.NoError(t, err)
}

func TestEth_GetFilterChanges_NoLogs(t *testing.T) {
	param := make([]map[string][]string, 1)
	param[0] = make(map[string][]string)
	param[0]["topics"] = []string{}
	rpcRes, err := call(t, "eth_newFilter", param)
	require.NoError(t, err)

	var ID hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &ID)
	require.NoError(t, err)

	changesRes, err := call(t, "eth_getFilterChanges", []string{ID.String()})
	require.NoError(t, err)

	var logs []*ethtypes.Log
	err = json.Unmarshal(changesRes.Result, &logs)
	require.NoError(t, err)
}

func TestEth_GetFilterChanges_WrongID(t *testing.T) {
	changesRes, err := call(t, "eth_getFilterChanges", []string{"0x1122334400000077"})
	require.NoError(t, err)

	var logs []*ethtypes.Log
	err = json.Unmarshal(changesRes.Result, &logs)
	require.NotNil(t, err)
}

// deployTestContract deploys a contract that emits an event in the constructor
func deployTestContract(t *testing.T) hexutil.Bytes {
	from := getAddress(t)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["data"] = "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029"

	rpcRes, err := call(t, "eth_sendTransaction", param)
	require.NoError(t, err)

	var hash hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &hash)
	require.NoError(t, err)

	return hash
}

func TestEth_GetTransactionReceipt(t *testing.T) {
	hash := deployTestContract(t)
	time.Sleep(time.Second)

	t.Log(hash)

	param := []string{hash.String()}
	rpcRes, err := call(t, "eth_getTransactionReceipt", param)
	require.NoError(t, err)

	t.Log(rpcRes.Result)

	recpt := new(ethtypes.Receipt)
	err = json.Unmarshal(rpcRes.Result, recpt)
	require.NoError(t, err)

	t.Log(recpt)
}

func TestEth_GetTxLogs(t *testing.T) {
	deployTestContract(t)

	// time.Sleep(time.Second*2)

	// param := []string{hash.String()}
	// rpcRes, err := call(t, "eth_getTxLogs", param)
	// require.NoError(t, err)

	// t.Log(rpcRes.Result)

	// logs := new([]*ethtypes.Log)
	// err = json.Unmarshal(rpcRes.Result, logs)
	// require.NoError(t, err)

	// t.Log(logs)
}

func TestEth_GetFilterChanges_NoParams(t *testing.T) {
	// rpcRes, err := call(t, "eth_blockNumber", []string{})
	// require.NoError(t, err)

	// var res hexutil.Uint64
	// err = res.UnmarshalJSON(rpcRes.Result)
	// require.NoError(t, err)

	//t.Log(res)

	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	param[0]["topics"] = []string{}
	param[0]["fromBlock"] = "0x1" //res.String()
	rpcRes, err := call(t, "eth_newFilter", param)
	require.NoError(t, err)

	//fmt.Println(rpcRes.Result)

	var ID hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &ID)
	require.NoError(t, err)

	// deploy contract, emitting some event
	_ = deployTestContract(t)
	time.Sleep(time.Second)

	// get filter changes
	changesRes, err := call(t, "eth_getFilterChanges", []string{ID.String()})
	require.NoError(t, err)

	var logs []*ethtypes.Log
	err = json.Unmarshal(changesRes.Result, &logs)
	require.NoError(t, err)

	t.Log(logs)
}
