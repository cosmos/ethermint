<!--
order: 5
-->

# Upgrade Node

These instructions are for full nodes that have ran on previous versions of and would like to upgrade to the latest testnet.

## Software Upgrade

First, stop your instance of `gaiad`. Next, upgrade the software:

```bash
cd ethermint
git fetch --all && git checkout <new_version>
make install
```

::: tip
If you have issues at this step, please check that you have the latest stable version of GO installed.
:::

You will need to ensure that the version installed matches the one needed for th testnet. Check the Ethermint [release page](https://github.com/ChainSafe/ethermint/releases) for details on each release.

## Upgrade Genesis File

:::warning
If the new version you are upgrading to has breaking changes, you will have to restart your chain. If it is not breaking, you can skip to [Restart](#restart)
:::

To upgrade the genesis file, you can either fetch it from a trusted source or export it locally using the `export` command.

### Fetch from a Trusted Source

If you are joining an existing testnet, you can fetch the genesis from the appropriate testnet source/repository where the genesis file is hosted.

Save the new genesis as `new_genesis.json`. Then replace the old `genesis.json` with `new_genesis.json`

```bash
cd $HOME/.emintd/config
cp -f genesis.json new_genesis.json
mv new_genesis.json genesis.json
```

Then, go to the [reset data](#reset-data) section.

### Exporting State to a New Genesis Locally

If you were running a node in the previous version of the network and want to build your new genesis locally from a state of this previous network, use the following command:

```bash
cd $HOME/.emintd/config
emintd export --for-zero-height --height=<export-height> > new_genesis.json
```

The command above take a state at a certain height `<export-height>` and turns it into a new genesis file that can be used to start a new network. The new network will start at height 0 if the `--for-zero-height` is provided.

Then, replace the old `genesis.json` with `new_genesis.json`.

```bash
cp -f genesis.json new_genesis.json
mv new_genesis.json genesis.json
```

At this point, you might want to run a script to update the exported genesis into a genesis state that is compatible with your new version.

You can use the `migrate` command to migrate from a given version to the next one (eg: `v0.X.X` to `v1.X.X`):

```bash
emintd migrate [target-version] [/path/to/genesis.json] --chain-id=<new_chain_id> --genesis-time=<yyyy-mm-ddThh:mm:ssZ>
```

## Reset Data

First, remove the outdated files and reset the data.

```bash
rm $HOME/.emintd/config/addrbook.json $HOME/.emintd/config/genesis.json
emintd unsafe-reset-all
```

Your node is now in a pristine state while keeping the original `priv_validator.json` and `config.toml`. If you had any sentry nodes or full nodes setup before,
your node will still try to connect to them, but may fail if they haven't also
been upgraded.

::: danger Warning
Make sure that every node has a unique `priv_validator.json`. Do not copy the `priv_validator.json` from an old node to multiple new nodes. Running two nodes with the same `priv_validator.json` will cause you to double sign.
:::

### Restart

::: tip
If you are upgrading your node to a new version that is not breaking from the previous one, you can restart the chain instead of [resetting](#reset-data) the node.
:::

To restart your node once the new genesis has been updated, use the `start` command:

```bash
emintd start
```
