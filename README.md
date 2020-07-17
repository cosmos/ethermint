<!--
parent:
  order: false
-->

<div align="center">
  <h1> Ethermint </h1>
</div>

<div align="center">
  <a href="https://github.com/ChainSafe/ethermint/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/cosmos/ethermint.svg" />
  </a>
  <a href="https://github.com/ChainSafe/ethermint/blob/development/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/cosmos/ethermint.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/cosmos/ethermint?tab=doc">
    <img alt="GoDoc" src="https://godoc.org/github.com/cosmos/ethermint?status.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/cosmos/ethermint">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/cosmos/ethermint"/>
  </a>
  <a href="https://codecov.io/gh/cosmos/ethermint">
    <img alt="Code Coverage" src="https://codecov.io/gh/ChainSafe/ethermint/branch/development/graph/badge.svg"/>
  </a>
</div>
<div align="center">
  <a href="https://github.com/cosmos/ethermint">
    <img alt="Lines Of Code" src="https://tokei.rs/b1/github/cosmos/ethermint" />
  </a>
  <a href="https://discord.gg/AzefAFd">
    <img alt="Discord" src="https://img.shields.io/discord/669268347736686612.svg" />
  </a>
  <a href="https://github.com/ChainSafe/ethermint/actions?query=workflow%3ABuild">
    <img alt="Build Status" src="https://github.com/ChainSafe/ethermint/workflows/Build/badge.svg" />
  </a>
  <a href="https://github.com/ChainSafe/ethermint/actions?query=workflow%3ALint">
    <img alt="Lint Status" src="https://github.com/ChainSafe/ethermint/workflows/Lint/badge.svg" />
  </a>
</div>

Ethermint is a scalable, high-throughput Proof-of-Stake blockchain that is fully compatible and
interoperable with Ethereum. It's build using the the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) which runs on top of [Tendermint Core](https://github.com/tendermint/tendermint) consensus engine.

> **WARNING:** Ethermint is under VERY ACTIVE DEVELOPMENT and should be treated as pre-alpha software. This means it is not meant to be run in production, its APIs are subject to change without warning and should not be relied upon, and it should not be used to hold any value. We will remove this warning when we have a release that is stable, secure, and properly tested.

**Note**: Requires [Go 1.14+](https://golang.org/dl/)

## Quick Start

To learn how the Ethermint works from a high-level perspective, go to the [Introduction](./docs/intro/overview.md) section from the documentation.

For more, please refer to the [Ethermint Docs](./docs/), which are also hosted on [docs.ethermint.zone](https://docs.ethermint.zone/).

### Starting Ethermint Web3 RPC API

After the daemon is started, run (in another process):

```bash
emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key mykey
```

and to make sure the server has started correctly, try querying the current block number:

```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545
```

or point any dev tooling at `http://localhost:8545` or whatever port is chosen just as you would with an Ethereum node

#### Keyring backend options

The instructions above include commands to use `test` as the `keyring-backend`. This is an unsecured keyring that doesn't require entering a password and should not be used in production. Otherwise, Ethermint supports using a file or OS keyring backend for key storage. To create and use a file stored key instead of defaulting to the OS keyring, add the flag `--keyring-backend file` to any relevant command and the password prompt will occur through the command line. This can also be saved as a CLI config option with:

```bash
emintcli config keyring-backend file
```

### Exporting Ethereum private key from Ethermint

To export the private key from Ethermint to something like Metamask, run:

```bash
emintcli keys unsafe-export-eth-key mykey
```

Import account through private key, and to verify that the Ethereum address is correct with:

```bash
emintcli keys parse $(emintcli keys show mykey -a)
```

### Tests

Integration tests are invoked via:

```bash
make test
```

To run CLI tests, execute:

```bash
make test-cli
```

#### Ethereum Mainnet Import

There is an included Ethereum mainnet exported blockchain file in `importer/blockchain`
that includes blocks up to height `97638`. To execute and test a full import of
these blocks using the EVM module, execute:

```bash
make test-import
```

You may also provide a custom blockchain export file to test importing more blocks
via the `--blockchain` flag. See `TestImportBlocks` for further documentation.

### Community

The following chat channels and forums are a great spot to ask questions about Ethermint:

- [Cosmos Discord](https://discord.gg/W8trcGV)
- [Cosmos Forum](https://forum.cosmos.network)
