#!/bin/bash

echo "init ethermint..."
make install &&
rm -rf ~/.emint* &&
emintd init moniker --chain-id $1 &&
emintcli config chain-id $1 &&
emintcli config output json &&
emintcli config indent true &&
emintcli config trust-node true &&
emintcli keys add $2
emintcli emintkeys add $3
emintd add-genesis-account $(emintcli keys show $2 -a) 1000photon,100000000stake &&
emintd add-genesis-account $(emintcli emintkeys show $3 -a) 100000000photon,100000000stake &&
emintd gentx --name $2 &&
emintd collect-gentxs &&
emintd validate-genesis &&
emintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info,*:error"