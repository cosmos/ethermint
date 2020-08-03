package rpc

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	sdkcontext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	emintcrypto "github.com/cosmos/ethermint/crypto"
	params "github.com/cosmos/ethermint/rpc/args"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PersonalEthAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PersonalEthAPI struct {
	cliCtx    sdkcontext.CLIContext
	nonceLock *AddrLocker
	keys        []emintcrypto.PrivKeySecp256k1
	keybaseLock sync.Mutex
}

// NewPersonalEthAPI creates an instance of the public ETH Web3 API.
func NewPersonalEthAPI(cliCtx sdkcontext.CLIContext, nonceLock *AddrLocker, keys []emintcrypto.PrivKeySecp256k1) *PersonalEthAPI {
	return &PersonalEthAPI{
		cliCtx:    cliCtx,
		nonceLock: nonceLock,
		keys: keys,
	}
}

// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (e *PersonalEthAPI) ImportRawKey(privkey, password string) (common.Address, error) {
	return common.Address{}, nil
}

// ListAccounts will return a list of addresses for accounts this node manages.
func (e *PersonalEthAPI) ListAccounts() ([]common.Address, error) {
	e.keybaseLock.Lock()
	addrs := []common.Address{}

	if e.cliCtx.Keybase == nil {
		keybase, err := keyring.NewKeyring(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			viper.GetString(flags.FlagHome),
			e.cliCtx.Input,
		)
		if err != nil {
			return addrs, err
		}

		e.cliCtx.Keybase = keybase
	}

	infos, err := e.cliCtx.Keybase.List()
	if err != nil {
		return addrs, err
	}

	e.keybaseLock.Unlock()

	for _, info := range infos {
		addressBytes := info.GetPubKey().Address().Bytes()
		addrs = append(addrs, common.BytesToAddress(addressBytes))
	}

	return addrs, nil
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (e *PersonalEthAPI) LockAccount(addr common.Address) bool {
	return false
}

// NewAccount will create a new account and returns the address for the new account.
func (e *PersonalEthAPI) NewAccount(password string) (common.Address, error) {
	return common.Address{}, nil
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (e *PersonalEthAPI) UnlockAccount(ctx context.Context, addr common.Address, password string, duration *uint64) (bool, error) {
	return false, nil
}

// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.To. If the given passwd isn't
// able to decrypt the key it fails.
func (e *PersonalEthAPI) SendTransaction(ctx context.Context, args params.SendTxArgs, passwd string) (common.Hash, error) {
	return common.Hash{}, nil
}

// Sign calculates an Ethereum ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func (e *PersonalEthAPI) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	return nil, nil
}

// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be 27 or 28 for legacy reasons.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecove
func (e *PersonalEthAPI) EcRecover(ctx context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	return common.Address{}, nil
}
