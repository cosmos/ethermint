<!--
order: 7
-->

# Notice

Important information for Etheremint development.

## Cross-Chain Transaction Replay Attack Prevention

Given the fact that Ethermint uses the same transaction signature format as Ethereum for compatibility with Ethereum tooling, and Ethermint's chainIDEpoch corresponds to the chain ID on Ethereum,
any valid transaction that is able to broadcast to Ethereum is also able to broadcast to Ethermint, and vice versa, which could lead to the Transaction Replay Attack if Ethermint starts the network by using Ethereum reserved chainIDs.

To prevent this issue, simply choose a different chainID other than [chainIDs that are reserved for Ethereum](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md#list-of-chain-ids).
