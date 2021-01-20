<!--
order: 3
-->

# Store and Database

Learn about the store used on Ethermint. {synopsis}

Since Ethermint is using Cosmos-sdk as ABCI application, when starting an Ethermint node,
a Key-Value database will be initial and all Ethermint blockchain data will be persisted to the node machine.


## Store

There are two kinds of store in Ethermint:  

1. Cosmos SDK MultiStore
2. Ethermint CommitStateDB

Cosmos SDK MultiStore is used in most of Ethermint modules such as bank, governance and staking. The synchronization of this store among network is taken care by [Tendermint](https://github.com/tendermint/tendermint).

Different store prefix(store key) will be assigned to each Ethermint module since they all implement cosmos-sdk MultiStore interface. Cosmos MultiStore can be accessed by `sdk.Context` expect for `ethermint/evm` module.

For `evm` module, Ethermint creates a type `CommitStateDB` which implements `ethvm.StateDB` from go-ethereum lib which aims to be compatible with ethereum EVM.

## Database

Currently, Ethermint supports the [databases](https://github.com/tendermint/tm-db) that Tendermint supports.

Developers can modify the database config in `{HOME}/.ethermintd/config.toml -> db_backend` to enable other supported database.


