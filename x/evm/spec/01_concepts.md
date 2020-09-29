<!--
order: 1
-->

# Concepts

## EVM

EVM is the Ethereum Virtual Machine that the necessary tools to run or create a contract on a given state.

## State DB

The `StateDB` is represents an EVM database for full state querying of both contracts and accounts.

## State Object

## Chain Config

The `ChainConfig` is a custom type that contains the same fields as the go-ethereum ChainConfig
parameters, but using `sdk.Int` types instead of `*big.Int`. It also defines additional YAML tags
for pretty printing.

The `ChainConfig` type is not a configurable SDK `Param` since the SDK does not allow for validation
against a previous stored parameter values or `Context` fields. Since most of this type's fields
rely on the block height value, this limitation prevents the validation of of potential new
parameter values against the current block height (eg: to prevent updating the config block values
to a past block).

If you want to update the config values, use an software upgrade procedure.

```go
type ChainConfig struct {
  HomesteadBlock sdk.Int `json:"homestead_block" yaml:"homestead_block"` // Homestead switch block (< 0 no fork, 0 = already homestead)

  DAOForkBlock   sdk.Int `json:"dao_fork_block" yaml:"dao_fork_block"`     // TheDAO hard-fork switch block (< 0 no fork)
  DAOForkSupport bool    `json:"dao_fork_support" yaml:"dao_fork_support"` // Whether the nodes supports or opposes the DAO hard-fork

  // EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
  EIP150Block sdk.Int `json:"eip150_block" yaml:"eip150_block"` // EIP150 HF block (< 0 no fork)
  EIP150Hash  string  `json:"eip150_hash" yaml:"eip150_hash"`   // EIP150 HF hash (needed for header only clients as only gas pricing changed)

  EIP155Block sdk.Int `json:"eip155_block" yaml:"eip155_block"` // EIP155 HF block
  EIP158Block sdk.Int `json:"eip158_block" yaml:"eip158_block"` // EIP158 HF block

  ByzantiumBlock      sdk.Int `json:"byzantium_block" yaml:"byzantium_block"`           // Byzantium switch block (< 0 no fork, 0 = already on byzantium)
  ConstantinopleBlock sdk.Int `json:"constantinople_block" yaml:"constantinople_block"` // Constantinople switch block (< 0 no fork, 0 = already activated)
  PetersburgBlock     sdk.Int `json:"petersburg_block" yaml:"petersburg_block"`         // Petersburg switch block (< 0 same as Constantinople)
  IstanbulBlock       sdk.Int `json:"istanbul_block" yaml:"istanbul_block"`             // Istanbul switch block (< 0 no fork, 0 = already on istanbul)
  MuirGlacierBlock    sdk.Int `json:"muir_glacier_block" yaml:"muir_glacier_block"`     // Eip-2384 (bomb delay) switch block (< 0 no fork, 0 = already activated)

  YoloV1Block sdk.Int `json:"yoloV1_block" yaml:"yoloV1_block"` // YOLO v1: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
  EWASMBlock  sdk.Int `json:"ewasm_block" yaml:"ewasm_block"`   // EWASM switch block (< 0 no fork, 0 = already activated)
}
```

## Genesis State

The `x/evm` module `GenesisState` defines the state necessary for initializing the chain from a previous exported height.

<!-- TODO: write about txs logs persistence -->

```go
// GenesisState defines the evm module genesis state
type GenesisState struct {
  Accounts    []GenesisAccount  `json:"accounts"`
  TxsLogs     []TransactionLogs `json:"txs_logs"`
  ChainConfig ChainConfig       `json:"chain_config"`
  Params      Params            `json:"params"`
}
```

The `GenesisAccount` type corresponds to an adaptation of the Ethereum `GenesisAccount` type. Its
main difference is that the one on Ethermint uses a custom `Storage` type that uses a slice instead of maps for the evm `State`,and that it doesn't contain the private key field.

It is also important to note that since the `auth` and `bank` SDK modules manage the accounts and balance state,  the `Address` must correspond to an `EthAccount` that is stored in the auth AccountKeeper and the balance must match the balance of the `EvmDenom` token denomination  defined on the `GenesisState`'s `Param`.

```go
type GenesisAccount struct {
  Address ethcmn.Address `json:"address"`
  Balance *big.Int       `json:"balance"`
  Code    hexutil.Bytes  `json:"code,omitempty"`
  Storage Storage        `json:"storage,omitempty"`
}
```
