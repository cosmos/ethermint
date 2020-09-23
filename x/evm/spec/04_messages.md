<!--
order: 4
-->

# Messages

## MsgEthereumTx

An EVM state transition can be achieved by using the `MsgEthereumTx`:

```go
// MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
type MsgEthereumTx struct {
  Data TxData

  // caches
  size atomic.Value
  from atomic.Value
}
```

```go
// TxData implements the Ethereum transaction data structure. It is used
// solely as intended in Ethereum abiding by the protocol.
type TxData struct {
  AccountNonce uint64          `json:"nonce"`
  Price        *big.Int        `json:"gasPrice"`
  GasLimit     uint64          `json:"gas"`
  Recipient    *ethcmn.Address `json:"to" rlp:"nil"` // nil means contract creation
  Amount       *big.Int        `json:"value"`
  Payload      []byte          `json:"input"`

  // signature values
  V *big.Int `json:"v"`
  R *big.Int `json:"r"`
  S *big.Int `json:"s"`

  // hash is only used when marshaling to JSON
  Hash *ethcmn.Hash `json:"hash" rlp:"-"`
}
```

This message is expected to fail if:

- `Data.Price` (i.e gas price) is ≤ 0.
- `Data.Amount` is negative
