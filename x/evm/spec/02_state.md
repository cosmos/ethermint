<!--
order: 2
-->

# State


<!-- TODO: structs that are set on the evm store -->

## `CommitStateDB`


## State Objects

State objects are used by the VM which is unable to deal with database-level errors. Any error that occurs during a database read is memoized here and will eventually be returned by `StateDB.Commit`.

The Ethermint `stateObject` is a concrete type that mimics the functionality from the `go-ethereum`
private `stateObject` type. It keeps track of the interim values for the contract bytecode, storage
state and balance of an `EthAccount`.

The storage entries (original and "dirty") for each state object are represented as slices instead
of maps since latter can cause non-deterministic block app hashes, which result in the chain
halting.

When a `stateObject` is committed during `EndBlock`. It sets sets the account contract code to store, as well as the dirty storage state. The account's nonce and the account balance are updated by calling the `auth` and `bank` module setter functions, respectively.

<!-- TODO: paste code from Ethermint State Object -->

The functionalities provided by the Ethermint `stateObject` are:

* Storage state getter and setter (temporary)
* Contract bytecode getter and setter (temporary)
* Balance getter and setter (temporary)
* Balance accounting (temporary)
* Account nonce and address getter and setter (temporary)
* Auxiliary functions: copy, RLP encoding, empty
* Commit state object (final)
