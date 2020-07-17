#!/bin/bash

FKEY="faucet"
TESTKEY="test"

CHAINID=123
MONIKER="localbenchmarktestnet"

# remove existing daemon and client
rm -rf ~/.emint*
pkill -f "emint*"

type "emintd" 2> /dev/null || make install
type "emintcli" 2> /dev/null || make install

emintcli config keyring-backend test

# Set up config for CLI
emintcli config chain-id $CHAINID
emintcli config output json
emintcli config indent true
emintcli config trust-node true

# if $KEY exists it should be deleted
emintcli keys add $FKEY

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
emintd init $MONIKER --chain-id $CHAINID

# Use a custom genesis with pre-generated keys
cp ./benchmark_ethmint_genesis.json $HOME/.emintd/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
emintd add-genesis-account $(emintcli keys show $FKEY -a) 1000000000000000000photon,1000000000000000000stake

# Sign genesis transaction
emintd gentx --name $FKEY --keyring-backend test

# Collect genesis tx
emintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
emintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to emintd.log
emintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace > emintd.log &

sleep 1

# Start the rest server with unlocked faucet key in background and log to emintcli.log 
emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $FKEY --chain-id $CHAINID --trace > emintcli.log &

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

## need to get the private key from the account in order to check this functionality.
# cd contracts/counter && go get && go build && ./counter $ACCT