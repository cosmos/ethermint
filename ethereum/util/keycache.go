package util

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"errors"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyCache interface {
	SetPath(account common.Address, path string) bool
	UnsetPath(account common.Address, path string)
	PrivateKey(account common.Address, password string) (key *ecdsa.PrivateKey, ok bool)
	SetPrivateKey(account common.Address, pk *ecdsa.PrivateKey)
	UnsetKey(account common.Address, password string)
	SignerFn(account common.Address, password string) bind.SignerFn
}

func NewKeyCache() KeyCache {
	return &keyCache{
		paths:    make(map[common.Address]string),
		pathsMux: new(sync.RWMutex),
		keys:     make(map[string]*ecdsa.PrivateKey),
		keysMux:  new(sync.RWMutex),
		guard:    NewUniquify(),
	}
}

type keyCache struct {
	paths    map[common.Address]string
	pathsMux *sync.RWMutex
	keys     map[string]*ecdsa.PrivateKey
	keysMux  *sync.RWMutex
	guard    Uniquify
}

// SetPath sets the wallet path for a given account. Returns true if the new path
// has been added or was changed.
func (k *keyCache) SetPath(account common.Address, path string) bool {
	k.pathsMux.Lock()
	prevPath, existing := k.paths[account]
	k.paths[account] = path
	k.pathsMux.Unlock()
	return !existing || prevPath != path
}

func (k *keyCache) UnsetPath(account common.Address, path string) {
	k.pathsMux.Lock()
	delete(k.paths, account)
	k.pathsMux.Unlock()
}

var (
	ErrNoKeyStore = errors.New("no keystore or file for account")
	ErrKeyDecrypt = errors.New("private key decryption failed")
)

func (k *keyCache) UnsetKey(account common.Address, password string) {
	h := hashAccountPass(account, password)
	k.keysMux.Lock()
	delete(k.keys, string(h))
	k.keysMux.Unlock()
}

func (k *keyCache) SetPrivateKey(account common.Address, pk *ecdsa.PrivateKey) {
	h := hashAccountPass(account, "")
	k.keysMux.Lock()
	k.keys[string(h)] = pk
	k.keysMux.Unlock()
}

func (k *keyCache) PrivateKey(account common.Address, password string) (key *ecdsa.PrivateKey, ok bool) {
	h := hashAccountPass(account, password)
	if err := k.guard.Call(string(h), func() error {
		k.keysMux.RLock()
		key, ok = k.keys[string(h)]
		k.keysMux.RUnlock()
		if ok {
			return nil
		}
		k.pathsMux.RLock()
		path, pathOk := k.paths[account]
		k.pathsMux.RUnlock()
		if !pathOk {
			return ErrNoKeyStore
		}
		if strings.HasPrefix(path, "keystore://") {
			path = strings.TrimPrefix(path, "keystore://")
		}
		keyJSON, err := ioutil.ReadFile(path)
		if err != nil {
			return ErrNoKeyStore
		}
		pk, err := keystore.DecryptKey(keyJSON, password)
		if err != nil {
			return ErrKeyDecrypt
		}
		k.keysMux.Lock()
		k.keys[string(h)] = pk.PrivateKey
		k.keysMux.Unlock()
		key = pk.PrivateKey
		ok = true
		return nil
	}); err != nil {
		return nil, false
	}
	return key, ok
}

func (k *keyCache) SignerFn(account common.Address, password string) bind.SignerFn {
	key, ok := k.PrivateKey(account, password)
	if !ok {
		return nil
	}
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != keyAddr {
			return nil, errors.New("not authorized to sign this account")
		}
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
}

func SignerFnForPk(privKey *ecdsa.PrivateKey) bind.SignerFn {
	keyAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	return func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != keyAddr {
			return nil, errors.New("not authorized to sign this account")
		}
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privKey)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
}

var hashSep = []byte("-")

func hashAccountPass(account common.Address, password string) []byte {
	h := sha1.New()
	h.Write(account[:])
	h.Write(hashSep)
	h.Write([]byte(password))
	return h.Sum(nil)
}
