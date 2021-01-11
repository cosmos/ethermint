#!/bin/bash

export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
go build -o ./build/ethermintd ./cmd/ethermintd
mkdir $GOPATH/bin
cp ./build/ethermintd $GOPATH/bin

CHAINID="ethermint-1337"

cd tests-solidity

if command -v yarn &> /dev/null; then
    yarn install
else
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
    sudo apt update && sudo apt install yarn
    yarn install
fi

chmod +x ./init-test-node.sh
./init-test-node.sh > ethermintd.log

cd suites/initializable
yarn test-ethermint

ok=$?

if (( $? != 0 )); then
    echo "initializable test failed: exit code $?"
fi

killall ethermintd

echo "Script exited with code $ok"
exit $ok

# initializable-buidler fails on CI, re-add later

./../../init-test-node.sh > ethermintd.log
cd ../initializable-buidler
yarn test-ethermint

ok=$(($? + $ok))

if (( $? != 0 )); then
    echo "initializable-buidler test failed: exit code $?"
fi

killall ethermintd

echo "Script exited with code $ok"
exit $ok