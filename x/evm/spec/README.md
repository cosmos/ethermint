<!--
order: 0
title: EVM Overview
parent:
  title: "evm"
-->

# `evm`

## Abstract

This document defines the specification of the Ethereum Virtual Machine (EVM) as a Cosmos SDK module.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[State Transitions](03_state_transitions.md)**
4. **[Messages](04_messages.md)**
5. **[ABCI](05_abci.md)**
6. **[Events](06_events.md)**
7. **[Parameters](07_params.md)**

## Module Architecture

> **NOTE for auditors**: If you're not familiar with the overall module structure from
the SDK modules, please check this [document](https://docs.cosmos.network/master/building-modules/structure.html) as
prerequisite reading.

```shell
evm/
├── client
│   └── cli
│       ├── query.go
│       └── tx.go
├── keeper
│   ├── keeper.go
│   └── querier.go
├── types
│   ├── chain_config.go
│   ├── codec.go          # Type registration for encoding
│   ├── errors.go         # Module-specific errors
│   ├── events.go         # Events exposed to the Tendermint PubSub/Websocket
│   ├── genesis.go        # Genesis state for the module
│   ├── journal.go        # Ethereum Journal of state transitions
│   ├── msg.go            # EVM module transaction messages
│   ├── params.go         # Module parameters that can be customized with governance parameter change proposals
│   ├── keys.go           # Store keys and utility functions
│   └── tx_data.go        # Ethereum transaction data types
├── abci.go               # ABCI BeginBlock and EndBlock logic
├── handler.go            # Message routing
└── module.go             # Module setup for the module manager
```
