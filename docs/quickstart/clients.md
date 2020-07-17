<!--
order: 6
-->

# Clients

Learn how to connect a client to a running node. {synopsis}

## Pre-requisite Readings

- [Run a Node](./run_node.md) {prereq}

## Start Client

After the Ethermint daemon is started, run (in another process):

```bash
emintcli rest-server --laddr "tcp://localhost:8545" --unlock-key $KEY --chain-id $CHAINID --trace
```

You should see the logs from the REST and the RPC server.

```bash
I[2020-07-17|16:54:35.037] Starting application REST service (chain-id: "8")... module=rest-server
I[2020-07-17|16:54:35.037] Starting RPC HTTP server on 127.0.0.1:8545   module=rest-server
```

or point any dev tooling at `http://localhost:8545` or whatever port is chosen just as you would with an Ethereum node

## Client Integrations

### Command Line Interface

Ethermint is integrated with a CLI client that can be used to send transactions and query the state from each module.

```bash
# available query commands
emintcli query -h

# available transaction commands
emintcli tx -h
```

### REST server

### Ethereum JSON-RPC
