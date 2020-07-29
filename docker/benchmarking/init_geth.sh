#!/bin/bash

PW="1234567"
ACCTS=30

go get -d github.com/ethereum/go-ethereum
go install github.com/ethereum/go-ethereum/cmd/geth

geth --identity "geth-benchmark" init --datadir data bench-geth-genesis.json

# Generate 30 accounts
echo -e 'Generating Accounts\n'
for i in $(seq 1 $ACCTS)
do
  geth --datadir "data" account new --password <(echo $PW)
  # i dont know what the windows command would be to get the last created account thats what this does in windows
  for /f "delims=" %%i in ('geth --datadir "data" account list') do set account=%%i
  ./benchmarking ag -ac %account% -am 10000000
  # echo "Generated test$i account"
done

# # echo "starting geth node"
geth --identity "geth-benchmark" --http --http.port "8545" --http.corsdomain "*" --datadir "data" --port "30303" --nodiscover --http.api "db,eth,net,web3,personal,miner,admin" --networkid 1900 --nat "any" 2>> "gethd.log" &

# # echo "starting geth attach"
geth --datadir "data" attach http://localhost:8545 2>> "gethcli.log" &