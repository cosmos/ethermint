<!--
order: 3
-->

# Testnet

Learn how to deploy a local testnet or connect to an existing one {synopsis}

## Starting a New Node

> NOTE: If you ran a full node on a previous testnet, please skip to [Upgrading From Previous Testnet](#upgrading-from-previous-testnet).

To start a new node, the mainnet instructions apply:

- [Join the mainnet](./join-mainnet.md)
- [Deploy a validator](../validators/validator-setup.md)

The only difference is the SDK version and genesis file. See the [testnet repo](https://github.com/cosmos/testnets) for information on testnets, including the correct version of the Cosmos-SDK to use and details about the genesis file.

## Upgrading Your Node

These instructions are for full nodes that have ran on previous versions of and would like to upgrade to the latest testnet.

### Reset Data

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

## Requesting tokens though the testnet faucet

Once the ethermint daemon is up and running, you can request tokens to your address using the `faucet` module:

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

Learn about Ethermint [accounts](./../basic/accounts.md) {hide}
