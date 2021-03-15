#!/bin/bash

CHAINID="ethermint-1337"
MONIKER="localtestnet"

# localKey address 0x71c767afdadad36f1e6706b9c5d412ad52a9499c
VAL_KEY="localkey"
VAL_MNEMONIC="gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"

# user1 address 0x872622f88c502a7c3c534659b5ebdc79ec2c581b
USER1_KEY="user1"
USER1_MNEMONIC="copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom"

# user2 address 0xac78c8e82ce54b9630995b4a04aa68d7ca991a50
USER2_KEY="user2"
USER2_MNEMONIC="maximum display century economy unlock van census kite error heart snow filter midnight usage egg venture cash kick motor survey drastic edge muffin visual"

# remove existing daemon and client
rm -rf ~/.ethermint*

# Import keys from mnemonics
echo $VAL_MNEMONIC | ethermintd keys add $VAL_KEY --recover --keyring-backend test
echo $USER1_MNEMONIC | ethermintd keys add $USER1_KEY --recover --keyring-backend test
echo $USER2_MNEMONIC | ethermintd keys add $USER2_KEY --recover --keyring-backend test

ethermintd init $MONIKER --chain-id $CHAINID

# Allocate genesis accounts (cosmos formatted addresses)
ethermintd add-genesis-account $(ethermintd keys show $VAL_KEY -a --keyring-backend test) 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test
ethermintd add-genesis-account $(ethermintd keys show $USER1_KEY -a --keyring-backend test) 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test
ethermintd add-genesis-account $(ethermintd keys show $USER2_KEY -a --keyring-backend test) 1000000000000000000000aphoton,1000000000000000000stake --keyring-backend test

# Sign genesis transaction
ethermintd gentx $VAL_KEY 1000000000000000000stake --amount=1000000000000000000000aphoton --chain-id $CHAINID --keyring-backend test

# Collect genesis tx
ethermintd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
ethermintd validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
ethermintd start --pruning=nothing --rpc.unsafe --rpc-api "web3, eth, personal, net" --keyring-backend test --trace --log_level "info"
