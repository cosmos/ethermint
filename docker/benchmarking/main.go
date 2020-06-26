package main

import (
	"log"
	"os"

	"github.com/araskachoi/ethermint/docker/benchmarking/tx"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "benchmark"
	app.Usage = "Benchmarking suite for ethermint"
	app.Commands = []cli.Command{
		tx.SendTx,
		tx.GenerateAccts
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
