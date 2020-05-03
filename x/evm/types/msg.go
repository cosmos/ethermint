package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ethermint/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	_ sdk.Msg = MsgEthermint{}
	_ sdk.Msg = MsgEthereumTx{}
	_ sdk.Tx  = MsgEthereumTx{}
)

var big8 = big.NewInt(8)

// message type and route constants
const (
	// TypeMsgEthereumTx defines the type string of an Ethereum tranasction
	TypeMsgEthereumTx = "ethereum"
	// TypeMsgEthermint defines the type string of Ethermint message
	TypeMsgEthermint = "ethermint"
)

// NewMsgEthermint returns a reference to a new Ethermint transaction
func NewMsgEthermint(
	nonce uint64, to sdk.AccAddress, amount sdk.Int,
	gasLimit uint64, gasPrice sdk.Int, payload []byte, from sdk.AccAddress,
) MsgEthermint {
	return MsgEthermint{
		AccountNonce: nonce,
		Price:        sdk.IntProto{Int: gasPrice},
		GasLimit:     gasLimit,
		Recipient:    to,
		Amount:       sdk.IntProto{Int: amount},
		Payload:      payload,
		From:         from,
	}
}

// Route should return the name of the module
func (msg MsgEthermint) Route() string { return RouterKey }

// Type returns the action of the message
func (msg MsgEthermint) Type() string { return TypeMsgEthermint }

// GetSignBytes encodes the message for signing
func (msg MsgEthermint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgEthermint) ValidateBasic() error {

	if msg.Price.Int.Sign() != 1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "price must be positive %s", msg.Price.Int)
	}

	// Amount can be 0
	if msg.Amount.Int.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", msg.Amount.Int)
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgEthermint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg MsgEthermint) To() *ethcmn.Address {
	if msg.Recipient == nil {
		return nil
	}

	addr := ethcmn.BytesToAddress(msg.Recipient.Bytes())
	return &addr
}

// sigCache is used to cache the derived sender and contains the signer used
// to derive it.
type sigCache struct {
	signer ethtypes.Signer
	from   ethcmn.Address
}

// NewMsgEthereumTx returns a reference to a new Ethereum transaction message.
func NewMsgEthereumTx(
	nonce uint64, to *ethcmn.Address, amount *big.Int,
	gasLimit uint64, gasPrice *big.Int, payload []byte,
) MsgEthereumTx {
	return newMsgEthereumTx(nonce, to, amount, gasLimit, gasPrice, payload)
}

// NewMsgEthereumTxContract returns a reference to a new Ethereum transaction
// message designated for contract creation.
func NewMsgEthereumTxContract(
	nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, payload []byte,
) MsgEthereumTx {
	return newMsgEthereumTx(nonce, nil, amount, gasLimit, gasPrice, payload)
}

func newMsgEthereumTx(
	nonce uint64, to *ethcmn.Address, amount *big.Int,
	gasLimit uint64, gasPrice *big.Int, payload []byte,
) MsgEthereumTx {
	if len(payload) > 0 {
		payload = ethcmn.CopyBytes(payload)
	}

	txData := TxData{
		AccountNonce: nonce,
		Recipient:    to.Bytes(),
		Payload:      payload,
		GasLimit:     gasLimit,
		Amount:       []byte{},
		Price:        []byte{},
		V:            []byte{},
		R:            []byte{},
		S:            []byte{},
	}

	if amount != nil {
		txData.Amount = new(big.Int).Set(amount).Bytes()
	}
	if gasPrice != nil {
		txData.Price = new(big.Int).Set(gasPrice).Bytes()
	}

	return MsgEthereumTx{Data: txData}
}

// Route returns the route value of an MsgEthereumTx.
func (msg MsgEthereumTx) Route() string { return RouterKey }

// Type returns the type value of an MsgEthereumTx.
func (msg MsgEthereumTx) Type() string { return TypeMsgEthereumTx }

// ValidateBasic implements the sdk.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg MsgEthereumTx) ValidateBasic() error {
	gasPrice := new(big.Int).SetBytes(msg.Data.Price)
	if gasPrice.Sign() != 1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "price must be positive %s", gasPrice)
	}

	// Amount can be 0
	amount := new(big.Int).SetBytes(msg.Data.Amount)
	if amount.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", amount)
	}

	return nil
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg MsgEthereumTx) To() *ethcmn.Address {
	recipient := ethcmn.BytesToAddress(msg.Data.Recipient)
	return &recipient
}

