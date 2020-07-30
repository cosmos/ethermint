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

# Allocate genesis accounts (cosmos formatted addresses)
emintd add-genesis-account $(emintcli keys show $FKEY -a) 1000000000000000000photon,1000000000000000000stake

# Generate 30 accounts
echo -e 'Generating Accounts\n'
for i in $(seq 1 $ACCTS)
do
  emintcli keys add $TESTKEY$i
  emintd add-genesis-account $(emintcli keys show $TESTKEY$i -a) 1000000000000000000photon,1000000000000000000stake
  UNLOCKKEYS+=$TESTKEY$i","
  # echo "Generated test$i account"
done

# Sign genesis transaction
emintd gentx --name $FKEY --keyring-backend test

# Collect genesis tx
emintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
emintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to emintd.log
emintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace > emintd.log &

# Start the rest server with unlocked faucet key in background and log to emintcli.log 
emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $FKEY,${UNLOCKKEYS::${#UNLOCKKEYS}-1} --chain-id $CHAINID --trace > emintcli.log &