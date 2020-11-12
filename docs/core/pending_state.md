<!--
order: 6
-->

# Pending State

## Ethermint vs Ethereum

In Ethereum, pending blocks are generated as they are queued for production by miners. These pending
blocks include pending transactions that are picked out by miners, based on the highest reward paid
in gas. This mechanism exists as block finality is not possible on the Ethereum network. Blocks are
committed with probabilistic finality, which means that transactions and blocks become less likely
to become reverted as more time (and blocks) passes.

Ethermint is designed quite differently on this front as there is no concept of a "pending state"
and blocks are produced by block producers, who include transactions into blocks in a
first-in-first-out (FIFO) fashion. Transactions cannot be cherry picked out of the queue of
confirmed transactions. All of this can be achieved because Ethermint (using Cosmos-SDK) offers
instant finality with its PoS based consensus algorithm. For this reason, Etheremint does not
require a pending state mechanism. However, this causes a few hiccups in terms of the queries that
can be made to pending state.

## Pending State Queries

Ethermint will make queries which will account for any unconfirmed transactions present in a node's
transaction mempool. A pending state query made will be subjective and the query will be made on the
target node's mempool. Thus, the pending state will not be the same for the same query to two
different nodes.

## RPC Calls on Pending Transactions

- [`eth_getBalance`](./../basics/json_rpc.md#eth_getbalance)
- [`eth_getTransactionCount`](./../basics/json_rpc.md#eth_gettransactioncount)
- [`eth_getBlockTransactionCountByNumber`](./../basics/json_rpc.md#eth_getblocktransactioncountbynumber)
- [`eth_getBlockByNumber`](./../basics/json_rpc.md#eth_getblockbynumber)
- [`eth_getTransactionByHash`](./../basics/json_rpc.md#eth_gettransactionbyhash)
- `eth_getTransactionByBlockNumberAndIndex`
- [`eth_sendTransaction`](./../basics/json_rpc.md#eth_sendtransaction)
