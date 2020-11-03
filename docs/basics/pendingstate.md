<!--
order: 6
-->

# Pending State

In Ethereum, pending blocks are generated as they are queued for production by miners. These pending blocks include pending transactions that are picked out by miners, based on the highest reward paid in gas. Ethermint is designed quite differently on this front as there is no "pending state" and blocks are produced by block producers, who include transactions into blocks in a first-in-first-out (FIFO) fashion. Transactions cannot be cherry picked out of the queue of confirmed transactions. For this reason, Etheremint does not require a pending state mechanism. However, this causes a few hiccups in terms of the queries that can be made to pending state.

## Pending State Queries

Ethermint will make queries which will account for the unconfirmed transactions present in a node's transaction mempool. The pending state query made will be subjective and the query will be made on that node's mempool. Thus, the pending state will not be the same for the same query to two different nodes. 

## RPC Calls on Pending Transactions

- [`eth_getBalance`](https://github.com/cosmos/ethermint/blob/development/docs/basics/json_rpc.md#eth_getbalance)
- [`eth_getTransactionCount`](https://github.com/cosmos/ethermint/blob/development/docs/basics/json_rpc.md#eth_gettransactioncount)
- [`eth_getBlockTransactionCountByNumber`](https://github.com/cosmos/ethermint/blob/development/docs/basics/json_rpc.md#eth_getblocktransactioncountbynumber)
- [`eth_getBlockByNumber`](https://github.com/cosmos/ethermint/blob/development/docs/basics/json_rpc.md#eth_getblockbynumber)
- [`eth_getTransactionByHash`](https://github.com/cosmos/ethermint/blob/development/docs/basics/json_rpc.md#eth_gettransactionbyhash)
- `eth_getTransactionByBlockNumberAndIndex`
- [`eth_sendTransaction`](https://github.com/cosmos/ethermint/blob/development/docs/basics/json_rpc.md#eth_sendtransaction)

