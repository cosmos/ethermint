package rpc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	params "github.com/cosmos/ethermint/rpc/args"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/crypto/hd"
	params "github.com/cosmos/ethermint/rpc/args"
)

// PersonalEthAPI is the personal_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PersonalEthAPI struct {
	ethAPI   *PublicEthAPI
	keyInfos []keyring.Info // all keys, both locked and unlocked. unlocked keys are stored in ethAPI.keys
}

// NewPersonalEthAPI creates an instance of the public Personal Eth API.
func NewPersonalEthAPI(ethAPI *PublicEthAPI) *PersonalEthAPI {
	api := &PersonalEthAPI{
		ethAPI: ethAPI,
	}

	infos, err := api.getKeybaseInfo()
	if err != nil {
		return api
	}

	api.keyInfos = infos
	return api
}

func (e *PersonalEthAPI) getKeybaseInfo() ([]keyring.Info, error) {
	e.ethAPI.keybaseLock.Lock()
	defer e.ethAPI.keybaseLock.Unlock()

	if e.ethAPI.clientCtx.Keybase == nil {
		keybase, err := keyring.New(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			viper.GetString(flags.FlagHome),
			e.ethAPI.clientCtx.Input,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return nil, err
		}

		e.ethAPI.clientCtx.Keybase = keybase
	}

	return e.ethAPI.clientCtx.Keybase.List()
}

// ImportRawKey armors and encrypts a given raw hex encoded ECDSA key and stores it into the key directory.
// The name of the key will have the format "personal_<length-keys>", where <length-keys> is the total number of
// keys stored on the keyring.
// NOTE: The key will be both armored and encrypted using the same passphrase.
func (e *PersonalEthAPI) ImportRawKey(privkey, password string) (common.Address, error) {
	e.ethAPI.logger.Debug("personal_importRawKey")
	priv, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}

	privKey := &ethsecp256k1.PrivKey{Key: crypto.FromECDSA(priv)}

	armor := sdkcrypto.EncryptArmorPrivKey(privKey, password, ethsecp256k1.KeyType)

	// ignore error as we only care about the length of the list
	list, _ := e.ethAPI.clientCtx.Keybase.List()
	privKeyName := fmt.Sprintf("personal_%d", len(list))

	if err := e.ethAPI.clientCtx.Keybase.ImportPrivKey(privKeyName, armor, password); err != nil {
		return common.Address{}, err
	}

	addr := common.BytesToAddress(privKey.PubKey().Address().Bytes())

	info, err := e.ethAPI.clientCtx.Keybase.Get(privKeyName)
	if err != nil {
		return common.Address{}, err
	}

	// append key and info to be able to lock and list the account
	//e.ethAPI.keys = append(e.ethAPI.keys, privKey)
	e.keyInfos = append(e.keyInfos, info)
	e.ethAPI.logger.Info("key successfully imported", "name", privKeyName, "address", addr.String())

	return addr, nil
}

// ListAccounts will return a list of addresses for accounts this node manages.
func (e *PersonalEthAPI) ListAccounts() ([]common.Address, error) {
	e.ethAPI.logger.Debug("personal_listAccounts")
	addrs := []common.Address{}
	for _, info := range e.keyInfos {
		addressBytes := info.GetPubKey().Address().Bytes()
		addrs = append(addrs, common.BytesToAddress(addressBytes))
	}

	return addrs, nil
}

// LockAccount will lock the account associated with the given address when it's unlocked.
// It removes the key corresponding to the given address from the API's local keys.
func (e *PersonalEthAPI) LockAccount(address common.Address) bool {
	e.ethAPI.logger.Debug("personal_lockAccount", "address", address.String())

	for i, key := range e.ethAPI.keys {
		if !bytes.Equal(key.PubKey().Address().Bytes(), address.Bytes()) {
			continue
		}

		tmp := make([]ethsecp256k1.PrivKey, len(e.ethAPI.keys)-1)
		copy(tmp[:i], e.ethAPI.keys[:i])
		copy(tmp[i:], e.ethAPI.keys[i+1:])
		e.ethAPI.keys = tmp

		e.ethAPI.logger.Debug("account unlocked", "address", address.String())
		return true
	}

	return false
}

