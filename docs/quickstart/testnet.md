<!--
order: 3
-->

# Testnet

Learn how to deploy a local testnet or connect to an existing one {synopsis}

## Pre-requisite Readings

- [Run Node](./run_node.md) {prereq}

## Add Seed Nodes

Your node needs to know how to find peers. You'll need to add healthy seed nodes to `$HOME/.emintd/config/config.toml`. If those seeds aren't working, you can find more seeds and persistent peers on an existing explorer.

For more information on seeds and peers, you can the Tendermint [P2P documentation](https://docs.tendermint.com/master/spec/p2p/peer.html).

## Testnet faucet

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

Learn about how to setup a [validator](./validator-setup.md) on Ethermint {hide}
