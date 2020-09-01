#!/bin/bash

KEY="mykey"
TESTKEY="test"
CHAINID=123321
MONIKER="localtestnet"

# stop and remove existing daemon and client data and process(es)
rm -rf ~/.ethermint*
pkill -f "ethermint*"

type "ethermintd" 2> /dev/null || make install
type "ethermintcli" 2> /dev/null || make install

ethermintcli config keyring-backend test

# Set up config for CLI
ethermintcli config chain-id $CHAINID
ethermintcli config output json
ethermintcli config indent true
ethermintcli config trust-node true

# if $KEY exists it should be deleted
ethermintcli keys add $KEY

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
ethermintd init $MONIKER --chain-id $CHAINID

cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["faucet"]["enable_faucet"]=true' >  $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account $(ethermintcli keys show $KEY -a) 1000000000000000000aphoton,1000000000000000000stake

# Sign genesis transaction
ethermintd gentx --name $KEY --keyring-backend test

# Collect genesis tx
ethermintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to file
ethermintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace > ethermintd.log &

sleep 1

# Start the rest server with unlocked faucet key in background and log to file
ethermintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $KEY --chain-id $CHAINID --trace > ethermintcli.log &

solc --abi contracts/counter/counter.sol --bin -o contracts/counter
mv contracts/counter/contracts_counter_counter_sol_Counter.abi contracts/counter/counter_sol.abi
mv contracts/counter/contracts_counter_counter_sol_Counter.bin contracts/counter/counter_sol.bin
abigen --bin=contracts/counter/counter_sol.bin --abi=contracts/counter/counter_sol.abi --pkg=main --out=contracts/counter/counter.go
sed -i '1s/^/0x/' contracts/counter/counter_sol.bin

# sleep 5

# TXHASH=$(curl --fail --silent -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"'$(curl --fail --silent -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545 | grep -o '\0x[^"]*' 2>&1)'", "data":"'$(cat contracts/abi_bin/contracts_counter_sol_Counter.bin)'"}],"id":1}' -H "Content-Type: application/json" http://localhost:8545 | grep -o '\0x[^"]*' 2>&1)

# echo $TXHASH

# sleep 5

# CONTRACTTX=$(curl --fail --silent -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["'$TXHASH'"],"id":1}' -H "Content-Type: application/json" http://localhost:8545)

# echo $CONTRACTTX

ACCT=$(curl --fail --silent -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545 | grep -o '\0x[^"]*' 2>&1)

echo $ACCT


curl -X POST --data '{"jsonrpc":"2.0","method":"personal_unlockAccount","params":["'$ACCT'", ""],"id":1}' -H "Content-Type: application/json" http://localhost:8545

PRIVKEY=$(ethermintcli keys unsafe-export-eth-key $KEY)

echo $PRIVKEY

## need to get the private key from the account in order to check this functionality.
cd contracts/counter && go get && go build && ./counter $ACCT
