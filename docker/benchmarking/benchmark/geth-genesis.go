package benchmark

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/urfave/cli"
)

type Genesis struct {
	Config struct {
		ChainID        int64  `json:"chainId"`
		Eip150Block    int64  `json:"eip150Block"`
		Eip150Hash     string `json:"eip150Hash"`
		Eip155Block    int64  `json:"eip155Block"`
		Eip158Block    int64  `json:"eip158Block"`
		HomesteadBlock int64  `json:"homesteadBlock"`
	} `json:"config"`
	Difficulty string      `json:"difficulty"`
	GasLimit   string      `json:"gasLimit"`
	Alloc      interface{} `json:"alloc"`
}
type Balance struct {
	Balance string `json:"balance"`
}

var (
	AddGenesis = cli.Command{
		Name:      "add-genesis-geth",
		ShortName: "ag",
		Usage:     "add geth account to genesis file",
		Action:    addAccountGeth,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "account, ac", Hidden: false, Usage: "add genesis account"},
			cli.IntFlag{Name: "amount, am", Hidden: false, Usage: "add balance to account"},
		},
	}
)

func addAccountGeth(ctx *cli.Context) error {
	st := ctx.String("account")
	r := regexp.MustCompile(`{(.+)?}`)
	res := r.FindStringSubmatch(st)
	addr := res[1]

	genesis, err := readGenesis()
	if err != nil {
		return err
	}

	alloc, ok := genesis.Alloc.(map[string]interface{})
	if ok {
		alloc[addr] = Balance{strconv.Itoa(ctx.Int("amount"))}
	}

	genesis.Alloc = alloc
	err = writeGenesis(genesis)
	if err != nil {
		return err
	}

	return nil
}

func readGenesis() (Genesis, error) {
	jsonFile, err := os.Open("bench-geth-genesis.json")
	if err != nil {
		return Genesis{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var genesis Genesis
	json.Unmarshal(byteValue, &genesis)
	return genesis, nil
}

func writeGenesis(data Genesis) error {
	file, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("bench-geth-genesis.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}
