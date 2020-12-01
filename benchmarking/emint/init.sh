#!/bin/bash

FKEY="mykey"
TESTKEY="test"
PW="12345678"
ACCTS=30
UNLOCKKEYS=""

CHAINID="ethermint-127001"
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
ethermintcli keys add $FKEY

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
ethermintd init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to aphoton
cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aphoton"' > $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json
cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aphoton"' > $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json
cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aphoton"' > $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json
cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aphoton"' > $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

# Enable faucet
cat $HOME/.ethermintd/config/genesis.json | jq '.app_state["faucet"]["enable_faucet"]=true' >  $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account $(ethermintcli keys show $FKEY -a) 100000000000000000000aphoton

# Generate 30 accounts
echo -e 'Generating Accounts\n'
for i in $(seq 1 $ACCTS)
do
  ethermintcli keys add $TESTKEY$i
  ethermintd add-genesis-account $(ethermintcli keys show $TESTKEY$i -a) 100000000000000000000aphoton
  UNLOCKKEYS+=$TESTKEY$i","
  # echo "Generated test$i account"
done

# Sign genesis transaction
ethermintd gentx --name $FKEY --amount=1000000000000000000aphoton --keyring-backend test

# Collect genesis tx
ethermintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to emintd.log
ethermintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace > emintd.log &

sleep 1

# Start the rest server with unlocked faucet key in background and log to emintcli.log 
ethermintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $FKEY,${UNLOCKKEYS::${#UNLOCKKEYS}-1} --chain-id $CHAINID --trace > emintcli.log &