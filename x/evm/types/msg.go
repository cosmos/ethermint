package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/cosmos/ethermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	_ sdk.Msg = MsgEthereumTx{}
	_ sdk.Tx  = MsgEthereumTx{}
)

var big8 = big.NewInt(8)

// TypeMsgEthereumTx defines the type string of an Ethereum transaction
const TypeMsgEthereumTx = "ethereum"

// MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
type MsgEthereumTx struct {
	Data TxData

	// caches
	size atomic.Value
	from atomic.Value
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
		Recipient:    to,
		Payload:      payload,
		GasLimit:     gasLimit,
		Amount:       new(big.Int),
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}

	if amount != nil {
		txData.Amount.Set(amount)
	}
	if gasPrice != nil {
		txData.Price.Set(gasPrice)
	}

	return MsgEthereumTx{Data: txData}
}

// String implements the fmt.Stringer interface
func (msg MsgEthereumTx) String() string {
	return msg.Data.String()
}

// Route returns the route value of an MsgEthereumTx.
func (msg MsgEthereumTx) Route() string { return RouterKey }

// Type returns the type value of an MsgEthereumTx.
func (msg MsgEthereumTx) Type() string { return TypeMsgEthereumTx }

// ValidateBasic implements the sdk.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg MsgEthereumTx) ValidateBasic() error {
	if msg.Data.Price.Cmp(big.NewInt(0)) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be 0")
	}

	if msg.Data.Price.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "gas price cannot be negative %s", msg.Data.Price)
	}

	// Amount can be 0
	if msg.Data.Amount.Sign() == -1 {
		return sdkerrors.Wrapf(types.ErrInvalidValue, "amount cannot be negative %s", msg.Data.Amount)
	}

	return nil
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg MsgEthereumTx) To() *ethcmn.Address {
	return msg.Data.Recipient
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
		msg.Data.Price,
		msg.Data.GasLimit,
		msg.Data.Recipient,
		msg.Data.Amount,
		msg.Data.Payload,
		chainID, uint(0), uint(0),
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

	msg.Data.V = v
	msg.Data.R = r
	msg.Data.S = s
	return nil
}

// VerifySig attempts to verify a Transaction's signature for a given chainID.
// A derived address is returned upon success or an error if recovery fails.
// The returned address will be cached in the message from field if the verification is successful.
func (msg *MsgEthereumTx) VerifySig(chainID *big.Int) (ethcmn.Address, error) {
	signer := ethtypes.NewEIP155Signer(chainID)

	if sc := msg.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous call is not the same as
		// used current, invalidate the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	// do not allow recovery for transactions with an unprotected chainID
	if chainID.Sign() == 0 {
		return ethcmn.Address{}, errors.New("chainID cannot be zero")
	}

	chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))
	V := new(big.Int).Sub(msg.Data.V, chainIDMul)
	V.Sub(V, big8)

	sigHash := msg.RLPSignBytes(chainID)
	sender, err := recoverEthSig(msg.Data.R, msg.Data.S, V, sigHash)
	if err != nil {
		return ethcmn.Address{}, err
	}

	msg.from.Store(sigCache{signer: signer, from: sender})
	return sender, nil
}

// GetGas implements the GasTx interface. It returns the GasLimit of the transaction.
func (msg MsgEthereumTx) GetGas() uint64 {
	return msg.Data.GasLimit
}

// Fee returns gasprice * gaslimit.
func (msg MsgEthereumTx) Fee() *big.Int {
	return new(big.Int).Mul(msg.Data.Price, new(big.Int).SetUint64(msg.Data.GasLimit))
}

// ChainID returns which chain id this transaction was signed for (if at all)
func (msg *MsgEthereumTx) ChainID() *big.Int {
	return deriveChainID(msg.Data.V)
}

// Cost returns amount + gasprice * gaslimit.
func (msg MsgEthereumTx) Cost() *big.Int {
	total := msg.Fee()
	total.Add(total, msg.Data.Amount)
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (msg MsgEthereumTx) RawSignatureValues() (v, r, s *big.Int) {
	return msg.Data.V, msg.Data.R, msg.Data.S
}

// From loads the ethereum sender address from the signature cache and returns an
// sdk.AccAddress from its bytes. The address needs to be cached previously using
// the VerifySig function.
//
// NOTE: if the sender address hasn't been cached, this function will return an
// nil AccAddress.
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