// NewAccount will create a new account and returns the address for the new account.
func (e *PersonalEthAPI) NewAccount(password string) (common.Address, error) {
	e.ethAPI.logger.Debug("personal_newAccount")
	_, err := e.getKeybaseInfo()
	if err != nil {
		return common.Address{}, err
	}

	name := "key_" + time.Now().UTC().Format(time.RFC3339)
	info, _, err := e.ethAPI.clientCtx.Keybase.CreateMnemonic(name, keyring.English, password, hd.EthSecp256k1)
	if err != nil {
		return common.Address{}, err
	}

	e.keyInfos = append(e.keyInfos, info)

	addr := common.BytesToAddress(info.GetPubKey().Address().Bytes())
	e.ethAPI.logger.Info("Your new key was generated", "address", addr.String())
	e.ethAPI.logger.Info("Please backup your key file!", "path", os.Getenv("HOME")+"/.ethermintd/"+name)
	e.ethAPI.logger.Info("Please remember your password!")
	return addr, nil
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
// It exports the private key corresponding to the given address from the keyring and stores it in the API's local keys.
func (e *PersonalEthAPI) UnlockAccount(_ context.Context, addr common.Address, password string, _ *uint64) (bool, error) { // nolint: interfacer
	e.ethAPI.logger.Debug("personal_unlockAccount", "address", addr.String())
	// TODO: use duration

	var keyInfo keyring.Info

	for _, info := range e.keyInfos {
		addressBytes := info.GetPubKey().Address().Bytes()
		if bytes.Equal(addressBytes, addr[:]) {
			keyInfo = info
			break
		}
	}

	if keyInfo == nil {
		return false, fmt.Errorf("cannot find key with given address %s", addr.String())
	}

	// exporting private key only works on local keys
	if keyInfo.GetType() != keyring.TypeLocal {
		return false, fmt.Errorf("key type must be %s, got %s", keyring.TypeLedger.String(), keyInfo.GetType().String())
	}

	armor, err := e.ethAPI.clientCtx.Keybase.ExportPrivKeyArmor(keyInfo.GetName(), password)
	if err != nil {
		return err
	}

	privKey, algo, err := sdkcrypto.UnarmorDecryptPrivKey(armor, password)
	if err != nil {
		return err
	}

	if algo != ethsecp256k1.KeyType {
		return fmt.Errorf("invalid key algorithm, got %s, expected %s", algo, ethsecp256k1.KeyType)
	}

	ethermintPrivKey, ok := privKey.(*ethsecp256k1.PrivKey)
	if !ok {
		return fmt.Errorf("invalid private key type %T, expected %T", privKey, &ethsecp256k1.PrivKey{})
	}

	e.ethAPI.keys = append(e.ethAPI.keys, ethermintPrivKey)
	e.ethAPI.logger.Debug("account unlocked", "address", addr.String())
	return true, nil
}

// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.To. If the given password isn't
// able to decrypt the key it fails.
func (e *PersonalEthAPI) SendTransaction(_ context.Context, args params.SendTxArgs, _ string) (common.Hash, error) {
	return e.ethAPI.SendTransaction(args)
}

// Sign calculates an Ethereum ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func (e *PersonalEthAPI) Sign(_ context.Context, data hexutil.Bytes, addr common.Address, _ string) (hexutil.Bytes, error) {
	e.ethAPI.logger.Debug("personal_sign", "data", data, "address", addr.String())

	key, ok := checkKeyInKeyring(e.ethAPI.keys, addr)
	if !ok {
		return nil, fmt.Errorf("cannot find key with address %s", addr.String())
	}

	sig, err := crypto.Sign(accounts.TextHash(data), key.ToECDSA())
	if err != nil {
		return nil, err
	}

	sig[crypto.RecoveryIDOffset] += 27 // transform V from 0/1 to 27/28
	return sig, nil
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
func (e *PersonalEthAPI) EcRecover(_ context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	e.ethAPI.logger.Debug("personal_ecRecover", "data", data, "sig", sig)

	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	pubkey, err := crypto.SigToPub(accounts.TextHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}
