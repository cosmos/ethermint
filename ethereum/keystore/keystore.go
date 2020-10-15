package keystore

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	log "github.com/xlab/suplog"

	eth "github.com/cosmos/ethermint/ethereum/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type EthKeyStore interface {
	PrivateKey(account common.Address, password string) (key *ecdsa.PrivateKey, ok bool)
	SignerFn(account common.Address, password string) bind.SignerFn
	UnsetKey(account common.Address, password string)
	Accounts() []common.Address
	AddPath(keybase string) error
	RemovePath(keybase string)
	Paths() []string
}

func New(paths ...string) (EthKeyStore, error) {
	ks := &keyStore{
		cache:                      eth.NewKeyCache(),
		notifyWalletSubscribersMux: new(sync.RWMutex),
		paths:                      make(map[string]struct{}),
		pathsMux:                   new(sync.RWMutex),
	}
	for _, path := range paths {
		ks.paths[path] = struct{}{}
	}
	ks.checkPaths()
	return ks, nil
}

type keyStore struct {
	cache                      eth.KeyCache
	notifyWalletSubscribers    []chan<- *WalletSpec
	notifyWalletSubscribersMux *sync.RWMutex

	paths    map[string]struct{}
	pathsMux *sync.RWMutex
}

func (ks *keyStore) checkPaths() {
	paths := ks.Paths()
	for _, keybasePath := range paths {
		err := ks.forEachWallet(keybasePath, func(spec *WalletSpec) error {
			if isNew := ks.cache.SetPath(spec.HexToAddress(), spec.Path); isNew {
				subs := ks.getNotifyWalletSubscribers()
				for _, notifyC := range subs {
					select {
					case notifyC <- spec:
					default:
					}
				}
			}
			return nil
		})
		if err != nil {
			log.WithFields(log.Fields{
				"keybasePath": keybasePath,
				"fn":          "checkPaths",
			}).WithError(err).Warningln("failed to lookup")
		}
	}
}

func (ks *keyStore) PrivateKey(account common.Address, password string) (key *ecdsa.PrivateKey, ok bool) {
	return ks.cache.PrivateKey(account, password)
}

func (ks *keyStore) SignerFn(account common.Address, password string) bind.SignerFn {
	return ks.cache.SignerFn(account, password)
}

func (ks *keyStore) UnsetKey(account common.Address, password string) {
	ks.cache.UnsetKey(account, password)
}

func (ks *keyStore) Accounts() []common.Address {
	paths := ks.Paths()
	var accounts []common.Address
	for _, keybasePath := range paths {
		if err := ks.forEachWallet(keybasePath, func(spec *WalletSpec) error {
			accounts = append(accounts, spec.HexToAddress())
			return nil
		}); err != nil {
			log.WithFields(log.Fields{
				"keybasePath": keybasePath,
				"fn":          "Accounts",
			}).WithError(err).Warningln("failed to lookup")
		}
	}
	return accounts
}

var errRangeStop = errors.New("stop")

func (ks *keyStore) forEachWallet(keybasePath string, fn func(spec *WalletSpec) error) error {
	return filepath.Walk(keybasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if path == keybasePath {
			return nil
		} else if info.IsDir() {
			return filepath.SkipDir
		}
		var spec *WalletSpec
		if data, err := ioutil.ReadFile(path); err != nil {
			return err
		} else if err = json.Unmarshal(data, &spec); err != nil {
			return err
		}
		if len(spec.Address) == 0 {
			return fmt.Errorf("failed to load address from %s", path)
		} else if !common.IsHexAddress(spec.Address) {
			return fmt.Errorf("wrong (not hex) address from %s", path)
		}
		spec.Path = path
		return fn(spec)
	})
}

func (ks *keyStore) AddPath(keybase string) error {
	f, err := os.Stat(keybase)
	if err != nil {
		return err
	} else if !f.IsDir() {
		return fmt.Errorf("%s is not a directory", keybase)
	}
	ks.pathsMux.Lock()
	ks.paths[keybase] = struct{}{}
	ks.pathsMux.Unlock()
	return nil
}

func (ks *keyStore) RemovePath(keybase string) {
	ks.pathsMux.Lock()
	delete(ks.paths, keybase)
	ks.pathsMux.Unlock()
}

func (ks *keyStore) Paths() []string {
	ks.pathsMux.RLock()
	paths := make([]string, 0, len(ks.paths))
	for p := range ks.paths {
		paths = append(paths, p)
	}
	ks.pathsMux.RUnlock()
	sort.Sort(sort.StringSlice(paths))
	return paths
}

func (ks *keyStore) NewWalletSubscribeNotify(notifyC chan<- *WalletSpec) {
	ks.notifyWalletSubscribersMux.Lock()
	ks.notifyWalletSubscribers = append(ks.notifyWalletSubscribers, notifyC)
	ks.notifyWalletSubscribersMux.Unlock()
}

func (ks *keyStore) getNotifyWalletSubscribers() []chan<- *WalletSpec {
	ks.notifyWalletSubscribersMux.RLock()
	subs := ks.notifyWalletSubscribers
	ks.notifyWalletSubscribersMux.RUnlock()
	return subs
}

type WalletSpec struct {
	Address string `json:"address"`
	ID      string `json:"id"`
	Version int    `json:"version"`
	Path    string `json:"-"`
}

func (spec *WalletSpec) HexToAddress() common.Address {
	return common.HexToAddress(spec.Address)
}
