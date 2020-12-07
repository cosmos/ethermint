<!--
order: 1
-->

# Concepts

## EVM

EVM is the Ethereum Virtual Machine that the necessary tools to run or create a contract on a given state.

## State DB

The `StateDB` interface from geth represents an EVM database for full state querying of both contracts and accounts. The concrete type that fulfills this interface on Ethermint is the `CommitStateDB`.

## State Object

`stateObject` represents an Ethereum account which is being modified.
The usage pattern is as follows:

First you need to obtain a state object.
Account values can be accessed and modified through the object.

## Genesis State

The `x/evm` module `GenesisState` defines the state necessary for initializing the chain from a previous exported height.

```go
// GenesisState defines the evm module genesis state
type GenesisState struct {
  Accounts    []GenesisAccount  `json:"accounts"`
  TxsLogs     []TransactionLogs `json:"txs_logs"`
  ChainConfig ChainConfig       `json:"chain_config"`
  Params      Params            `json:"params"`
}
```

### Genesis Accounts

The `GenesisAccount` type corresponds to an adaptation of the Ethereum `GenesisAccount` type. Its
main difference is that the one on Ethermint uses a custom `Storage` type that uses a slice instead of maps for the evm `State`,and that it doesn't contain the private key field.

It is also important to note that since the `auth` and `bank` SDK modules manage the accounts and balance state,  the `Address` must correspond to an `EthAccount` that is stored in the `auth`'s module `AccountKeeper` and the balance must match the balance of the `EvmDenom` token denomination  defined on the `GenesisState`'s `Param`. The values for the address and the balance amount maintain the same format as the ones from the SDK to make manual inspections easier on the genesis.json.

```go
type GenesisAccount struct {
  Address string        `json:"address"`
  Balance sdk.Int       `json:"balance"`
  Code    hexutil.Bytes `json:"code,omitempty"`
  Storage Storage       `json:"storage,omitempty"`
}
```

### Transaction Logs

On every Ethermint transaction, its result contains the Ethereum `Log`s from the state machine
execution that are used by the JSON-RPC Web3 server for for filter querying. Since Cosmos upgrades
don't persist the transactions on the blockchain state, we need to persist the logs the EVM module
state to prevent the queries from failing.

`TxsLogs` is the field that contains all the transaction logs that need to be persisted after an upgrade. It uses an array instead of a map to ensure determinism on the iteration.

```go
type TransactionLogs struct {
  Hash ethcmn.Hash     `json:"hash"`
  Logs []*ethtypes.Log `json:"logs"`
}
```

### Chain Config

The `ChainConfig` is a custom type that contains the same fields as the go-ethereum `ChainConfig`
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

  YoloV2Block sdk.Int `json:"yoloV2_block" yaml:"yoloV2_block"` // YOLO v2: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
  EWASMBlock  sdk.Int `json:"ewasm_block" yaml:"ewasm_block"`   // EWASM switch block (< 0 no fork, 0 = already activated)
}
```

### Params

See the [params](07_params.md) document for further information about parameters.
