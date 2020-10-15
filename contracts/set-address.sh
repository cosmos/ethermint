#!/bin/sh

set -e

# privkey=`yes 12345678 | injective-cli keys unsafe-export-eth-key genesis`
privkey="C0BBCB000CAB4A0475E27064AE82B57AD1137E2F2655A19BAAB2F29653EE2FCA"

if [ -z $1 ]
then
	echo "No contract name is provided"
	exit 1
fi

if [ -z $2 ]
then
	echo "No contract address is provided"
	exit 1
fi

change_request=`evm-deploy-contract \
	--endpoint "http://localhost:1317" \
	--name Registry \
	--source Registry.sol \
	--privkey "$privkey" \
	--gas-price 0 \
	tx 0x5C7e1fc74fe17242a077DB7DFd962e897Ed4e39a requestContractAddressChange "$1" "$2"`

logs=`evm-deploy-contract \
	--endpoint "http://localhost:1317" \
	--name Registry \
	--source Registry.sol \
	logs 0x5C7e1fc74fe17242a077DB7DFd962e897Ed4e39a "$change_request" ContractAddressLocked`

lock_id=`echo "$logs" | jq --raw-output ".[0].lockId | implode" | iconv -c -f utf-8 -t latin1 | xxd -ps -c 32 | head -n 1`

if [ -z $lock_id ]
then
	echo "No lock ID found in the event logs"
	exit 1
fi

echo "Lock ID: $lock_id"

evm-deploy-contract \
	--endpoint "http://localhost:1317" \
	--name Registry \
	--source Registry.sol \
	--privkey "$privkey" \
	--gas-price 0 \
	tx 0x5C7e1fc74fe17242a077DB7DFd962e897Ed4e39a confirmContractAddressChange "$lock_id"

echo "Contract $1 changed address to $2"
