package benchmark

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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

type resource struct {
	timestamp     int
	cpuPercentage float64
	ramPercentage float64
	process       string
}

var (
	HOST = os.Getenv("HOST")

	SendTx = cli.Command{
		Name:      "sendtx",
		ShortName: "s",
		Usage:     "Command to send transactions",
		Action:    sendTx,
		Flags: []cli.Flag{
			cli.IntFlag{Name: "duration, d", Value: 10000, Hidden: false, Usage: "test duration in seconds"},
			cli.IntFlag{Name: "txcount, c", Value: 100, Hidden: false, Usage: "test transaction count"},
			cli.StringFlag{Name: "protocol, p", Value: "emint", Hidden: false, Usage: "geth or emint"},
		},
	}
	Analyze = cli.Command{
		Name:      "analyze",
		ShortName: "a",
		Usage:     "Analyze the receipts.json file. Output will be the blocks and corresponding transactions included in those blocks.",
		Action:    analyze,
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "verbose, v", Hidden: false, Usage: "Shows the additional details of analysis."},
			cli.IntFlag{Name: "start, s", Usage: "Returns the metrics averaged from given start time."},
			cli.IntFlag{Name: "end, e", Usage: "Returns the metrics averaged until given end time."},
			cli.StringFlag{Name: "protocol, p", Value: "emint", Hidden: false, Usage: "geth or emint"},
		},
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
	time.Sleep(1000000 * time.Nanosecond)
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
	for i := 0; i < 10; i++ {
		receipt, err := getTransactionReceipt(hash)
		if receipt != nil {
			return receipt, err
		}

		time.Sleep(time.Second)
	}
	return nil, nil
}

func getAllReceipts(hashes []hexutil.Bytes) []map[string]interface{} {
	receipts := make([]map[string]interface{}, 0)
	for _, hash := range hashes {
		receipt, err := waitForReceipt(hash)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		receipts = append(receipts, receipt)
	}
	return receipts
}

func checkRepeats(list []string, item string) []string {
	exist := false
	for _, litem := range list {
		if litem == item {
			exist = true
		}
	}
	if !exist {
		list = append(list, item)
	}
	return list
}

func hexToInt(hexStr string) int {
	result, _ := strconv.ParseUint(strings.Replace(hexStr, "0x", "", -1), 16, 64)
	return int(result)
}

func average(resourcelist []float64, timestamplist []int, start, end int) float64 {
	var sum float64
	for i, val := range timestamplist {
		if val >= start && val <= end {
			sum += resourcelist[i]
		}
	}
	return sum / float64(end-start)
}

func calcTotalTPS(transactions, timestamps []int) float64 {
	sum := 0
	for _, val := range transactions {
		sum += val
	}
	return float64(sum) / float64(timestamps[len(timestamps)-2]-timestamps[len(timestamps)-1])
}

func calcRangedTps(transactions, timestamps []int, startTime, endTime int) float64 {
	if endTime > timestamps[len(timestamps)-2] {
		endTime = timestamps[len(timestamps)-2]
	}
	startBlockIndex := 0
	endBlockIndex := 0
	startTxIndex := 0
	if startTime >= endTime {
		return -1
	}
	for i, ts := range timestamps {
		if startTime == ts {
			startBlockIndex = i
		}
		if endTime == ts {
			endBlockIndex = i
		}
	}
	if startBlockIndex == len(timestamps)-1 {
		startTxIndex = 0
	} else {
		startTxIndex = startBlockIndex + 1
	}
	endTxIndex := endBlockIndex + 1

	sum := 0
	if endTxIndex == startTxIndex {
		sum = transactions[startTxIndex]
	} else {
		for i := startTxIndex; i <= endTxIndex-1; i++ {
			sum += transactions[i]
		}
	}

	return float64(sum) / float64(endTime-startTime)
}

