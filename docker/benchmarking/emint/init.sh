#!/bin/bash

FKEY="faucet"
TESTKEY="test"
PW="12345678"
ACCTS=30
UNLOCKKEYS=""

CHAINID=123
MONIKER="localbenchmarktestnet"

# remove existing daemon and client
rm -rf ~/.ethermint*

pkill -f "ethermint*"

make install

ethermintcli config keyring-backend test

# Set up config for CLI
ethermintcli config chain-id $CHAINID
ethermintcli config output json
ethermintcli config indent true
ethermintcli config trust-node true

# if $KEY exists it should be deleted
ethermintcli keys add $FKEY --algo "eth_secp256k1"

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
ethermintd init $MONIKER --chain-id $CHAINID

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account $(ethermintcli keys show $FKEY -a) 1000000000000000000photon,1000000000000000000stake

# Generate 30 accounts
echo -e 'Generating Accounts\n'
for i in $(seq 1 $ACCTS)
do
  ethermintcli keys add $TESTKEY$i --algo "eth_secp256k1"
  ethermintd add-genesis-account $(ethermintcli keys show $TESTKEY$i -a) 1000000000000000000photon,1000000000000000000stake
  UNLOCKKEYS+=$TESTKEY$i","
  # echo "Generated test$i account"
done

# Sign genesis transaction
ethermintd gentx --name $FKEY --keyring-backend test

# Collect genesis tx
ethermintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to emintd.log
ethermintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace > emintd.log &

# Start the rest server with unlocked faucet key in background and log to emintcli.log 
ethermintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $FKEY,${UNLOCKKEYS::${#UNLOCKKEYS}-1} --chain-id $CHAINID --trace > emintcli.log &