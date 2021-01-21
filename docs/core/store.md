<!--
order: 3
-->

# Store and Database

Learn about the store and database concept on Ethermint. {synopsis}

Since Ethermint uses the Cosmos-SDK as an ABCI application, when starting an Ethermint node,
a Key-Value database will be initialized and all Ethermint blockchain data will be persisted into the node host machine.


## Store

Store is the abstract data structure in Ethermint to organize and persist blockchain data into a physical disk.

There are two kinds of store in Ethermint:  

1. Cosmos SDK MultiStore.
2. Ethermint CommitStateDB.

Cosmos SDK MultiStore is used in most Ethermint modules such as bank, governance and staking. The synchronization of this store on the network is taken care of by [Tendermint](https://github.com/tendermint/tendermint).

A unique store prefix(also referred to as the store key) will be assigned to each Ethermint module since they all implement the cosmos-sdk MultiStore interface. Cosmos MultiStore can be accessed by `sdk.Context` expect `ethermint/evm` module.

For `evm` module, Ethermint creates a type called `CommitStateDB` which implements `ethvm.StateDB` interface from go-ethereum, and it aims to be compatible with ethereum EVM transactions.

## Database

Database is the actual Key-Value database dependency that is running on the backend.  

Currently, Ethermint supports the [databases](https://github.com/tendermint/tm-db) that Tendermint supports.

Developers can modify the database config in `~/.ethermintd/config.toml -> db_backend` to enable another supported database.

## Next {hide}

Learn how to deploy a Solidity smart contract on Ethermint using [Truffle](./../guides/truffle.md) {hide}