func sendTx(ctx *cli.Context) error {
	log.Println(fmt.Sprintf("Starting transactions. Sending %d transactions, timeout %d seconds", ctx.Int("txcount"), ctx.Int("duration")))
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

	//remove facuet from list of accts
	accts = accts[1:]

	var nonces = make([]int, len(accts))

	value := "0x3B9ACA00"   //
	gasLimit := "0x5208"    //
	gasPrice := "0x15EF3C0" //
	txs := 0

	txTicker := time.NewTicker(time.Duration(500000) * time.Nanosecond)
	defer txTicker.Stop()
	testDuration := time.NewTicker(time.Duration(ctx.Int("duration")) * time.Second)
	defer testDuration.Stop()

	echan := make(chan error)

	hashes := []hexutil.Bytes{}
	nonceIncIndex := 0
	var wg sync.WaitGroup

	startTime := time.Now()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		select {
		case <-txTicker.C:
			txs++
			if txs >= ctx.Int("txcount")+1 {
				wg.Wait()
				endTime := time.Now()

				txTicker.Stop()
				testDuration.Stop()
				receipts := getAllReceipts(hashes)

				hashesf, err := os.Create(cwd + "/hashes.json")
				if err != nil {
					return err
				}
				hashesJSON, err := json.Marshal(hashes)
				if err != nil {
					return err
				}
				_, err = hashesf.Write(hashesJSON)
				if err != nil {
					return err
				}

				receiptsf, err := os.Create(cwd + "/receipts.json")
				if err != nil {
					return err
				}
				receiptsJSON, err := json.Marshal(receipts)
				if err != nil {
					return err
				}
				_, err = receiptsf.Write(receiptsJSON)
				if err != nil {
					return err
				}

				startTimef, err := os.Create(cwd + "/start.txt")
				if err != nil {
					return err
				}
				_, err = startTimef.Write([]byte(fmt.Sprintf("%d", startTime.Unix())))
				if err != nil {
					return err
				}

				endTimef, err := os.Create(cwd + "/end.txt")
				if err != nil {
					return err
				}
				_, err = endTimef.Write([]byte(fmt.Sprintf("%d", endTime.Unix())))
				if err != nil {
					return err
				}

				log.Println(fmt.Sprintf("Transactions completed. Duration: %d [ns]", endTime.UnixNano()-startTime.UnixNano()))
				log.Println(fmt.Sprintf("Tx Start Time: %d [unix], Tx End Time: %d [unix]", startTime.Unix(), endTime.Unix()))

				return nil
			}

			wg.Add(1)

			fromIndex := getRandAcct(0, len(accts)-1)
			nonceIncIndex = fromIndex
			from := accts[fromIndex]
			toIndex := getRandAcct(0, len(accts)-1)
			to := accts[toIndex]

			if from == to {
				to = accts[getRandAcct(0, len(accts)-1)]
			}

			param := make([]map[string]string, 1)
			param[0] = make(map[string]string)

			param[0]["from"] = from
			param[0]["to"] = to
			param[0]["value"] = value
			param[0]["gasLimit"] = gasLimit
			param[0]["gasPrice"] = gasPrice
			param[0]["nonce"] = "0x" + fmt.Sprintf("%x", nonces[fromIndex])

			rpcTxRes, err := call("eth_sendTransaction", param)
			if err != nil {
				log.Panic(err)
				echan <- err
			}

			var hash hexutil.Bytes
			err = json.Unmarshal(rpcTxRes.Result, &hash)
			if err != nil {
				fmt.Println(err)
				echan <- err
			}
			hashes = append(hashes, hash)

			nonces[nonceIncIndex]++

			wg.Done()

		case <-testDuration.C:
			txTicker.Stop()
			testDuration.Stop()
			return nil
		case err := <-echan:
			fmt.Printf("received err on channel:\n%v", err)
			return err
		}
	}
}

