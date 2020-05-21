<!--
order: 4
-->

# Gas and Fees

This document describes the default strategies to handle gas and fees within a Cosmos SDK application. {synopsis}

## Introduction to `Gas` in the SDK

<!-- TODO: -->

## Matching EVM Gas consumption

<!-- TODO: -->

## Gas refunds

<!-- TODO: -->

## AnteHandler

The `AnteHandler` is a special `handler` that is run for every transaction during `CheckTx` and `DeliverTx`, before the `handler` of each `message` in the transaction. `AnteHandler`s have a different signature than `handler`s

<!-- TODO: -->

## Next {hide}

Learn more about the [Lifecycle of a transaction](./tx-lifecycle.md) {hide}
