<!--
order: 5
-->

# JSON-RPC Server

Check the JSON-RPC methods and namespaces supported on Ethermint. {synopsis}

## Pre-requisite Readings

- [Ethereum JSON-RPC](https://eth.wiki/json-rpc/API) {prereq}
- [Geth JSON-RPC APIs](https://geth.ethereum.org/docs/rpc/server) {prereq}

## JSON-RPC Methods

| Method                                                                            | Namespace | Implemented |
|-----------------------------------------------------------------------------------|-----------|-------------|
| [`web3_clientVersion`](#web3_clientVersion)                                       | Web3      | ✔           |
| [`web3_sha3`](#web3_sha3)                                                         | Web3      | ✔           |
| [`net_version`](#net_version)                                                     | Net       | ✔           |
| `net_peerCount`                                                                   | Net       |             |
| `net_listening`                                                                   | Net       |             |
| [`eth_protocolVersion`](#eth_protocolVersion)                                     | Eth       | ✔           |
| [`eth_syncing`](#eth_syncing)                                                     | Eth       | ✔           |
| [`eth_gasPrice`](#eth_gasPrice)                                                   | Eth       | ✔           |
| [`eth_accounts`](#eth_accounts)                                                   | Eth       | ✔           |
| [`eth_blockNumber`](#eth_blockNumber)                                             | Eth       | ✔           |
| [`eth_getBalance`](#eth_getBalance)                                               | Eth       | ✔           |
| [`eth_getStorageAt`](#eth_getStorageAt)                                           | Eth       | ✔           |
| [`eth_getTransactionCount`](#eth_getTransactionCount)                             | Eth       | ✔           |
| [`eth_getBlockTransactionCountByNumber`](#eth_getBlokTransactionCountByNumber)    | Eth       | ✔           |
| [`eth_getBlockTransactionCountByHash`](#eth_getBlockTransactionCountByHash)       | Eth       | ✔           |
| [`eth_getCode`](#eth_getCode)                                                     | Eth       | ✔           |
| [`eth_sign`](#eth_sign)                                                           | Eth       | ✔           |
| [`eth_sendTransaction`](#eth_sendTransaction)                                     | Eth       | ✔           |
| [`eth_sendRawTransaction`](#eth_sendRawTransaction)                               | Eth       | ✔           |
| [`eth_call`](#eth_call)                                                           | Eth       | ✔           |
| [`eth_estimateGas`](#eth_estimateGas)                                             | Eth       | ✔           |
| [`eth_getBlockByNumber`](#eth_getBlockByNumber)                                   | Eth       | ✔           |
| [`eth_getBlockByHash`](#eth_getBlockByHash)                                       | Eth       | ✔           |
| [`eth_getTransactionByHash`](#eth_getTransactionByHash)                           | Eth       | ✔           |
| [`eth_getTransactionByBlockHashAndIndex`](#eth_getTransactionByBlockHashAndIndex) | Eth       | ✔           |
| [`eth_getTransactionReceipt`](#eth_getTransactionReceipt)                         | Eth       | ✔           |
| [`eth_newFilter`](#eth_newFilter)                                                 | Eth       | ✔           |
| [`eth_newBlockFilter`](#eth_newBlockFilter)                                       | Eth       | ✔           |
| [`eth_newPendingTransactionFilter`](#eth_newPendingTransactionFilter)             | Eth       | ✔           |
| [`eth_uninstallFilter`](#eth_uninstallFilter)                                     | Eth       | ✔           |
| [`eth_getFilterChanges`](#eth_getFilterChanges)                                   | Eth       | ✔           |
| [`eth_getLogs`](#eth_getLogs)                                                     | Eth       | ✔           |
| [`eth_subscribe`](#eth_subscribe)                                                 | Websocket | ✔           |
| [`eth_unsubscribe`](#eth_unsubscribe)                                             | Websocket | ✔           |
| `eth_getTransactionbyBlockNumberAndIndex`                                         | Eth       |             |
| `eth_getWork`                                                                     | Eth       |             |
| `eth_submitWork`                                                                  | Eth       |             |
| `eth_submitHashrate`                                                              | Eth       |             |
| `eth_getCompilers`                                                                | Eth       |             |
| `eth_compileLLL`                                                                  | Eth       |             |
| `eth_compileSolidity`                                                             | Eth       |             |
| `eth_compileSerpent`                                                              | Eth       |             |
| `eth_signTransaction`                                                             | Eth       |             |

:::tip
Block Number can be entered as a Hex string, `"earliest"`, `"latest"` or `"pending"`.
:::

Below is a list of the RPC methods, the parameters and an example response from the namespaces.

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

- The data to convert into a SHA3 hash

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
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0", "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x36354d5575577c8000"}
```

### eth_getStorageAt

Returns the storage address for a given account address. // i need to learn how to find the address storage key so i can get a real response.

#### Parameters

- Accout Address

- Address Storage Key

- Block Number

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getStorageAt","params":["0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0", "0xb47cde69de5130ac4310768396858d7fc20ee04b75e353ac8d5a991f3fbf5691"  "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

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
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "to":"0x0f54f47bf9b8e317b214ccd6a7c3e38b893cd7f0", "value":"0x16345785d8a0000", "gasLimit":"0x5208", "gasPrice":"0x55ae82600"}],"id":1}'  -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":"0x33653249db68ebe5c7ae36d93c9b2abc10745c80a72f591e296f598e2d4709f6"}
```

### eth_sendRawTransaction

Creates new message call transaction or a contract creation for signed transactions.

You can get signed transaction data using the personal_sign method

#### Parameters

-  The signed transaction data

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf9ff74c86aefeb5f6019d77280bbb44fb695b4d45cfe97e6eed7acd62905f4a85034d5c68ed25a2e7a8eeb9baf1b8401e4f865d92ec48c1763bf649e354d900b1c"],"id":1}'  -H "Content-Type: application/json" http://localhost:8545

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

- Block Number

- If true it returns the full transaction objects, if false only the hashes of the transactions.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1", false],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":{"difficulty":null,"extraData":"0x0","gasLimit":"0xffffffff","gasUsed":null,"hash":"0xabac6416f737a0eb54f47495b60246d405d138a6a64946458cf6cbeae0d48465","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x0000000000000000000000000000000000000000","nonce":null,"number":"0x1","parentHash":"0x","sha3Uncles":null,"size":"0x9b","stateRoot":"0x","timestamp":"0x5f5bd3e5","totalDifficulty":null,"transactions":[],"transactionsRoot":"0x","uncles":[]}}
```

### eth_getBlockByHash

Returns the block info given the hash found in the command above and a bool.

#### Parameters 
- Hash of a block.

- If true it returns the full transaction objects, if false only the hashes of the transactions.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4", false],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":{"difficulty":null,"extraData":"0x0","gasLimit":"0xffffffff","gasUsed":null,"hash":"0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4","logsBloom":"0x00000000100000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000002000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x0000000000000000000000000000000000000000","nonce":null,"number":"0xc","parentHash":"0x404e58f31a9ede1b614b98701d6b0fbf1450f186842dbcf6426dd16811a5ca0d","sha3Uncles":null,"size":"0x307","stateRoot":"0x599ccdb111fc62c6398dc39be957df8e97bf8ab72ce6c06ff10641a92b754627","timestamp":"0x5f5fdbbd","totalDifficulty":null,"transactions":["0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615"],"transactionsRoot":"0x4764dba431128836fa919b83d314ba9cc000e75f38e1c31a60484409acea777b","uncles":[]}}
```

### eth_getTransactionByHash

Returns transaction details given the ethereum tx something.

#### Parameters

- hash of a transaction

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xec5fa15e1368d6ac314f9f64118c5794f076f63c02e66f97ea5fe1de761a8973"],"id":1}' -H "Content-Type: application/json" http://localhost:8545
 
// Result
{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x7a7398cc11d9c4c8e6f53e0c73824297aceafdab62db9e4b867a0da694384864","blockNumber":"0x188","from":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70","gas":"0x147ee","gasPrice":"0x3b9aca00","hash":"0xec5fa15e1368d6ac314f9f64118c5794f076f63c02e66f97ea5fe1de761a8973","input":"0x6dba746c","nonce":"0x18","to":"0xa655256f589060437e5ffe2246dec385d040f148","transactionIndex":"0x0","value":"0x0","v":"0xa96","r":"0x6db399d694a452fb4106419140a6e5dbbe6817743a0f6f695a651e6576e59a5e","s":"0x25dd6ab1f936d0280d2fed0caeb0ebe5b9a46de6d8cb08ad8fd2c88deb55fc31"}}
```

### eth_getTransactionByBlockHashAndIndex

Returns transaction details given the block hash and the transaction index. 

#### Parameters

- Hash of a block.

- Transaction index position.

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockHashAndIndex","params":["0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4", "0x0"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4","blockNumber":"0xc","from":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d","gas":"0x4c4b40","gasPrice":"0x3b9aca00","hash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","input":"0x4f2be91f","nonce":"0x0","to":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","transactionIndex":"0x0","value":"0x0","v":"0xa96","r":"0xced57d973e58b0f634f776d57daf41d3d3387ceb347a3a72ca0746e5ec2b709e","s":"0x384e89e209a5eb147a2bac3a4e399507400ac7b29cd155531f9d6203a89db3f2"}}
```

### eth_getTransactionReceipt

Returns the receipt of a transaction by transaction hash.

#### Parameters

- Hash of a transaction

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea614"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x1b9911f57c13e5160d567ea6cf5b545413f96b95e43ec6e02787043351fb2cc4","blockNumber":"0xc","contractAddress":"0x0000000000000000000000000000000000000000","cumulativeGasUsed":null,"from":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d","gasUsed":"0x5289","logs":[{"address":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","topics":["0x64a55044d1f2eddebe1b90e8e2853e8e96931cefadbfa0b2ceb34bee36061941"],"data":"0x0000000000000000000000000000000000000000000000000000000000000002","blockNumber":"0xc","transactionHash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","topics":["0x938d2ee5be9cfb0f7270ee2eff90507e94b37625d9d2b3a61c97d30a4560b829"],"data":"0x0000000000000000000000000000000000000000000000000000000000000002","blockNumber":"0xc","transactionHash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x1","removed":false}],"logsBloom":"0x00000000100000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000002000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0x439c697e0742a0ddb124a376efd62a72a94ac35a","transactionHash":"0xae64961cb206a9773a6e5efeb337773a6fd0a2085ce480a174135a029afea615","transactionIndex":"0x0"}}
```

### eth_newFilter

Create new filter using topics of some kind.

#### Parameters

- hash of a transaction

```json
// Request
curl localhost:8545 -H "Content-Type:application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_newFilter","params":[{"topics":["0x0000000000000000000000000000000000000000000000000000000012341234"]}],"id":1}' -H "Content-Type: application/json" http://localhost:8545

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

// Result
{"jsonrpc":"2.0","id":1,"result":"0x9daacfb5893d946997d3801ea18e9902"}
```

### eth_uninstallFilter

Removes the filter with the given filter id. Returns true if the filter was successfully uninstalled, otherwise false.

#### Parameters

- The filter id

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_uninstallFilter","params":["0xb91b6608b61bf56288a661a1bd5eb34a"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":true}
```

### eth_getFilterChanges

Polling method for a filter, which returns an array of logs which occurred since last poll.

#### Parameters

- The filter id

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getFilterChanges","params":["0x127e9eca4f7751fb4e5cb5291ad8b455"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":["0xc6f08d183a81e149896fc5317c872f9092068e88e956ca1864e9bd4c81c09b44","0x3ca6dfb5be15549d721d1b3d10c1bec50ed6217c9ac7b61df361fac9692a27e5","0x776fffac134171acb1ebf2e59856625501ad5ccc5c4c8fe0359e0d4dff8919f2","0x84123103704dbd738c089276ab2b04b5936330b24f6e78453c4ba8bf4848aaf9","0xffddbe5bd8e8aa41e44002daa9ea89ade9e6980a0d83f51d104cf16498827eca","0x53430e49963e8ae32605d8f22dec2e757a691e6436d593854ca4d9383eeab86a","0x975948058c9351a91fbec332ca00dda39d1a919f5f16b996a4c7e30c38ba423b","0x619e37e32024c8efef7f7220e6caff4ee1d682ea78b2ac91e0a6b30850dc0677","0x31a5d985a40d08303ac68000ce008df512bcd1a911c497415c97f0624b4a271a","0x91dcf1fce4503a8dbb3e6fb61073f25cd31d69c766ecba639fefde4436e59d07","0x606d9e0143cfdb410a6812c590a8135b5c6b5c59eec26d760d5cd930aa47257d","0xd3c00b859b29b20ba654415eef648ef58251389c73a138580db87675b0d5465f","0x954391f0eb50888be90489898016ebb54f750f612f3adec2a00854955d5e52d8","0x698905f06aff921a9e9fcef39b8b0d107747c3e6204d2ea79cf4c12debf8d253","0x9fcafec5721938a06eb8e2951ede4b6ef8fae54a8c8f85f3166ec9782a0032b5","0xaec6d3364e47a5716ba69e4705f3c705d017f81298859589591183bfea87be7a","0x91bf2ee13319b6eaca96ed89c126437b66c4df1b13560c6a9bb18556ee3b7e1f","0x4f426dc1fc0ea8149052033065b237892d2d34927b2d558ab50c5a7fb98d6e79","0xdd809fb07e5aab638fef5311371b4e2b27c9c9a6183fde0cdd2b7724f6d2a89b","0x7e12fc92ab953e233a304959a2a8474d96195e71efd9388fdceb1326a577811a","0x30618ef6b490c3cc9979c47163459db37c1a1e0aa5793c56accd417f9d89973b","0x614609f06ee24bae7408e45895b1a25e6b19a8159aeea7a95c9d1339d9ba286f","0x115ddc6d533620040791d241f01f1c5ae3d9d1a8f64b15af5e9793e4d9096e22","0xb7458c9323beeca2cd54f32a6af5671f3cd5a7a251aed9d82bdd6ebe5f56305b","0x573dd48a5ba7bf4cc3d49597cd7419f75ecc9897258f1ebadebd670446d0d358","0xcb6670918439f9698413b53f3b5336d82ca4be152fdefaacf45e052fff6262fc","0xf3fe2a8945abafd269ab97bfdc80b3dbff2202ffdce59a227f952874b966b230","0x989980707007533cc0840a079f77f261a2e818abae1a1ffd3af02f3fff1d35fd","0x886b6ae365fec996be8a9a2c31cf4cda97ff8352908be2c83f17abd66ef1591e","0xfd90df68706ef95a62b317de93d6899a9bd6c80416e42d007f5c30fcdedfce24","0x7af8491fbb0373886d9032bb74e0ef52ed9e100f260b79bd15f46126b38cbede","0x91d1e2cd55533cf7dd5de86c9aa73295e811b1279be193d429bbd6ba83810e16","0x6b65b3128c2104005a04923288fe2aa33a2477a4962bef70532f94cab582f2a7"]}
```

<!-- 
### eth_getFilterLogs

Polling method for a filter, which returns an array of logs which occurred since last poll.

#### Parameters

- The filter id

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getFilterLogs","params":["0x127e9eca4f7751fb4e5cb5291ad8b455"],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"filter 0x35b64c227ce30e84fc5c7bd347be380e doesn't have a LogsSubscription type: got 5"}} 
``` -->

### eth_getLogs

Returns an array of all logs matching a given filter object.

#### Parameters

- Object containing:

    fromBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.

    toBlock: QUANTITY|TAG - (optional, default: "latest") Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.

    address: DATA|Array, 20 Bytes - (optional) Contract address or a list of addresses from which logs should originate.

    topics: Array of DATA, - (optional) Array of 32 Bytes DATA topics. Topics are order-dependent. Each topic can also be an array of DATA with “or” options.
    
    blockhash: (optional, future) With the addition of EIP-234, blockHash will be a new filter option which restricts the logs returned to the single block with the 32-byte hash blockHash. Using blockHash is equivalent to fromBlock = toBlock = the block number with hash blockHash. If blockHash is present in in the filter criteria, then neither fromBlock nor toBlock are allowed.

```json
// Request
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":["0x775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd738898","0x0000000000000000000000000000000000000000000000000000000000000011"], "fromBlock":"latest"}],"id":1}' -H "Content-Type: application/json" http://localhost:8545

// Result
{"jsonrpc":"2.0","id":1,"result":[]}
```

## WebSocket Methods

Read about websockets in [events](./../quickstart/events.md) {hide}

### eth_subscribe

subscribe using JSON-RPC notifications. This allows clients to wait for events instead of polling for them.

It works by subscribing to particular events. The node will return a subscription id. For each event that matches the subscription a notification with relevant data is send together with the subscription id.

#### Parameters

- Subscription Name

- Optional Arguments

```json
// Request
{"id": 1, "method": "eth_subscribe", "params": ["newHeads", {"includeTransactions": true}]}

// Result
< {"jsonrpc":"2.0","result":"0x34da6f29e3e953af4d0c7c58658fd525","id":1}
```

### eth_unsubscribe

Unsubscribe from an event using the subscription id

#### Parameters

- Subscription ID

```json
// Request
{"id": 1, "method": "eth_unsubscribe", "params": ["0x34da6f29e3e953af4d0c7c58658fd525"]}

// Result
{"jsonrpc":"2.0","result":true,"id":1}
```

## Next {hide}

Learn about the [encoding](./../core/encoding.md) formats used on Ethermint {hide}
