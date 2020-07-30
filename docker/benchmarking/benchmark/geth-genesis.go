package benchmark

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

type Genesis struct {
	Config struct {
		ChainID             int64       `json:"chainId"`
		Eip150Block         int64       `json:"eip150Block"`
		Eip155Block         int64       `json:"eip155Block"`
		Eip158Block         int64       `json:"eip158Block"`
		HomesteadBlock      int64       `json:"homesteadBlock"`
		ByzantiumBlock      int64       `json:"byzantiumBlock"`
		ConstantinopleBlock int64       `json:"constantinopleBlock"`
		PetersburgBlock     int64       `json:"petersburgBlock"`
		Consensus           interface{} `json:"clique"`
	} `json:"config"`
	Difficulty string      `json:"difficulty"`
	GasLimit   string      `json:"gasLimit"`
	ExtraData  string      `json:"extraData"`
	Alloc      interface{} `json:"alloc"`
}
type Balance struct {
	Balance string `json:"balance"`
}

var (
	AddAcctGenesis = cli.Command{
		Name:      "geth-add-genesis-acct",
		ShortName: "gaa",
		Usage:     "add geth account to genesis file and allocate funds",
		Action:    addAccountGeth,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "account, ac", Hidden: false, Usage: "add genesis account"},
			cli.IntFlag{Name: "amount, am", Hidden: false, Usage: "add balance to account"},
		},
	}
	AddSignerGenesis = cli.Command{
		Name:      "geth-add-genesis-signer",
		ShortName: "gas",
		Usage:     "add geth account to genesis file as a signer",
		Action:    addSignerGeth,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "account, ac", Hidden: false, Usage: "add genesis account"},
		},
	}
)

func addSignerGeth(ctx *cli.Context) error {
	signer := ctx.String("account")

	genesis, err := readGenesis()
	if err != nil {
		return err
	}

	genesis.ExtraData = "0x0000000000000000000000000000000000000000000000000000000000000000" + signer + "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

	err = writeGenesis(genesis)
	if err != nil {
		return err
	}

	return nil
}

func addAccountGeth(ctx *cli.Context) error {
	addr := ctx.String("account")

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
	jsonFile, err := os.Open("genesis.json")
	if err != nil {
		jsonFile, err = os.Open("templ-genesis.json")
		if err != nil {
			return Genesis{}, err
		}
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

	err = ioutil.WriteFile("genesis.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}
