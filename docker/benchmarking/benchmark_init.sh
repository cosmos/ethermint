#!/bin/bash

FKEY="faucet"
TESTKEY="test"
PW="12345678"
ACCTS=30
UNLOCKKEYS=""

CHAINID=123
MONIKER="localbenchmarktestnet"

# remove existing daemon and client
rm -rf ~/.emint*

pkill -f "emint*"

make install

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

# Generate 30 accounts
for i in $(seq 1 $ACCTS)
do
  emintcli keys add $TESTKEY$i
  emintd add-genesis-account $(emintcli keys show $TESTKEY$i -a) 1000000000000000000photon,1000000000000000000stake
  UNLOCKKEYS+=$TESTKEY$i","
  echo "Generated test$i account"
done

# Sign genesis transaction
emintd gentx --name $FKEY --keyring-backend test

# Collect genesis tx
emintd collect-gentxs

# Enable faucet
# cat  $HOME/.emintd/config/genesis.json | jq '.app_state["faucet"]["enable_faucet"]=true' >  $HOME/.emintd/config/tmp_genesis.json && mv $HOME/.emintd/config/tmp_genesis.json $HOME/.emintd/config/genesis.json

echo -e '\n\ntestnet faucet enabled'
echo -e 'to transfer tokens to your account address use:'
echo -e "emintcli tx faucet request 100photon --from $FKEY\n"

# Run this to ensure everything worked and that the genesis file is setup correctly
emintd validate-genesis

# Command to run the rest server in a different terminal/window
echo -e '\nrun the following command in a different terminal/window to run the REST server and JSON-RPC:'
echo -e "emintcli rest-server --laddr \"tcp://localhost:8545\" --unlock-key $FKEY --chain-id $CHAINID --trace\n"

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to emintd.log
emintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace > emintd.log &

# Start the rest server with unlocked faucet key in background and log to emintcli.log 
emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $FKEY,${UNLOCKKEYS::${#UNLOCKKEYS}-1} --chain-id $CHAINID --trace > emintcli.log &