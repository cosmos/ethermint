package faucet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagListenAddr = "listen"
)

func GetCmds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "commands for running faucet and requesting funds for testnet",
	}
	cmd.AddCommand(
		startCmd(),
		requestCmd(),
	)
	return cmd
}

// startCmd starts the faucet to listen on the given flagListenAddr address.
func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [key-name] [amount]",
		Short: "listens on a port for requests for tokens",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			faucet, err := config.Get(args[0])
			if err != nil {
				return err
			}
			info, err := faucet.keyring.Key(args[0])
			if err != nil {
				return err
			}

			// TODO: consider using int
			amount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			listenAddr, err := cmd.Flags().GetString(flagListenAddr)
			if err != nil {
				return err
			}

			r := mux.NewRouter()
			r.HandleFunc("/", faucet.Handler(info.GetAddress(), amount)).Methods("POST")

			server := &http.Server{
				Handler:      r,
				Addr:         listenAddr,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			return server.ListenAndServe()
		},
	}
	cmd.Flags().StringP(flagListenAddr, "l", "0.0.0.0:8000", "sets the faucet listener addresss")
	if err := viper.BindPFlag(flagListenAddr, cmd.Flags().Lookup(flagListenAddr)); err != nil {
		panic(err)
	}
	return cmd
}

func requestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [key]",
		Short: "request tokens from the faucet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			faucet, err := config.Get(args[0])
			if err != nil {
				return err
			}

			urlString, err := cmd.Flags().GetString(flagURL)
			if err != nil {
				return err
			}

			if urlString == "" {
				u, err := url.Parse(faucet.RPCAddr)
				if err != nil {
					return err
				}

				host, port, err := net.SplitHostPort(u.Host)
				if err != nil {
					return err
				}

				urlString = fmt.Sprintf("%s://%s:%d", u.Scheme, host, port)
			}

			var keyName string
			if len(args) == 2 {
				keyName = args[1]
			} else {
				keyName = faucet.Key
			}

			info, err := faucet.keyring.Key(keyName)
			if err != nil {
				return err
			}

			// send request using an sdk.AccAddress
			body, err := json.Marshal(Request{Address: info.GetAddress().String()})
			if err != nil {
				return err
			}

			resp, err := http.Post(urlString, "application/json", bytes.NewBuffer(body))
			if err != nil {
				return err
			}

			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			fmt.Println(string(respBody))
			return nil
		},
	}
	return urlFlag(cmd)
}

// // Service creates a faucet service for
// func Service() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "faucet [user] [home] [chain-id] [key-name] [amount]",
// 		Short: "faucet returns a sample faucet service file",
// 		Args:  cobra.ExactArgs(5),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			chain, err := config.Chains.Get(args[2])
// 			if err != nil {
// 				return err
// 			}
// 			_, err = chain.Keybase.Key(args[3])
// 			if err != nil {
// 				return err
// 			}
// 			_, err = sdk.ParseCoin(args[4])
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Printf(`[Unit]
// Description=faucet
// After=network.target
// [Service]
// Type=simple
// User=%s
// WorkingDirectory=%s
// ExecStart=%s/go/bin/rly testnets faucet %s %s %s
// Restart=on-failure
// RestartSec=3
// LimitNOFILE=4096
// [Install]
// WantedBy=multi-user.target
// `, args[0], args[1], args[1], args[2], args[3], args[4])
// 			return nil
// 		},
// 	}
// 	return cmd
// }
