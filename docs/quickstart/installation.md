<!--
order: 1
-->

# Installation

## Binaries

Clone and build Ethermint using `git`:

```bash
git clone https://github.com/ChainSafe/ethermint.git
cd ethermint
make install
```

Check that the binaries have been successfuly installed:

```bash
emintd -h
emintcli -h
```

<!-- ## Docker -->

<!-- TODO: -->

## Releases

::: warning
Ethermint is under VERY ACTIVE DEVELOPMENT and should be treated as pre-alpha software. This means it is not meant to be run in production, its APIs are subject to change without warning and should not be relied upon, and it should not be used to hold any value. We will remove this warning when we have a release that is stable, secure, and properly tested.
:::

You can also download a specific release available on the [Ethermint repository](https://github.com/ChainSafe/ethermint/releases)

## Next {hide}

Learn how to [run a node](./.run_node.md) {hide}
