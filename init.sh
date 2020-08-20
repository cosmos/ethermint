#!/bin/bash

KEY="mykey"
CHAINID=8
MONIKER="localtestnet"

# remove existing daemon and dent
rm -rf ~/.ethermint*

make install

ethermintd config keyring-backend test

# Set up config for d
ethermintd config chain-id $CHAINID
ethermintd config output json
ethermintd config indent true
ethermintd config trust-node true

# if $KEY exists it should be deleted
ethermintd keys add $KEY --algo "eth_secp256k1"

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
ethermintd init $MONIKER --chain-id $CHAINID

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account $(ethermintd keys show $KEY -a) 1000000000000000000photon,1000000000000000000stake

# Sign genesis transaction
ethermintd gentx --name $KEY --keyring-backend test

# Collect genesis tx
ethermintd collect-gentxs

# Enable faucet
cat  $HOME/.ethermintd/config/genesis.json | jq '.app_state["faucet"]["enable_faucet"]=true' >  $HOME/.ethermintd/config/tmp_genesis.json && mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

echo -e '\n\ntestnet faucet enabled'
echo -e 'to transfer tokens to your account address use:'
echo -e "ethermintd tx faucet request 100photon --from $KEY\n"


# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Command to run the rest server in a different terminal/window
echo -e '\nrun the following command in a different terminal/window to run the REST server and JSON-RPC:'
echo -e "ethermintd rest-server --laddr \"tcp://localhost:8545\" --unlock-key $KEY --chain-id $CHAINID --trace\n"

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
ethermintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace
