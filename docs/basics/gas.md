<!--
order: 3
-->

# Gas and Fees

Learn about the differences between `Gas` and `Fees` in Ethereum and Cosmos. {synopsis}

The concept of gas in Ethereum was created to disallow the EVM (Ethereum Virtual Machine) from running infinite loops by allocating a small amount of monetary value into the system. A unit of gas, usually in a form as a fraction of the native coin, is consumed for every operation on the EVM and requires a user to pay for these operations. These operations can be anything done that changes the state of the EVM will require gas to do (e.g. sending a transaction, contract call). 

Exactly like Ethereum, Cosmos utilizes the concept of gas and this is how Cosmos tracks the resource usage of operations during execution. Operations on Cosmos are any read or writes done to the store. In Cosmos, a fee is calculated and charged to the user during a message execution. This fee is calculated from the sum of all gas consumed in an message execution (fee = gas * gas-price). 

In both networks, gas is used to make sure that operations do not require an excess amount of computational power to complete and as a way to deter bad-acting users from spamming the network. 

## Introduction to `Gas` in the SDK

In the Cosmos SDK, gas is tracked in the `main gas meter` and the `block gas meter`. 

The `main gas meter` keeps track of the gas consumed during executions that lead to state transitions.
The `block gas meter` keeps track of the gas consumed in a block and enforces that the gas does not go over a predefined limit. 

More information regarding gas in Cosmos SDK can be found [here](https://docs.cosmos.network/master/basics/gas-fees.html).

## Matching EVM Gas consumption

As Ethermint is intended to simulate the EVM, gas consumption must be equitable in order to accurately calculate the state transition hashes and exact the behaviour that would be seen on the main Ethereum network (main net). The gas calculated in Ethermint is done by go-ethereum's `IntrinsicGas` method. This allows Ethermint to generate the expected gas costs for operations done in the network and scale the gas costs as it would in the Ethereum network. 

<!-- need someone to read over -->

## Gas Refunds

In Ethereum, gas can be specified prior to execution and the remaining gas will be refunded back to the user if any gas is left over - should fail with out of gas if not enough gas was provided. In Ethermint, the concept of gas refunds does not exist and the fees paid is not refunded in part back to the user. The fees exacted on a transaction will be collected by the validator and no refunds are issued. Thus, it is extremely important to use the correct gas. 

To prevent overspending on fees, providing the `--gas-adjustment` flag for a cosmos transactions will determine the fees automatically. Also the `eth_estimateGas` rpc call can be used to manually get the correct gas costs for a transaction.

## 0 Fee Transactions

In Ethermint, it is possible to send transactions with 0 fees. The fees are determined by the validator of the network, and each validator can specify a different value for their fees. 

EVM transactions cannot have 0 fees as gas is required inherently. Cosmos-SDK related transactions, such as staking, would be transactions that are eligible for 0 fee execution, given that the validator has specified 0 fees.

## AnteHandler

The `AnteHandler` is a special `handler` that is run for every transaction during `CheckTx` and `DeliverTx`, before the `handler` of each `message` in the transaction. `AnteHandler`s have a different signature than `handler`s

<!-- TODO: -->

## Next {hide}

Learn about the [Photon](./photon.md) token {hide}
