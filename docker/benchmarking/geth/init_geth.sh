#!/bin/bash

PW="1234567"
ACCTS=3
UNLOCKACCTS=""

go get -d github.com/ethereum/go-ethereum
go install github.com/ethereum/go-ethereum/cmd/geth

# Generate 30 accounts
echo -e 'Generating Accounts\n'
for i in $(seq 1 $ACCTS)
do
  geth --datadir "data" account new --password <(echo $PW)
  ACCT=$(geth --datadir "data" --nousb account list | cut -d\  -f3 | sed -e 's/^{//' | cut -d\  -f3 | sed -e 's/}//' | tail -1)

  ../benchmarking gaa -ac $ACCT -am 1000000000000000000

  UNLOCKACCTS+=$ACCT","

  if [ $i -eq 1 ]
  then
    echo "!!!!!ADDING SIGNER!!!! -- "$ACCT
    ../benchmarking gas -ac $ACCT
  fi
done

geth --identity "geth-benchmark" init --datadir data genesis.json

# # echo "starting geth node"
geth --identity "geth-benchmark" --nousb --http --http.port "8545" --http.corsdomain "*" --datadir "data" --port "30303" --http.api "db,eth,net,web3,personal,miner,admin" --networkid 123 --nat "any" --allow-insecure-unlock --unlock ${UNLOCKACCTS::${#UNLOCKACCTS}-1} --password pw.txt --mine 2>> "geth.log" &

sleep 5

# # echo "starting geth attach"
geth --datadir "data" attach http://localhost:8545 2>> "gethcli.log" &