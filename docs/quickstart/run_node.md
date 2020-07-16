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

## Restart

To restart your node, just type:

```bash
emintd start
```

::: tip
If you are upgrading your node to a new version that is not breaking from the previous one, you can restart the chain instead of [resetting](#reset-data) the node.
:::

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

## Clearing data from chain

### Reset Data

Alternatively, you can **reset** the blockchain database, and remove the node's address book files, and resets `priv_validator.json` to the genesis state.

::: danger
If you are running a **validator node**, always be careful when doing `emintd unsafe-reset-all`. You should never use this command if you are not switching `chain-id`.
:::

::: danger
**IMPORTANT**: Make sure that every node has a unique `priv_validator.json`. **Do not** copy the `priv_validator.json` from an old node to multiple new nodes. Running two nodes with the same `priv_validator.json` will cause you to double sign!
:::

First, remove the outdated files and reset the data.

```bash
emintd unsafe-reset-all
```

Your node is now in a pristine state while keeping the original `priv_validator.json` and `config.toml`. If you had any sentry nodes or full nodes setup before, your node will still try to connect to them, but may fail if they haven't also been upgraded.

### Delete Data

Data for the Daemon and CLI binaries should be stored at `~/.emintd` and `~/.emintcli`, respectively by default. To **delete** the existing binaries and configuration, run:

```bash
rm -rf ~/.emint*
```

To clear all data except key storage (if keyring backend chosen) and then you can rerun the full node installation commands from above to start the node again.

## Next {hide}

Learn about Ethermint [accounts](./../basic/accounts.md) {hide}
