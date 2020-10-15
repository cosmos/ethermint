#!/bin/sh

set -e

# privkey=`yes 12345678 | injective-cli keys unsafe-export-eth-key genesis`
privkey="C0BBCB000CAB4A0475E27064AE82B57AD1137E2F2655A19BAAB2F29653EE2FCA"

registry_address=`evm-deploy-contract \
	--endpoint "http://localhost:1317" \
	--name Registry \
	--source Registry.sol \
	--privkey "$privkey" \
	--gas-price 0 \
	deploy`

echo "Deployed registry contract: $registry_address"

proxy_address=`evm-deploy-contract \
	--endpoint "http://localhost:1317" \
	--name UpgradeableProxy \
	--source util/UpgradeableProxy.sol \
	--privkey "$privkey" \
	--gas-price 0 \
	deploy 0xbeefE2577fFDecD66b073AAEAb627BA35Ef0378d "$registry_address"`

echo "Deployed proxy registry contract: $proxy_address"

evm-deploy-contract \
	--endpoint "http://localhost:1317" \
	--name Registry \
	--source Registry.sol \
	--privkey "$privkey" \
	--gas-price 0 \
	tx "$proxy_address" init 0xbeefE2577fFDecD66b073AAEAb627BA35Ef0378d

echo "Contract init done! Enjoy"
