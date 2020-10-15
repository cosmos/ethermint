package util

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type NonceCache interface {
	Serialize(account common.Address, fn func() error) error
	Sync(account common.Address, syncFn func() (uint64, error))

	Set(account common.Address, nonce int64)
	Get(account common.Address) int64
	Incr(account common.Address) int64
	Decr(account common.Address) int64
}

func NewNonceCache() NonceCache {
	return &nonceCache{
		mux:    new(sync.RWMutex),
		nonces: make(map[common.Address]int64),
		locks:  make(map[common.Address]*sync.RWMutex),
		guard:  NewUniquify(),
	}
}

type nonceCache struct {
	mux    *sync.RWMutex
	nonces map[common.Address]int64
	locks  map[common.Address]*sync.RWMutex
	guard  Uniquify
}

// Serialize serializes access to the nonce cache for all goroutines, all nonce increments should be done
// in this context. If a transaction increments nonce, but has not been submitted,
// it will have exclusive right to decrease nonce back for other transactions.
func (n nonceCache) Serialize(account common.Address, fn func() error) error {
	return n.guard.Call(account.Hex(), fn)
}

func (n nonceCache) Get(account common.Address) int64 {
	n.mux.RLock()
	lock, ok := n.locks[account]
	if !ok {
		n.mux.RUnlock()
		return 0
	}
	lock.RLock()
	n.mux.RUnlock()
	nonce := n.nonces[account]
	lock.RUnlock()
	return nonce
}

func (n nonceCache) Set(account common.Address, nonce int64) {
	n.mux.Lock()
	lock, ok := n.locks[account]
	if !ok {
		n.locks[account] = new(sync.RWMutex)
		n.nonces[account] = nonce
		n.mux.Unlock()
		return
	}
	lock.Lock()
	n.mux.Unlock()
	n.nonces[account] = nonce
	lock.Unlock()
}

func (n nonceCache) Incr(account common.Address) int64 {
	n.mux.Lock()
	lock, ok := n.locks[account]
	if !ok {
		n.nonces[account] = 1
		n.locks[account] = new(sync.RWMutex)
		n.mux.Unlock()
		return 0
	}
	lock.Lock()
	n.mux.Unlock()
	nonce := n.nonces[account]
	n.nonces[account]++
	lock.Unlock()
	return nonce
}

func (n nonceCache) Decr(account common.Address) int64 {
	n.mux.Lock()
	lock, ok := n.locks[account]
	if !ok {
		n.nonces[account] = 0
		n.locks[account] = new(sync.RWMutex)
		n.mux.Unlock()
		return 0
	}
	lock.Lock()
	n.mux.Unlock()
	nonce := n.nonces[account]
	n.nonces[account]--
	lock.Unlock()
	return nonce
}

func (n nonceCache) Sync(account common.Address, syncFn func() (uint64, error)) {
	n.mux.Lock()
	lock, ok := n.locks[account]
	if !ok {
		n.nonces[account] = 0
		n.locks[account] = new(sync.RWMutex)
		lock = n.locks[account]
	}
	lock.Lock()
	n.mux.Unlock()
	prevNonce, prevOk := n.nonces[account]
	{
		n.mux.RLock()
		nextNonce, nextOk := n.nonces[account]
		n.mux.RUnlock()
		if !prevOk && nextOk {
			// we're not fist here to sync - skip
			lock.Unlock()
			return
		} else if nextNonce != prevNonce {
			lock.Unlock()
			return
		}
		if nonce, err := syncFn(); err == nil {
			n.mux.Lock()
			n.nonces[account] = int64(nonce)
			n.mux.Unlock()
		}
	}
	lock.Unlock()
}
