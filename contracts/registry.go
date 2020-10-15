// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lockId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"contractName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"}],\"name\":\"ContractAddressConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lockId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"contractName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"}],\"name\":\"ContractAddressLocked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractName\",\"type\":\"string\"}],\"name\":\"getContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// GetContractAddress is a free data retrieval call binding the contract method 0x04433bbc.
//
// Solidity: function getContractAddress(string contractName) view returns(address)
func (_Registry *RegistryCaller) GetContractAddress(opts *bind.CallOpts, contractName string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "getContractAddress", contractName)
	return *ret0, err
}

// GetContractAddress is a free data retrieval call binding the contract method 0x04433bbc.
//
// Solidity: function getContractAddress(string contractName) view returns(address)
func (_Registry *RegistrySession) GetContractAddress(contractName string) (common.Address, error) {
	return _Registry.Contract.GetContractAddress(&_Registry.CallOpts, contractName)
}

// GetContractAddress is a free data retrieval call binding the contract method 0x04433bbc.
//
// Solidity: function getContractAddress(string contractName) view returns(address)
func (_Registry *RegistryCallerSession) GetContractAddress(contractName string) (common.Address, error) {
	return _Registry.Contract.GetContractAddress(&_Registry.CallOpts, contractName)
}

// RegistryContractAddressConfirmedIterator is returned from FilterContractAddressConfirmed and is used to iterate over the raw logs and unpacked data for ContractAddressConfirmed events raised by the Registry contract.
type RegistryContractAddressConfirmedIterator struct {
	Event *RegistryContractAddressConfirmed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryContractAddressConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryContractAddressConfirmed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryContractAddressConfirmed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryContractAddressConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryContractAddressConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryContractAddressConfirmed represents a ContractAddressConfirmed event raised by the Registry contract.
type RegistryContractAddressConfirmed struct {
	LockId       [32]byte
	ContractName string
	NewValue     common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterContractAddressConfirmed is a free log retrieval operation binding the contract event 0xbe01e41b60eded6aaa26c1682f75fa4e8edac9df1a3ba2a37258dc6027f00408.
//
// Solidity: event ContractAddressConfirmed(bytes32 lockId, string contractName, address newValue)
func (_Registry *RegistryFilterer) FilterContractAddressConfirmed(opts *bind.FilterOpts) (*RegistryContractAddressConfirmedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ContractAddressConfirmed")
	if err != nil {
		return nil, err
	}
	return &RegistryContractAddressConfirmedIterator{contract: _Registry.contract, event: "ContractAddressConfirmed", logs: logs, sub: sub}, nil
}

// WatchContractAddressConfirmed is a free log subscription operation binding the contract event 0xbe01e41b60eded6aaa26c1682f75fa4e8edac9df1a3ba2a37258dc6027f00408.
//
// Solidity: event ContractAddressConfirmed(bytes32 lockId, string contractName, address newValue)
func (_Registry *RegistryFilterer) WatchContractAddressConfirmed(opts *bind.WatchOpts, sink chan<- *RegistryContractAddressConfirmed) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ContractAddressConfirmed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryContractAddressConfirmed)
				if err := _Registry.contract.UnpackLog(event, "ContractAddressConfirmed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContractAddressConfirmed is a log parse operation binding the contract event 0xbe01e41b60eded6aaa26c1682f75fa4e8edac9df1a3ba2a37258dc6027f00408.
//
// Solidity: event ContractAddressConfirmed(bytes32 lockId, string contractName, address newValue)
func (_Registry *RegistryFilterer) ParseContractAddressConfirmed(log types.Log) (*RegistryContractAddressConfirmed, error) {
	event := new(RegistryContractAddressConfirmed)
	if err := _Registry.contract.UnpackLog(event, "ContractAddressConfirmed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryContractAddressLockedIterator is returned from FilterContractAddressLocked and is used to iterate over the raw logs and unpacked data for ContractAddressLocked events raised by the Registry contract.
type RegistryContractAddressLockedIterator struct {
	Event *RegistryContractAddressLocked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RegistryContractAddressLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryContractAddressLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RegistryContractAddressLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RegistryContractAddressLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryContractAddressLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryContractAddressLocked represents a ContractAddressLocked event raised by the Registry contract.
type RegistryContractAddressLocked struct {
	LockId       [32]byte
	ContractName string
	NewValue     common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterContractAddressLocked is a free log retrieval operation binding the contract event 0x6db4f66c7a48e2b5a06ad599bd86d5147fcf77be050594899e1c38c9da75c8e0.
//
// Solidity: event ContractAddressLocked(bytes32 lockId, string contractName, address newValue)
func (_Registry *RegistryFilterer) FilterContractAddressLocked(opts *bind.FilterOpts) (*RegistryContractAddressLockedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "ContractAddressLocked")
	if err != nil {
		return nil, err
	}
	return &RegistryContractAddressLockedIterator{contract: _Registry.contract, event: "ContractAddressLocked", logs: logs, sub: sub}, nil
}

// WatchContractAddressLocked is a free log subscription operation binding the contract event 0x6db4f66c7a48e2b5a06ad599bd86d5147fcf77be050594899e1c38c9da75c8e0.
//
// Solidity: event ContractAddressLocked(bytes32 lockId, string contractName, address newValue)
func (_Registry *RegistryFilterer) WatchContractAddressLocked(opts *bind.WatchOpts, sink chan<- *RegistryContractAddressLocked) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "ContractAddressLocked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryContractAddressLocked)
				if err := _Registry.contract.UnpackLog(event, "ContractAddressLocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContractAddressLocked is a log parse operation binding the contract event 0x6db4f66c7a48e2b5a06ad599bd86d5147fcf77be050594899e1c38c9da75c8e0.
//
// Solidity: event ContractAddressLocked(bytes32 lockId, string contractName, address newValue)
func (_Registry *RegistryFilterer) ParseContractAddressLocked(log types.Log) (*RegistryContractAddressLocked, error) {
	event := new(RegistryContractAddressLocked)
	if err := _Registry.contract.UnpackLog(event, "ContractAddressLocked", log); err != nil {
		return nil, err
	}
	return event, nil
}
