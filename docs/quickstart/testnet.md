<!--
order: 3
-->

# Testnet

Learn how to deploy a local testnet or connect to an existing public one {synopsis}

## Pre-requisite Readings

- [Install Aragon-Chain](./installation.md) {prereq}
- [Install Docker](https://docs.docker.com/engine/installation/)  {prereq}
- [Install docker-compose](https://docs.docker.com/compose/install/)  {prereq}
<!-- - [Install `jq`](https://stedolan.github.io/jq/download/) {prereq} -->

## Single-node, Local, Manual Testnet

This guide helps you create a single validator node that runs a network locally for testing and other development related uses.

### Initialize node

```bash
$MONIKER=testing
$KEY=mykey
$CHAINID=8

emintd init $MONIKER --chain-id=$CHAINID
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
# specified in this config (e.g. 10photon).

minimum-gas-prices = ""
```

### Genesis Procedure

```bash
# Create a key to hold your account
emintcli keys add $KEY

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some coins
# with the genesis.app_state.staking.params.bond_denom denom, the default is staking
emintd add-genesis-account $(emintcli keys show validator -a) 1000000000stake,10000000000photon

# Generate the transaction that creates your validator
emintd gentx --name $KEY

# Add the generated bonding transaction to the genesis file
emintd collect-gentxs

# Finally, check the correctness of the genesis.json file
emintd validate-genesis
```

### Run Testnet

Now its safe to start the daemon:

```bash
emintd start
```

## Multi-node, Local, Automated Testnet

### Build & Run Testnet

To build start a 4 node testnet run:

```bash
make localnet-start
```

This command creates a 4-node network using the `emintdnode` Docker image.
The ports for each node are found in this table:

| Node ID      | P2P Port | REST/RPC Port |
|--------------|----------|---------------|
| `emintnode0` | `26656`  | `26657`       |
| `emintnode1` | `26659`  | `26660`       |
| `emintnode2` | `26661`  | `26662`       |
| `emintnode3` | `26663`  | `26664`       |

To update the binary, just rebuild it and restart the nodes:

```bash
make localnet-start
```

### Stop Testnet

Once you are done, execute:

```bash
make localnet-stop
```

### Configuration

The `make localnet-start` creates files for a 4-node testnet in `./build` by
calling the `emintd testnet` command. This outputs a handful of files in the
`./build` directory:

```bash
tree -L 3 build/

build/
├── emintcli
├── emintd
├── gentxs
│   ├── node0.json
│   ├── node1.json
│   ├── node2.json
│   └── node3.json
├── node0
│   ├── emintcli
│   │   ├── key_seed.json
│   │   └── keyring-test-cosmos
│   └── emintd
│       ├── config
│       ├── data
│       └── emintd.log
├── node1
│   ├── emintcli
│   │   ├── key_seed.json
│   │   └── keyring-test-cosmos
│   └── emintd
│       ├── config
│       ├── data
│       └── emintd.log
├── node2
│   ├── emintcli
│   │   ├── key_seed.json
│   │   └── keyring-test-cosmos
│   └── emintd
│       ├── config
│       ├── data
│       └── emintd.log
└── node3
    ├── emintcli
    │   ├── key_seed.json
    │   └── keyring-test-cosmos
    └── emintd
        ├── config
        ├── data
        └── emintd.log
```

Each `./build/nodeN` directory is mounted to the `/emintd` directory in each container.

### Logging

Logs are saved under each `./build/nodeN/emintd/gaia.log`. You can also watch logs
directly via Docker, for example:

```bash
docker logs -f emintdnode0
```

### Keys & Accounts

To interact with `emintcli` and start querying state or creating txs, you use the
`emintcli` directory of any given node as your `home`, for example:

```bash
emintcli keys list --home ./build/node0/emintcli
```

Now that accounts exists, you may create new accounts and send those accounts
funds!

::: tip
**Note**: Each node's seed is located at `./build/nodeN/emintcli/key_seed.json` and can be restored to the CLI using the `emintcli keys add --restore` command
:::

### Special Binaries

If you have multiple binaries with different names, you can specify which one to run with the BINARY environment variable. The path of the binary is relative to the attached volume. For example:

```bash
# Run with custom binary
BINARY=aragon-chain make localnet-start
```

## Multi-node, Public, Manual Testnet

If you are looking to connect to a persistent public testnet. You will need to manually configure your node.

### Genesis and Seeds

#### Copy the Genesis File

::: tip
If you want to start a network from scratch, you will need to start the [genesis procedure](#genesis-procedure) by creating a `genesis.json` and submit + collect the genesis transactions from the [validators](./validator-setup.md).
:::

If you want to connect to an existing testnet, fetch the testnet's `genesis.json` file and copy it into the `emintd`'s config directory (i.e `$HOME/.emintd/config/genesis.json`).

Then verify the correctness of the genesis configuration file:

```bash
emintd validate-genesis
```

#### Add Seed Nodes

Your node needs to know how to find peers. You'll need to add healthy seed nodes to `$HOME/.emintd/config/config.toml`. If those seeds aren't working, you can find more seeds and persistent peers on an existing explorer.

For more information on seeds and peers, you can the Tendermint [P2P documentation](https://docs.tendermint.com/master/spec/p2p/peer.html).

#### Start testnet

The final step is to [start the nodes](./run_node.md#start-node). Once enough voting power (+2/3) from the genesis validators is up-and-running, the testnet will start producing blocks.

## Testnet faucet

Once the aragon-chain daemon is up and running, you can request tokens to your address using the `faucet` module:

```bash
# query your initial balance
emintcli q bank balances $(emintcli keys show <mykey> -a)  

# send a tx to request tokens to your account address
emintcli tx faucet request 100photon --from <mykey>

# query your balance after the request
emintcli q bank balances $(emintcli keys show <mykey> -a)
```

You can also check to total amount funded by the faucet and the total supply of the chain via:

```bash
# total amount funded by the faucet
emintcli q faucet funded

# total supply
emintcli q supply total
```

## Next {hide}

Learn about how to setup a [validator](./validator-setup.md) node on Aragon-Chain {hide}
