<!--
order: 3
-->

# Use Cases

## EVM Module

The EVM module packaged inside Ethermint can be used separately as its own standalone module. This can be added as a dependency to any Cosmos chain, which will allow for smart contract support. A project can integrate the EVM module and connect over a bridge network (like [Chainbridge](https://github.com/ChainSafe/ChainBridge), or a [Peg Zone](https://github.com/cosmos/peggy)) to offer full smart contract support. This will relieve those projects from creating their own smart contract functionality. The EVM module will also provide existing projects, such as enterprise projects, with fast finality as well as proof of stake consensus.

## Ethermint Client

The Ethermint client can be used to generate an EVM compatible chain and it offers high security with its fast finality. Ethermint also offers built-in interporability functionalities (under development) with its IBC ([Inter-Blockchain Communication](https://cosmos.network/ibc)) implementation. 

## Trade offs

Either option above will allow for fast finality, using a PoS consensus engine. Using the EVM module as a dependency will require the importing of the EVM and the maintaining of the chain (including validator sets, code upgrades/conformance, community engagement, and participation incetivizations). This will allow for granular control over the network and specific configurations/features that may not be available in the Ethermint client.

Using Ethermint will allow for the direct deployment of smart contracts to the Ethermint network. Utilizing the Ethermint client will defer the chain maintenance to the Ethermint network and allow for the participation in a more mature blockchain. The Ethermint client will also offer (in the near future) IBC compatibility which allows for interoperability between different network. 

<!-- please read over! -->