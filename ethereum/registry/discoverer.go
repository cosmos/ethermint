package registry

import (
	"context"
	"sync"
	"time"

	log "github.com/xlab/suplog"

	"github.com/ethereum/go-ethereum/common"
)

type ContractsSet struct {
	CoordinatorContract common.Address
	DevUtilsContract    common.Address
	ExchangeContract    common.Address
	FuturesContract     common.Address
}

type ContractDiscoverer interface {
	// GetContracts returns only when all contracts in the set are not empty.
	GetContracts() ContractsSet
}

type contractDiscoverer struct {
	finalSet    *ContractsSet
	finalSetMux *sync.RWMutex

	// lazyProvider func() provider.EVMProvider
}

// // NewContractDiscoverer inits only from provider func that returns when provider is ready.
// func NewContractDiscoverer(lazyProvider func() provider.EVMProvider) ContractDiscoverer {
// 	return &contractDiscoverer{
// 		lazyProvider: lazyProvider,
// 		finalSetMux:  new(sync.RWMutex),
// 	}
// }

func (c *contractDiscoverer) GetContracts() ContractsSet {
	c.finalSetMux.RLock()
	if c.finalSet != nil {
		return *c.finalSet
	}
	c.finalSetMux.RUnlock()

	// ethProvider := c.lazyProvider()

	var set ContractsSet
	for {
		ts := time.Now()
		_, cancelFn := context.WithTimeout(context.Background(), defaultRPCTimeout)
		// set.ExchangeContract = discoverContractAddress(discoverCtx, "@0x/exchange", ethProvider)
		// set.DevUtilsContract = discoverContractAddress(discoverCtx, "@0x/devutils", ethProvider)
		// set.CoordinatorContract = discoverContractAddress(discoverCtx, "@injective/coordinator", ethProvider)
		// set.FuturesContract = discoverContractAddress(discoverCtx, "@injective/futures", ethProvider)
		cancelFn()
		log.Infoln("Contract addresses discovered in", time.Since(ts))

		if hasEmptyAddresses(
			set.ExchangeContract,
			set.DevUtilsContract,
			set.CoordinatorContract,
			set.FuturesContract,
		) {
			log.WithFields(log.Fields{
				"exchangeContract":    set.ExchangeContract.Hex(),
				"devUtilsContract":    set.DevUtilsContract.Hex(),
				"coordinatorContract": set.CoordinatorContract.Hex(),
				"futuresContract":     set.FuturesContract.Hex(),
			}).Println("Still have some addresses empty")

			time.Sleep(10 * time.Second)

			continue
		}

		c.finalSetMux.Lock()
		c.finalSet = &set
		c.finalSetMux.Unlock()

		return set
	}
}

const defaultRPCTimeout = 10 * time.Second

var (
	defaultRegistryAddress = common.HexToAddress("0x5C7e1fc74fe17242a077DB7DFd962e897Ed4e39a")
	defaultReadonlyAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func hasEmptyAddresses(addresses ...common.Address) bool {
	for _, addr := range addresses {
		if (addr == common.Address{}) {
			return true
		}
	}

	return false
}

// func discoverContractAddress(ctx context.Context, name string, provider provider.EVMProvider) common.Address {
// 	caller, _ := contracts.NewRegistryCaller(defaultRegistryAddress, provider)

// 	opts := &bind.CallOpts{
// 		From:    defaultReadonlyAddress,
// 		Context: ctx,
// 	}

// 	address, err := caller.GetContractAddress(opts, name)
// 	if err != nil {
// 		return common.Address{}
// 	}

// 	return address
// }
