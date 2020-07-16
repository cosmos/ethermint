<!--
order: 2
-->

# Run a Node

Run a local node and start the REST and JSON-RPC clients {synopsis}

## Script deployment

Run the local node with faucet enabled:

::: warning
The script below will remove any pre-existing binaries installed. Use the manual deploy if you want
to keep your binaries and configuration files.
:::

```bash
./init.sh
```

In another terminal window or tab, run the Ethereum JSON-RPC server as well as the SDK REST server:

```bash
emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key mykey --chain-id 8
```

## Manual setup

These instructions are for setting up a brand new full node from scratch.

First, initialize the node and create the necessary config files:

```bash
emintd init <your_custom_moniker>
```

::: warning
Monikers can contain only ASCII characters. Using Unicode characters will render your node unreachable.
:::

You can edit this `moniker` later, in the `$(HOME)/.emintd/config/config.toml` file:

```toml
# A custom human readable name for this node
moniker = "<your_custom_moniker>"
```

You can edit the `$HOME/.emintd/config/app.toml` file in order to enable the anti spam mechanism and reject incoming transactions with less than the minimum gas prices:

```toml
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### main base config options #####

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 10uatom).

minimum-gas-prices = ""
```

Your full node is now initiallized.

## Key Management

To run a node with the same key every time:
replace `emintcli keys add $KEY` in `./init.sh` with:

```bash
echo "your mnemonic here" | emintcli keys add $KEY --recover
```

::: tip
Ethermint currently only supports 24 word mnemonics.
:::

You can generate a new key/mnemonic with

```bash
emintcli keys add $KEY
```

To export your ethermint key as an ethereum private key (for use with Metamask for example):

```bash
emintcli keys unsafe-export-eth-key $KEY
```

## Next {hide}

Learn about Ethermint [accounts](./../basic/accounts.md) {hide}