func analyze(ctx *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	receiptsf, err := ioutil.ReadFile(cwd + "/receipts.json")
	if err != nil {
		fmt.Println("Unable to locate receipts.json file. Please run the sendtx command to generate this file.")
		return err
	}
	var receipts []map[string]interface{}
	var transactions []int
	var totalTx int
	var resourceUsage []resource

	err = json.Unmarshal(receiptsf, &receipts)
	if err != nil {
		return err
	}

	blocks := []string{}
	timestamps := []int{}
	for _, receipt := range receipts {
		blockn := fmt.Sprintf("%s", receipt["blockNumber"])
		blocks = checkRepeats(blocks, blockn)
	}

	for _, block := range blocks {
		param := []interface{}{block}
		rpcResGetTx, err := call("eth_getBlockTransactionCountByNumber", param)
		if err != nil {
			return err
		}
		var txCount string
		err = json.Unmarshal(rpcResGetTx.Result, &txCount)
		if err != nil {
			return err
		}

		transactions = append(transactions, hexToInt(txCount))

		param = []interface{}{block, false}
		rpcResGetBlock, err := call("eth_getBlockByNumber", param)
		if err != nil {
			return err
		}
		jsonBlock := make(map[string]interface{})
		err = json.Unmarshal(rpcResGetBlock.Result, &jsonBlock)
		if err != nil {
			return err
		}

		timestamps = append(timestamps, hexToInt(fmt.Sprintf("%s", jsonBlock["timestamp"])))

		if ctx.Bool("verbose") {
			fmt.Println(jsonBlock)
		}
	}

	// get the block before the first blocks with transactions; for timestamp.
	param := []interface{}{"0x" + fmt.Sprintf("%x", hexToInt(blocks[0])-1), false}
	rpcResGetBlock, err := call("eth_getBlockByNumber", param)
	if err != nil {
		return err
	}
	jsonBlock := make(map[string]interface{})
	err = json.Unmarshal(rpcResGetBlock.Result, &jsonBlock)
	if err != nil {
		return err
	}
	timestamps = append(timestamps, hexToInt(fmt.Sprintf("%s", jsonBlock["timestamp"])))

	for _, tx := range transactions {
		totalTx += tx
	}

	// parse resource usage file
	resourcef, err := os.Open(cwd + "/resource.log")
	if err != nil {
		fmt.Println("Unable to locate resource.log file. Please run the sendtx command to generate this file.")
		return err
	}
	scanner := bufio.NewScanner(resourcef)
	re := regexp.MustCompile(`\s+`)

	var currentTimeStamp int
	var emintdCPUUsage []float64
	var emintdRAMUsage []float64
	var emintcliCPUUsage []float64
	var emintcliRAMUsage []float64
	var gethCPUUsage []float64
	var gethRAMUsage []float64
	var emintdTimestamps []int
	var emintcliTimestamps []int
	var gethTimestamps []int

	for scanner.Scan() {
		s := re.ReplaceAllString(scanner.Text(), " ")
		spl := strings.Split(s, " ")

		switch l := len(spl[0]); {
		case l == 10:
			currentTimeStamp, err = strconv.Atoi(spl[0])
			if err != nil {
				return err
			}
		case spl[0] == "" && l >= 4:
			cpu, err := strconv.ParseFloat(spl[1], 64)
			if err != nil {
				return err
			}
			ram, err := strconv.ParseFloat(spl[2], 64)
			if err != nil {
				return err
			}

			resourceUsage = append(resourceUsage,
				resource{
					timestamp:     currentTimeStamp,
					cpuPercentage: cpu,
					ramPercentage: ram,
					process:       spl[3],
				})
			switch process := spl[3]; {
			case process == "ethermintd":
				emintdCPUUsage = append(emintdCPUUsage, cpu)
				emintdRAMUsage = append(emintdRAMUsage, ram)
				emintdTimestamps = append(emintdTimestamps, currentTimeStamp)
			case process == "ethermintcli":
				emintcliCPUUsage = append(emintcliCPUUsage, cpu)
				emintcliRAMUsage = append(emintcliRAMUsage, ram)
				emintcliTimestamps = append(emintcliTimestamps, currentTimeStamp)
			case process == "geth":
				gethCPUUsage = append(gethCPUUsage, cpu)
				gethRAMUsage = append(gethRAMUsage, ram)
				gethTimestamps = append(gethTimestamps, currentTimeStamp)
			}
		case spl[0] != "":
			cpu, err := strconv.ParseFloat(spl[0], 64)
			if err != nil {
				return err
			}
			ram, err := strconv.ParseFloat(spl[1], 64)
			if err != nil {
				return err
			}
			resourceUsage = append(resourceUsage,
				resource{
					timestamp:     currentTimeStamp,
					cpuPercentage: cpu,
					ramPercentage: ram,
					process:       spl[2],
				})
			switch process := spl[2]; {
			case process == "ethermintd":
				emintdCPUUsage = append(emintdCPUUsage, cpu)
				emintdRAMUsage = append(emintdRAMUsage, ram)
				emintdTimestamps = append(emintdTimestamps, currentTimeStamp)
			case process == "ethermintcli":
				emintcliCPUUsage = append(emintcliCPUUsage, cpu)
				emintcliRAMUsage = append(emintcliRAMUsage, ram)
				emintcliTimestamps = append(emintcliTimestamps, currentTimeStamp)
			case process == "geth":
				gethCPUUsage = append(gethCPUUsage, cpu)
				gethRAMUsage = append(gethRAMUsage, ram)
				gethTimestamps = append(gethTimestamps, currentTimeStamp)
			}
		}
	}

	if ctx.Bool("verbose") {
		fmt.Printf("Resource Usage: %+v\n", resourceUsage)
	}
	fmt.Println("-------------------------")
	fmt.Println("Blocks with Tx: ", blocks)
	fmt.Println("Block Timestamps: ", timestamps) //last two timestamps: first-1 block timestamp
	fmt.Println("Transactions: ", transactions)
	fmt.Println("Total Transactions: ", totalTx)
	fmt.Println("Total TX per second (TPS): ", calcTotalTPS(transactions, timestamps))
	fmt.Println("-------------------------")
	if ctx.Int("start") > 0 && ctx.Int("end") > 0 {
		if ctx.String("protocol") == "emint" {
			fmt.Println("Average CPU Usage [emintd]: ", average(emintdCPUUsage, emintdTimestamps, ctx.Int("start"), ctx.Int("end")))
			fmt.Println("Average RAM Usage [emintd]: ", average(emintdRAMUsage, emintdTimestamps, ctx.Int("start"), ctx.Int("end")))
			fmt.Println("Average CPU Usage [emintcli]: ", average(emintcliCPUUsage, emintcliTimestamps, ctx.Int("start"), ctx.Int("end")))
			fmt.Println("Average RAM Usage [emintcli]: ", average(emintcliRAMUsage, emintcliTimestamps, ctx.Int("start"), ctx.Int("end")))
		} else {
			fmt.Println("Average CPU Usage [geth]: ", average(gethCPUUsage, gethTimestamps, ctx.Int("start"), ctx.Int("end")))
			fmt.Println("Average RAM Usage [geth]: ", average(gethRAMUsage, gethTimestamps, ctx.Int("start"), ctx.Int("end")))
		}

		fmt.Println(fmt.Sprintf("Ranged TX per second (TPS) [start:%d; end%d]: ", ctx.Int("start"), ctx.Int("end")), calcRangedTps(transactions, timestamps, ctx.Int("start"), ctx.Int("end")))
	}

	return nil
}
