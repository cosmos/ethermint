<!--
order: 5
-->


# Joining Chainsafe's Public Testnet

This document outlines the steps to join the public testnet hosted by [Chainsafe](https://chainsafe.io). 

## Steps:
1. Install the Ethermint binaries (ethermintd & ethermint cli):
```
git clone https://github.com/cosmos/ethermint
cd ethermint
make install
```

2. Create an Ethermint account:
```
ethermintcli keys add <keyname>
```

3. Copy genesis file:
Follow this [link](https://gist.github.com/araskachoi/43f86f3edff23729b817e8b0bb86295a) and copy it over to the directory ~/.ethermintd/config/genesis.json

4. Add peers:
Edit the file located in ~/.ethermintd/config/config.toml and edit line 350 (persistent_peers) to the following;
```
"f618ed93d19dd0b162ef8ba8707f924674454a58@54.210.246.165:26656,24891ad476acc000b20507128ed18e20c6faacc7@18.204.206.179:26656,b765eda477f864743887219d1536698005d7c13e@3.86.104.251:26656"
```

5. Validate genesis and start the Ethermint network:
```
ethermintd validate-genesis
```
```
ethermintd start --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info" --trace
```
(we recommend running the command in the background for convenience)

6. Start the RPC server:
```
ethermintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $KEY --chain-id etherminttestnet-1 --trace
```
where `$KEY` is the key name that was used in step 2.
(we recommend running the command in the background for convenience)

7. Request funds from the faucet:
You will need to know the Ethereum hex address, and it can be found with the following command:

```
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545
```
Using the output of the above command, you will then send the command with your valid Ethereum address:
```
curl --header "Content-Type: application/json" --request POST --data '{"address":"0xYouEthereumHexAddress"}' 3.95.21.91:3000
```