// GetMsgs returns a single MsgEthereumTx as an sdk.Msg.
func (msg MsgEthereumTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSigners returns the expected signers for an Ethereum transaction message.
// For such a message, there should exist only a single 'signer'.
//
// NOTE: This method panics if 'VerifySig' hasn't been called first.
func (msg MsgEthereumTx) GetSigners() []sdk.AccAddress {
	sender := msg.From()
	if sender.Empty() {
		panic("must use 'VerifySig' with a chain ID to get the signer")
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes returns the Amino bytes of an Ethereum transaction message used
// for signing.
//
// NOTE: This method cannot be used as a chain ID is needed to create valid bytes
// to sign over. Use 'RLPSignBytes' instead.
func (msg MsgEthereumTx) GetSignBytes() []byte {
	panic("must use 'RLPSignBytes' with a chain ID to get the valid bytes to sign")
}

// RLPSignBytes returns the RLP hash of an Ethereum transaction message with a
// given chainID used for signing.
func (msg MsgEthereumTx) RLPSignBytes(chainID *big.Int) ethcmn.Hash {
	return rlpHash([]interface{}{
		msg.Data.AccountNonce,
		new(big.Int).SetBytes(msg.Data.Price),
		msg.Data.GasLimit,
		ethcmn.BytesToAddress(msg.Data.Recipient),
		new(big.Int).SetBytes(msg.Data.Amount),
		msg.Data.Payload,
		chainID,
		uint(0),
		uint(0),
	})
}

// EncodeRLP implements the rlp.Encoder interface.
func (msg *MsgEthereumTx) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &msg.Data)
}

// DecodeRLP implements the rlp.Decoder interface.
func (msg *MsgEthereumTx) DecodeRLP(s *rlp.Stream) error {
	_, size, err := s.Kind()
	if err != nil {
		// return error if stream is too large
		return err
	}

	if err := s.Decode(&msg.Data); err != nil {
		return err
	}

	msg.size.Store(ethcmn.StorageSize(rlp.ListSize(size)))
	return nil
}

// Sign calculates a secp256k1 ECDSA signature and signs the transaction. It
// takes a private key and chainID to sign an Ethereum transaction according to
// EIP155 standard. It mutates the transaction as it populates the V, R, S
// fields of the Transaction's Signature.
func (msg *MsgEthereumTx) Sign(chainID *big.Int, priv *ecdsa.PrivateKey) error {
	txHash := msg.RLPSignBytes(chainID)

	sig, err := ethcrypto.Sign(txHash[:], priv)
	if err != nil {
		return err
	}

	if len(sig) != 65 {
		return fmt.Errorf("wrong size for signature: got %d, want 65", len(sig))
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])

	var v *big.Int

	if chainID.Sign() == 0 {
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	} else {
		v = big.NewInt(int64(sig[64] + 35))
		chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))

		v.Add(v, chainIDMul)
	}

	msg.Data.V = v.Bytes()
	msg.Data.R = r.Bytes()
	msg.Data.S = s.Bytes()

	return nil
}

// VerifySig attempts to verify a Transaction's signature for a given chainID.
// A derived address is returned upon success or an error if recovery fails.
func (msg *MsgEthereumTx) VerifySig(chainID *big.Int) (ethcmn.Address, error) {
	r := new(big.Int).SetBytes(msg.Data.R)
	s := new(big.Int).SetBytes(msg.Data.S)
	v := new(big.Int).SetBytes(msg.Data.V)

	signer := ethtypes.NewEIP155Signer(chainID)

	if msg.from != nil {
		// If the signer used to derive from in a previous call is not the same as
		// used current, invalidate the cache.
		// TODO: signer bytes -> Signer
		if signer.Equal(msg.from.Getsigner()) {
			return ethcmn.BytesToAddress(msg.from.Getfrom()), nil
		}
	}

	// do not allow recovery for transactions with an unprotected chainID
	if chainID.Sign() == 0 {
		return ethcmn.Address{}, errors.New("chainID cannot be zero")
	}

	chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))

	V := new(big.Int).Sub(v, chainIDMul)
	V.Sub(V, big8)

	sigHash := msg.RLPSignBytes(chainID)
	sender, err := recoverEthSig(r, s, V, sigHash)
	if err != nil {
		return ethcmn.Address{}, err
	}

	msg.from = &SigCache{signer: signer, from: sender.Bytes()}
	return sender, nil
}

// GetGas implements the GasTx interface. It returns the GasLimit of the transaction.
func (msg MsgEthereumTx) GetGas() uint64 {
	return msg.Data.GasLimit
}

// Fee returns gasprice * gaslimit.
func (msg MsgEthereumTx) Fee() *big.Int {
	gasPrice := new(big.Int).SetBytes(msg.Data.Price)
	gasLimit := new(big.Int).SetUint64(msg.Data.GasLimit)
	return new(big.Int).Mul(gasPrice, gasLimit)
}

// ChainID returns which chain id this transaction was signed for (if at all)
func (msg *MsgEthereumTx) ChainID() *big.Int {
	return deriveChainID(new(big.Int).SetBytes(msg.Data.V))
}

// Cost returns amount + gasprice * gaslimit.
func (msg MsgEthereumTx) Cost() *big.Int {
	total := msg.Fee()
	total.Add(total, new(big.Int).SetBytes(msg.Data.Amount))
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (msg MsgEthereumTx) RawSignatureValues() (v, r, s *big.Int) {
	return new(big.Int).SetBytes(msg.Data.V), new(big.Int).SetBytes(msg.Data.R), new(big.Int).SetBytes(msg.Data.S)
}

// From loads the ethereum sender address from the sigcache and returns an
// sdk.AccAddress from its bytes
func (msg *MsgEthereumTx) From() sdk.AccAddress {
	sc := msg.from.Load()
	if sc == nil {
		return nil
	}

	sigCache := sc.(sigCache)

	if len(sigCache.from.Bytes()) == 0 {
		return nil
	}

	return sdk.AccAddress(sigCache.from.Bytes())
}

// deriveChainID derives the chain id from the given v parameter
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
