# Ethermint

Ethermint is a scalable, high-throughput Proof-of-Stake blockchain that is fully compatible and
interoperable with Ethereum's virtual machine and soidity smart contract programming language. It's built using the the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) which runs on top of [Tendermint Core](https://github.com/tendermint/tendermint) consensus engine.

> **WARNING:** Ethermint is under VERY ACTIVE DEVELOPMENT and should be treated as pre-alpha software. This means it is not meant to be run in production, its APIs are subject to change without warning and should not be relied upon, and it should not be used to hold any value. We will remove this warning when we have a release that is stable, secure, and properly tested.

**Note**: Requires [Go 1.16+](https://golang.org/dl/)

## Quick Start

To learn how the Ethermint works from a high-level perspective, go to the [Introduction](./docs/intro/overview.md) section from the documentation.

For more, please refer to the [Ethermint Docs](./docs/), which are also hosted on [docs.ethermint.zone](https://docs.ethermint.zone/).

## Tests

Unit tests are invoked via:

```bash
make test
```

To run JSON-RPC tests, execute:

```bash
make test-rpc
```

There is also an included Ethereum mainnet exported blockchain file in `importer/blockchain`
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
