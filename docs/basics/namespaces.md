<!--
order: 5
-->

# Namespaces

A list of the rpc methods, the parameters and an example response.

## Json-RPC Methods
- [web3_clientVersion](#web3_clientVersion)

- [web3_sha3](#web3_sha3)

- [net_version](#net_version)

- [eth_protocolVersion](#eth_protocolVersion)

- [eth_syncing](#eth_syncing)

- [eth_gasPrice](#eth_gasPrice)

- [eth_accounts](#eth_accounts)

- [eth_blockNumber](#eth_blockNumber)

- [eth_getBalance](#eth_getBalance)

- [eth_getStorageAt](#eth_getStorageAt)

- [eth_getTransactionCount](#eth_getTransactionCount)

- [eth_getBlockTransactionCountByHash](#eth_getBlockTransactionCountByHash)

- [eth_getBlockTransactionCountByNumber](#eth_getBlockTransactionCountByNumber)

- [eth_getCode](#eth_getCode)

- [eth_sign](#eth_sign)

- [eth_sendTransaction](#eth_sendTransaction)

- [eth_sendRawTransaction](#eth_sendRawTransaction)

- [eth_call](#eth_call)

- [eth_estimateGas](#eth_estimateGas)

- [eth_getBlockByHash](#eth_getBlockByHash)

- [eth_getBlockByNumber](#eth_getBlockByNumber)

- [eth_getTransactionByHash](#eth_getTransactionByHash)

- [eth_getTransactionByBlockHashAndIndex](#eth_getTransactionByBlockHashAndIndex)

- [eth_getTransactionByBlockNumberAndIndex](#eth_getTransactionByBlockNumberAndIndex)

- [eth_getTransactionReceipt](#eth_getTransactionReceipt)

- [eth_newFilter](#eth_newFilter)

- [eth_newBlockFilter](#eth_newBlockFilter)

- [eth_newPendingTransactionFilter](#eth_newPendingTransactionFilter)

- [eth_uninstallFilter](#eth_uninstallFilter)

- [eth_getFilterChanges](#eth_getFilterChanges)

- [eth_getFilterLogs](#eth_getFilterLogs)

- [eth_getLogs](#eth_getLogs)

## Unused Methods
 - eth_mining
 - eth_coinbase
 - eth_hashrate
 - eth_getUncleCountByBlockHash
 - eth_getUncleCountByBlockNumber
 - eth_getUncleByBlockHashAndIndex
 - eth_getUncleByBlockNumberAndIndex

## Methods that are not implemented
 - net_peerCount 
 - net_listening	
 - eth_getTransactionbyBlockNumberAndIndex
 - eth_getWork
 - eth_submitWork
 - eth_submitHashrate
 - eth_getCompilers
 - eth_compileLLL
 - eth_compileSolidity
 - eth_compileSerpent
 - eth_signTransaction

## Web3 Methods

### web3_clientVersion

Get the web3 client version.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67}' -H "Content-Type: application/json" http://localhost:8545

// Result
 {"jsonrpc":"2.0","id":1,"result":"Ethermint/0.0.0+/linux/go1.14"}
```

### web3_sha3

Returns Keccak-256 (not the standardized SHA3-256) of the given data.

#### Parameters

- the data to convert into a SHA3 hash

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"web3_sha3","params":["0x67656c6c6f20776f726c64"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x1b84adea42d5b7d192fd8a61a85b25abe0757e9a65cab1da470258914053823f"}
```

## Net Methods

### net_version

Returns the current network id.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"8"}
```

## Eth Methods

### eth_protocolVersion

Returns the current ethereum protocol version.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x3f"}
```

### eth_syncing

The sync status object may need to be different depending on the details of Tendermint's sync protocol. However, the 'synced' result is simply a boolean, and can easily be derived from Tendermint's internal sync state.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":false}
```

### eth_gasPrice

Returns the current gas price in aphotons.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x0"}
```

### eth_accounts

Returns array of all eth accounts.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":["0x3b7252d007059ffc82d16d022da3cbf9992d2f70","0xddd64b4712f7c8f1ace3c145c950339eddaf221d","0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0"]}
```

### eth_blockNumber

Returns the current block height.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x66"}
```

### eth_getBalance

Returns the account balance for a given account address and Block Number.

#### Parameters

- Accout Address

- Block Number


```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x7bf7b17da59880d9bcca24915679668db75f9397", "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x36354d5575577c8000"}
```

### eth_getStorageAt

Returns the storage address for a given account address. // i need to learn how to find the address store key so i can get a real response.

#### Parameters

- Accout Address

- Address Store Key

- Block Number

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getStorageAt","params":["0x7bf7b17da59880d9bcca24915679668db75f9397", "0"  "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000000"}
```

### eth_getTransactionCount

Returns the total transaction for a given account address and Block Number.

#### Parameters

- Accout Address

- Block Number

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x7bf7b17da59880d9bcca24915679668db75f9397", "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x8"}
```

### eth_getBlockTransactionCountByNumber

Returns the total transaction count for a given block number.

#### Parameters

- Block number

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByNumber","params":["0x1"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
 {"jsonrpc":"2.0","id":1,"result":{"difficulty":null,"extraData":"0x0","gasLimit":"0xffffffff","gasUsed":"0x0","hash":"0x8101cc04aea3341a6d4b3ced715e3f38de1e72867d6c0db5f5247d1a42fbb085","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x0000000000000000000000000000000000000000","nonce":null,"number":"0x17d","parentHash":"0x70445488069d2584fea7d18c829e179322e2b2185b25430850deced481ca2e77","sha3Uncles":null,"size":"0x1df","stateRoot":"0x269bb17fe7adb8dd5f15f57b717979f82078d6b7a675c1ba1b0da2d27e415fcc","timestamp":"0x5f5ba97c","totalDifficulty":null,"transactions":[],"transactionsRoot":"0x","uncles":[]}}
```

### eth_getBlockTransactionCountByHash

Returns the total transaction count for a given block hash.

#### Parameters

- Block Hash

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByHash","params":["0x8101cc04aea3341a6d4b3ced715e3f38de1e72867d6c0db5f5247d1a42fbb085"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x3"}
```

### eth_getCode

Returns the code for a given account address and Block Number.

#### Parameters

- Accout Address

- Block Number

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getCode","params":["0x7bf7b17da59880d9bcca24915679668db75f9397", "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0xef616c92f3cfc9e92dc270d6acff9cea213cecc7020a76ee4395af09bdceb4837a1ebdb5735e11e7d3adb6104e0c3ac55180b4ddf5e54d022cc5e8837f6a4f971b"}
```

### eth_sign

The sign method calculates an Ethereum specific signature with: sign(keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))).

By adding a prefix to the message makes the calculated signature recognisable as an Ethereum specific signature. This prevents misuse where a malicious DApp can sign arbitrary data (e.g. transaction) and use the signature to impersonate the victim.

::: warning 
the address to sign with must be unlocked.
:::

#### Parameters

- Account Address 

- Message to sign

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sign","params":["0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "0xdeadbeaf"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x909809c76ed2a5d38733de39207d0f411222b9b49c64a192bf649cb13f63f37b45acb4f6939facb4f1c277bc70fb00407564140c0f18600ac44388f2c1dfd1dc1b"}
```


### eth_sendTransaction

Sends transaction from given account to a given account.

#### Parameters

 - Object containing: 

    from: DATA, 20 Bytes - The address the transaction is send from.

    to: DATA, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.

    gas: QUANTITY - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.

    gasPrice: QUANTITY - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas

    value: QUANTITY - value sent with this transaction

    data: DATA - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see Ethereum Contract ABI

    nonce: QUANTITY - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.


```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0", "to":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "value":"0x16345785d8a0000", "gasLimit":"0x5208", "gasPrice":"0x55ae82600"}],"id":1}'  -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x33653249db68ebe5c7ae36d93c9b2abc10745c80a72f591e296f598e2d4709f6"}
```

### eth_sendRawTransaction

Creates new message call transaction or a contract creation for signed transactions.

#### Parameters

-  The signed transaction data

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[I dont know how to get that],"id":1}'  -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000000"}
```

### eth_call

Executes a new message call immediately without creating a transaction on the block chain.

#### Parameters

- Object containing:

    from: DATA, 20 Bytes - (optional) The address the transaction is sent from.

    to: DATA, 20 Bytes - The address the transaction is directed to.

    gas: QUANTITY - gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.

    gasPrice: QUANTITY - gasPrice used for each paid gas

    value: QUANTITY - value sent with this transaction

    data: DATA - (optional) Hash of the method signature and encoded parameters. For details see Ethereum Contract ABI in the Solidity documentation

- Block number



```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_call","params":[{"from":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "to":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d", "gas":"0x5208", "gasPrice":"0x55ae82600", "value":"0x16345785d8a0000", "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"}, "0x0"],"id":1}'  -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x"}
```

### eth_estimateGas

Returns an estimate value of the gas required to send the transaction.

#### Parameters

- Object containing: 

    from: DATA, 20 Bytes - The address the transaction is send from.

    to: DATA, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.

    value: QUANTITY - value sent with this transaction

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{"from":"0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0", "to":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "value":"0x16345785d8a00000"}],"id":1}'  -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x1199b"}
```


### eth_getBlockByNumber

Returns information about a block by block number.

#### Parameters

- QUANTITY|TAG - block number, or the string "earliest", "latest" or "pending", as in the default block parameter.

- Boolean - If true it returns the full transaction objects, if false only the hashes of the transactions.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1", false],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":{"difficulty":null,"extraData":"0x0","gasLimit":"0xffffffff","gasUsed":null,"hash":"0xabac6416f737a0eb54f47495b60246d405d138a6a64946458cf6cbeae0d48465","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x0000000000000000000000000000000000000000","nonce":null,"number":"0x1","parentHash":"0x","sha3Uncles":null,"size":"0x9b","stateRoot":"0x","timestamp":"0x5f5bd3e5","totalDifficulty":null,"transactions":[],"transactionsRoot":"0x","uncles":[]}}
```

### eth_getBlockByHash

Returns the block info given the hash found in the command above and a bool.

#### Parameters 
- DATA, 32 Bytes - Hash of a block.

- Boolean - If true it returns the full transaction objects, if false only the hashes of the transactions.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4", false],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":{"difficulty":null,"extraData":"0x0","gasLimit":"0xffffffff","gasUsed":null,"hash":"0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4","logsBloom":"0x00000000100000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000002000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x0000000000000000000000000000000000000000","nonce":null,"number":"0xc","parentHash":"0x404e58f31a9ede1b614b98701d6b0fbf1450f186842dbcf6426dd16811a5ca0d","sha3Uncles":null,"size":"0x307","stateRoot":"0x599ccdb111fc62c6398dc39be957df8e97bf8ab72ce6c06ff10641a92b754627","timestamp":"0x5f5fdbbd","totalDifficulty":null,"transactions":["0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615"],"transactionsRoot":"0x4764dba431128836fa919b83d314ba9cc000e75f38e1c31a60484409acea777b","uncles":[]}}
```

### eth_getTransactionByHash

Returns transaction details given the ethereum tx something.

#### Parameters

- DATA, 32 Bytes - hash of a transaction

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"],"id":1}'

// Result
// no idea how to get this to work
```

### eth_getTransactionByBlockHashAndIndex

Returns transaction details given the block hash and the transaction index. 

#### Parameters

- DATA, 32 Bytes - hash of a block.

- QUANTITY - integer of the transaction index position.

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockHashAndIndex","params":["0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4", "0x0"],"id":1}'

// Result
{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4","blockNumber":"0xc","from":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d","gas":"0x4c4b40","gasPrice":"0x3b9aca00","hash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","input":"0x4f2be91f","nonce":"0x0","to":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","transactionIndex":"0x0","value":"0x0","v":"0xa96","r":"0xced57d973e58b0f634f776d57daf41d3d3387ceb347a3a72ca0746e5ec2b709e","s":"0x384e89e209a5eb147a2bac3a4e399507400ac7b29cd155531f9d6203a89db3f2"}}
```

### eth_getTransactionReceipt

Returns the receipt of a transaction by transaction hash.
 

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea614"],"id":1}'

// Result
{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4","blockNumber":"0xc","contractAddress":"0x0000000000000000000000000000000000000000","cumulativeGasUsed":null,"from":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d","gasUsed":"0x5289","logs":[{"address":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","topics":["0x64a55044d1f2eddebe1b90e8e2853e8e96931cefadbfa0b2ceb34bee36061941"],"data":"0x0000000000000000000000000000000000000000000000000000000000000002","blockNumber":"0xc","transactionHash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","topics":["0x938d2ee5be9cfb0f7270ee2eff90507e94b37625d9d2b3a61c97d30a4560b829"],"data":"0x0000000000000000000000000000000000000000000000000000000000000002","blockNumber":"0xc","transactionHash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x1","removed":false}],"logsBloom":"0x00000000100000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000002000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","transactionHash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","transactionIndex":"0x0"}}
```

### eth_newFilter

Create new filter using topics of some kind.
 
#### Parameters

- DATA, 32 Bytes - hash of a transaction

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_newFilter","params":[{"topics":["0x0000000000000000000000000000000000000000000000000000000012341234"]}],"id":1}'

// Result
{"jsonrpc":"2.0","id":1,"result":"0xdc714a4a2e3c39dc0b0b84d66a3ccb00"}
```

### eth_newBlockFilter

Creates a filter in the node, to notify when a new block arrives.
 

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_newBlockFilter","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x3503de5f0c766c68f78a03a3b05036a5"}
```

### eth_newPendingTransactionFilter

Creates a filter in the node, to notify when new pending transactions arrive.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_newPendingTransactionFilter","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"result":"0x9daacfb5893d946997d3801ea18e9902"}
```

### eth_uninstallFilter

Removes the filter with the given filter id. Returns true if the filter was successfully uninstalled, otherwise false.

#### Parameters

- The filter id

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_uninstallFilter","params":["0xb91b6608b61bf56288a661a1bd5eb34a"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"result":true}
```

### eth_getFilterChanges

Polling method for a filter, which returns an array of logs which occurred since last poll.

#### Parameters

- The filter id

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getFilterChanges","params":["0xb91b6608b61bf56288a661a1bd5eb34a"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"result":[]} // I need to find a result that isnt empty
```

### eth_getFilterLogs

Polling method for a filter, which returns an array of logs which occurred since last poll.

#### Parameters

- The filter id

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getFilterLogs","params":["0xb91b6608b61bf56288a661a1bd5eb34a"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"filter 0x35b64c227ce30e84fc5c7bd347be380e doesn't have a LogsSubscription type: got 5"}} // I couldnt get one that didnt error
```

### eth_getLogs

Returns an array of all logs matching a given filter object.

#### Parameters

- some string array of topics

- a block to check the logs from

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":["0x775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd738898","0x0000000000000000000000000000000000000000000000000000000000000011"], "fromBlock":"earliest"}],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"result":[]} // I need to find a result that isnt empty
```

### eth_getProof

Returns an array of all logs matching a given filter object. // I need to learn to get a store key so that i can get a real response

#### Parameters

- Accout Address

- Address Store Key

- Block Number

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getProof","params":["0x3b7252d007059ffc82d16d022da3cbf9992d2f70", ["0"],  "latest"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"result":{"address":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70","accountProof":["ops:\u003ctype:\"iavl:v\" key:\"\\001;rR\\320\\007\\005\\237\\374\\202\\321m\\002-\\243\\313\\371\\231-/p\" data:\"\\224\\002\\n\\221\\002\\n)\\010\\n\\020\\016\\030\\331\\010* |{\\234\\221\\314\\354k\\335!\\256$z\\216\\004g\\360\\315\\026\\000F\\217\\305\\200\\244F\\t\\306\\320\\250GzH\\n)\\010\\010\\020\\t\\030\\331\\010* \\034\\365\\230\\241\\336\\341*V\\200K\\217\\356\\306$\\304\\300z\\371\\274\\007\\316\\320\\340\\303\\235\\\\\\266:lK\\351\\036\\n(\\010\\006\\020\\006\\030\\021* f\\202\\222\\247*\\021\\\"a\\n\\034\\311X\\370\\221\\210\\201-q\\n\\332.\\n\\024\\311C\u003e\\345\u003e\\342\\315\\330\\004\\n(\\010\\004\\020\\003\\030\\021\\\" \\260\\202\\254\\240d\\323?\\024\\32157\\244\\212V\\240\\211\\220\\354p\\326\\355\\220\\036%\\004\\177\\346\\273\\177\\037\\332\\312\\n(\\010\\002\\020\\002\\030\\021* 5d\\1776\\005\\023\\010n}b\\347\\202]A\\014\\270\\302!e\\025\\344$AP\u003e;\\\"t\\034\\322T9\\032;\\n\\025\\001;rR\\320\\007\\005\\237\\374\\202\\321m\\002-\\243\\313\\371\\231-/p\\022 \\031\\336\\005PS\\342\\016\\247\\\\8V\\210o\\313qJ\\220u\\206\\021g\\256\\367P\\362\\345\\363\\220\\374\\201\\307\\211\\030\\021\" \u003e ops:\u003ctype:\"multistore\" key:\"acc\" data:\"\\325\\004\\n\\322\\004\\n1\\n\\006supply\\022'\\n%\\010\\331\\010\\022 \\354V.\\325\\245\\002\\002\\345\\030\\267\\365\\343\\234\\306Q\\246\\210i\\024\\264\\310U\\202\\331\\222\\3719S\\227\\254\\037\\r\\n.\\n\\003gov\\022'\\n%\\010\\331\\010\\022 `\\315\\023|\\031b\\354\\254ac\\211\\326\\2004\\203?)!P\\234-(Z\\245\\365\\0259\\227\\316\\226\\212t\\n\\021\\n\\010evidence\\022\\005\\n\\003\\010\\331\\010\\n7\\n\\014distribution\\022'\\n%\\010\\331\\010\\022 y[\\202l\\000\\007I\\256\\024\\016\\225y\\023\\206=U$\\006\\204\\022e\\000qG\\326\\263\\223\\274n/i0\\n.\\n\\003acc\\022'\\n%\\010\\331\\010\\022 \\037\\267\\232\\201\\030\\262\\306\\307\\272s\\300\\211\\t\\263\\303\\025O\\362\\350|\\277\\323H\\354c\\310w\\276\\225V\\362\\n\\n\\020\\n\\007upgrade\\022\\005\\n\\003\\010\\331\\010\\n/\\n\\004main\\022'\\n%\\010\\331\\010\\022 \\261\\234\\360\\230\\264;\u003c}Ty\\237\\357\\tD\\252\\225\\t\\237\\202\\335\\302\\031\\372\\014\\200\\222\\344R\\242\\202\\277A\\n2\\n\\007staking\\022'\\n%\\010\\331\\010\\022 \\317\\263\\342\\230\\010D\\254$\\226\\250\\004\\274\\t_\\rpf\\235e1\\224\\265V\\204\\331}f\\276C\\033O\\217\\n1\\n\\006faucet\\022'\\n%\\010\\331\\010\\022 \\024u\\250\\316\\2506\\351\\376?\\236\\315\\376y\\226eeg\\313]\\230\\034\\376\\215n\\341!\\311\\274|\\251\\214\\006\\n/\\n\\004mint\\022'\\n%\\010\\331\\010\\022 f\\005\\307D\\206.K\\220\\0056yt\\000\\026\\n\\2476\\217\\254\\274u\\271\\\\\\221\\270O\\375\\034\\240\\024\\000y\\n1\\n\\006params\\022'\\n%\\010\\331\\010\\022 \\274I'\\344\\227\\213\\032\\347\\310\\257\\342H\\205\\225\\343\\346\\311\\337\\r\\340\\203\\327\\377\\317._\\301\\244\u0026{@\\335\\n.\\n\\003evm\\022'\\n%\\010\\331\\010\\022 \\220\\213\\010\\375\\365\\364\\037=\\266\\223\\273+\\031\\020\\237\\236Q\\220@R\\350\\262\\2437\\346\\033u\\021@8\\3207\\n3\\n\\010slashing\\022'\\n%\\010\\331\\010\\022 h\\335-\\021\\223\\337\\000X\\036\\364h\\202\\264~\\\\\\004='\\265\\206\\307\\262\\233\\204\\317\\246\\342b\\301\\021%e\" \u003e "],"balance":"0x36354d5575577c8000","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470","nonce":"0x8","storageHash":"0x0000000000000000000000000000000000000000000000000000000000000000","storageProof":[{"key":"0","value":"0x0","proof":[""]}]}}
```

### eth_subscribe

Returns an array of all logs matching a given filter object.

#### Parameters

- something

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_subscribe","params":["newHeads", {}],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"notifications not supported"}} // I dont know how to get this to work
```

### eth_unsubscribe

Returns an array of all logs matching a given filter object.

#### Parameters

- something

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_unsubscribe","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Request
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"notifications not supported"}} // I dont know how to get this to work
```

