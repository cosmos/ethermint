<!--
order: 3
-->

# Store and Database

Learn about the store and database concept on Ethermint. {synopsis}

Since Ethermint is using Cosmos-sdk as ABCI application, when starting an Ethermint node,
a Key-Value database will be initialized and all Ethermint blockchain data will be persisted into the node machine.


## Store

Store is the abstract data structure in Ethermint to organize and persist blockchain data into physical disk.

There are two kinds of store in Ethermint:  

1. Cosmos SDK MultiStore.
2. Ethermint CommitStateDB.

Cosmos SDK MultiStore is used in most of Ethermint modules such as bank, governance and staking. The synchronization of this store among network is taken care by [Tendermint](https://github.com/tendermint/tendermint).

Different store prefix(store key) will be assigned to each Ethermint module since they all implement cosmos-sdk MultiStore interface. Cosmos MultiStore can be accessed by `sdk.Context` expect `ethermint/evm` module.

For `evm` module, Ethermint creates a type called `CommitStateDB` which implements `ethvm.StateDB` from go-ethereum lib, and it aims to be compatible with ethereum EVM transactions.

## Database

Database is the actual Key-Value database dependency that is running on the backend.  

Currently, Ethermint supports the [databases](https://github.com/tendermint/tm-db) that Tendermint supports.

Developers can modify the database config in `~/.ethermintd/config.toml -> db_backend` to enable other supported database.

## Next {hide}

Learn how to deploy a Solidity smart contract on Ethermint using [Truffle](./../guides/truffle.md) {hide}

