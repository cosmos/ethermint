package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/urfave/cli"
)

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

var (
	HOST = os.Getenv("HOST")

	SendTx = cli.Command{
		Name:   "sendtx",
		Usage:  "Command to send transactions",
		Action: sendTx,
		Flags:  []cli.Flag{},
	}
)

func getRandAcct(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func createRequest(method string, params interface{}) Request {
	return Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func call(method string, params interface{}) (*Response, error) {
	if HOST == "" {
		HOST = "http://localhost:8545"
	}

	req, err := json.Marshal(createRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *Response
	time.Sleep(1 * time.Second)
	/* #nosec */
	res, err := http.Post(HOST, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return rpcRes, nil
}

func getTransactionReceipt(hash hexutil.Bytes) (map[string]interface{}, error) {
	param := []string{hash.String()}
	rpcRes, err := call("eth_getTransactionReceipt", param)
	if err != nil {
		return nil, err
	}

	receipt := make(map[string]interface{})
	err = json.Unmarshal(rpcRes.Result, &receipt)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func waitForReceipt(hash hexutil.Bytes) (map[string]interface{}, error) {
	for i := 0; i < 600; i++ {
		receipt, err := getTransactionReceipt(hash)
		if receipt != nil {
			return receipt, err
		}

		time.Sleep(time.Second)
	}
	return nil, nil
}

// will get a list of addrs from some config file. these will be the addresses that are included in the genesis file.
func sendTx(ctx *cli.Context) error {

	rpcRes, err := call("eth_accounts", []string{})
	if err != nil {
		return err
	}

	var accts []string
	err = json.Unmarshal(rpcRes.Result, &accts)
	if err != nil {
		return err
	}

	if len(accts) == 0 {
		fmt.Println("no accounts available")
		return nil
	}

	value := "0x3B9ACA00"
	gasLimit := "0x5208"
	gasPrice := "0x15EF3C0"
	maxTx := 1000
	txs := 0

	ticker := time.NewTicker(time.Duration(100000) * time.Nanosecond)
	defer ticker.Stop()
	testDuration := time.NewTicker(time.Duration(10) * time.Second)
	defer testDuration.Stop()

	echan := make(chan error)

	hashes := []hexutil.Bytes{}
	receipts := []map[string]interface{}{}
	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			txs++
			if txs >= maxTx {
				wg.Wait()
				ticker.Stop()
				fmt.Println("hashes: ", hashes)
				fmt.Println("receipts: ", receipts)
				return nil
			}
			wg.Add(1)
			go func(e chan error) {
				fmt.Println(txs)

				from := accts[getRandAcct(0, len(accts)-1)]
				to := accts[getRandAcct(0, len(accts)-1)]

				if string(from) == string(to) {
					to = accts[getRandAcct(0, len(accts)-1)]
				}

				fmt.Println("from: ", from)
				fmt.Println("to: ", to)

				param := make([]map[string]string, 1)
				param[0] = make(map[string]string)

				param[0]["from"] = fmt.Sprintf("%s", from)
				param[0]["to"] = fmt.Sprintf("%s", to)
				param[0]["value"] = value
				param[0]["gasLimit"] = gasLimit
				param[0]["gasPrice"] = gasPrice

				fmt.Println(param)

				rpcTxRes, err := call("eth_sendTransaction", param)
				if err != nil {
					fmt.Println(err)
					echan <- err
				}

				var hash hexutil.Bytes
				err = json.Unmarshal(rpcTxRes.Result, &hash)
				if err != nil {
					fmt.Println(err)
					echan <- err
				}
				hashes = append(hashes, hash)

				receipt, err := waitForReceipt(hash)
				if err != nil {
					fmt.Println(err)
					echan <- err
				}
				receipts = append(receipts, receipt)

				wg.Done()
			}(echan)
		case <-testDuration.C:
			ticker.Stop()
			testDuration.Stop()
			return nil
		case err := <-echan:
			fmt.Printf("received err on channel:\n%v", err)
			return err
		}
	}
}
