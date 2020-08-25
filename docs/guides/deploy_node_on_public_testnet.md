<!--
order: 5
-->

# Deploy node to existing testnet

Learn how to deploy and connect a node in the cloud to public testnet. {synopsis}

## Pre-requisite Readings

- [Deploy Testnet on cloud](./testnet_on_cloud.md) {prereq}


## Deploy node to Digital Ocean

### Create a Droplet

Create a new droplet using the same steps in [Deploy Testnet on cloud](./testnet_on_cloud.md). 

Once this new droplet is created make sure to get its IP address to be used in the next steps.

### Connect to Droplet

Click on the started Droplet, and you'll see details about it. At the moment, we're interested in the IP address - this is the address that the Droplet is at on the internet.

To access it, we'll need to connect to it using our previously created private key. From the same folder as that private key, run:

```bash
$ ssh -i digital-ocean-key root@<IP_ADDRESS>
```

Now you are connected to the droplet. 

### Install Ethermint

Clone and build Ethermint in the droplet using `git`:

```bash
go install https://github.com/ChainSafe/ethermint.git
```

Check that the binaries have been successfuly installed:

```bash
ethermintd -h
ethermintcli -h
```

#### Copy the Genesis File

To connect the node to the existing testnet, fetch the testnet's `genesis.json` file and copy it into the new droplets config directory (i.e `$HOME/.ethermintd/config/genesis.json`).

To do this ssh into both the testnet droplet and the new node droplet. 

On your local machine copy the genesis.json file from the testnet droplet to the new droplet using 

```bash
scp -3 root@<TESTNET_IP_ADDRESS>:$HOME/.ethermintd/config/genesis.json root@<NODE_IP_ADDRESS>:$HOME/.ethermintd/config/genesis.json
```

### Start the Node

Once the genesis file is copied over run `ethermind start` inside the node droplet. 