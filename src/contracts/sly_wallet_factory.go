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

// Reference imports to suppress slyerrors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SLYWalletFactoryABI is the input ABI used to generate the binding from.
const SLYWalletFactoryABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_diamondCutFacet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_diamondLoupeFacet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_slyWalletFacet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_slyDiamondInit\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"diamond\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"}],\"name\":\"SLYWalletCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"createSLYWallet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_salt\",\"type\":\"bytes32\"}],\"name\":\"createSLYWalletWithSalt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"diamondCutFacet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"diamondLoupeFacet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllWallets\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"getWalletsByOwner\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWalletsCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_salt\",\"type\":\"bytes32\"}],\"name\":\"predictWalletAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"predictedAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slyDiamondInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slyWalletFacet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// SLYWalletFactory is an auto generated Go binding around an Ethereum contract.
type SLYWalletFactory struct {
	SLYWalletFactoryCaller     // Read-only binding to the contract
	SLYWalletFactoryTransactor // Write-only binding to the contract
	SLYWalletFactoryFilterer   // Log filterer for contract events
}

// SLYWalletFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SLYWalletFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLYWalletFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SLYWalletFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLYWalletFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SLYWalletFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLYWalletFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SLYWalletFactorySession struct {
	Contract     *SLYWalletFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SLYWalletFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SLYWalletFactoryCallerSession struct {
	Contract *SLYWalletFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SLYWalletFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SLYWalletFactoryTransactorSession struct {
	Contract     *SLYWalletFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// NewSLYWalletFactory creates a new instance of SLYWalletFactory, bound to a specific deployed contract.
func NewSLYWalletFactory(address common.Address, backend bind.ContractBackend) (*SLYWalletFactory, error) {
	contract, err := bindSLYWalletFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SLYWalletFactory{SLYWalletFactoryCaller: SLYWalletFactoryCaller{contract: contract}, SLYWalletFactoryTransactor: SLYWalletFactoryTransactor{contract: contract}, SLYWalletFactoryFilterer: SLYWalletFactoryFilterer{contract: contract}}, nil
}

// NewSLYWalletFactoryCaller creates a new read-only instance of SLYWalletFactory, bound to a specific deployed contract.
func NewSLYWalletFactoryCaller(address common.Address, caller bind.ContractCaller) (*SLYWalletFactoryCaller, error) {
	contract, err := bindSLYWalletFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SLYWalletFactoryCaller{contract: contract}, nil
}

// NewSLYWalletFactoryTransactor creates a new write-only instance of SLYWalletFactory, bound to a specific deployed contract.
func NewSLYWalletFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*SLYWalletFactoryTransactor, error) {
	contract, err := bindSLYWalletFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SLYWalletFactoryTransactor{contract: contract}, nil
}

// NewSLYWalletFactoryFilterer creates a new log filterer instance of SLYWalletFactory, bound to a specific deployed contract.
func NewSLYWalletFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*SLYWalletFactoryFilterer, error) {
	contract, err := bindSLYWalletFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SLYWalletFactoryFilterer{contract: contract}, nil
}

// bindSLYWalletFactory binds a generic wrapper to an already deployed contract.
func bindSLYWalletFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SLYWalletFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// DiamondCutFacet is a free data retrieval call binding the contract method 0x53a8825f.
func (_SLYWalletFactory *SLYWalletFactoryCaller) DiamondCutFacet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "diamondCutFacet")
	if err != nil {
		return *new(common.Address), err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// DiamondLoupeFacet is a free data retrieval call binding the contract method 0xa8cbe91a.
func (_SLYWalletFactory *SLYWalletFactoryCaller) DiamondLoupeFacet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "diamondLoupeFacet")
	if err != nil {
		return *new(common.Address), err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// GetAllWallets is a free data retrieval call binding the contract method 0x98f40e35.
func (_SLYWalletFactory *SLYWalletFactoryCaller) GetAllWallets(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "getAllWallets")
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address), nil
}

// GetWalletsByOwner is a free data retrieval call binding the contract method 0xfea80399.
func (_SLYWalletFactory *SLYWalletFactoryCaller) GetWalletsByOwner(opts *bind.CallOpts, _owner common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "getWalletsByOwner", _owner)
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address), nil
}

// GetWalletsCount is a free data retrieval call binding the contract method 0xa96abd9d.
func (_SLYWalletFactory *SLYWalletFactoryCaller) GetWalletsCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "getWalletsCount")
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}

// PredictWalletAddress is a free data retrieval call binding the contract method 0x0a43c8ea.
func (_SLYWalletFactory *SLYWalletFactoryCaller) PredictWalletAddress(opts *bind.CallOpts, _salt [32]byte) (common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "predictWalletAddress", _salt)
	if err != nil {
		return *new(common.Address), err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// SlyDiamondInit is a free data retrieval call binding the contract method 0xf3408064.
func (_SLYWalletFactory *SLYWalletFactoryCaller) SlyDiamondInit(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "slyDiamondInit")
	if err != nil {
		return *new(common.Address), err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// SlyWalletFacet is a free data retrieval call binding the contract method 0x1a34b0da.
func (_SLYWalletFactory *SLYWalletFactoryCaller) SlyWalletFacet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SLYWalletFactory.contract.Call(opts, &out, "slyWalletFacet")
	if err != nil {
		return *new(common.Address), err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// CreateSLYWallet is a paid mutator transaction binding the contract method 0xaa69c1ee.
func (_SLYWalletFactory *SLYWalletFactoryTransactor) CreateSLYWallet(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _SLYWalletFactory.contract.Transact(opts, "createSLYWallet", _owner)
}

// CreateSLYWalletWithSalt is a paid mutator transaction binding the contract method 0x72fe1f87.
func (_SLYWalletFactory *SLYWalletFactoryTransactor) CreateSLYWalletWithSalt(opts *bind.TransactOpts, _owner common.Address, _salt [32]byte) (*types.Transaction, error) {
	return _SLYWalletFactory.contract.Transact(opts, "createSLYWalletWithSalt", _owner, _salt)
}

// SLYWalletCreatedIterator is returned from FilterSLYWalletCreated and is used to iterate over the raw logs and unpacked data for SLYWalletCreated events raised by the SLYWalletFactory contract.
type SLYWalletCreatedIterator struct {
	Event *SLYWalletCreatedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for slyerrors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SLYWalletCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SLYWalletCreatedEvent)
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
		it.Event = new(SLYWalletCreatedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return false
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SLYWalletCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SLYWalletCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SLYWalletCreatedEvent represents a SLYWalletCreated event raised by the SLYWalletFactory contract.
type SLYWalletCreatedEvent struct {
	Diamond common.Address
	Owner   common.Address
	Creator common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSLYWalletCreated is a free log retrieval operation binding the contract event 0x69e8a87ea6f86f8bfd222fced5bb8c0b4e2180155bfabebbe01f6d43eaa28fd5.
//
// Solidity: event SLYWalletCreated(address indexed diamond, address indexed owner, address indexed creator)
func (_SLYWalletFactory *SLYWalletFactoryFilterer) FilterSLYWalletCreated(opts *bind.FilterOpts, diamond []common.Address, owner []common.Address, creator []common.Address) (*SLYWalletCreatedIterator, error) {
	var diamondRule []interface{}
	for _, diamondItem := range diamond {
		diamondRule = append(diamondRule, diamondItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SLYWalletFactory.contract.FilterLogs(opts, "SLYWalletCreated", diamondRule, ownerRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SLYWalletCreatedIterator{contract: _SLYWalletFactory.contract, event: "SLYWalletCreated", logs: logs, sub: sub}, nil
}

// WatchSLYWalletCreated is a free log subscription operation binding the contract event 0x69e8a87ea6f86f8bfd222fced5bb8c0b4e2180155bfabebbe01f6d43eaa28fd5.
//
// Solidity: event SLYWalletCreated(address indexed diamond, address indexed owner, address indexed creator)
func (_SLYWalletFactory *SLYWalletFactoryFilterer) WatchSLYWalletCreated(opts *bind.WatchOpts, sink chan<- *SLYWalletCreatedEvent, diamond []common.Address, owner []common.Address, creator []common.Address) (ethereum.Subscription, error) {
	var diamondRule []interface{}
	for _, diamondItem := range diamond {
		diamondRule = append(diamondRule, diamondItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SLYWalletFactory.contract.WatchLogs(opts, "SLYWalletCreated", diamondRule, ownerRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SLYWalletCreatedEvent)
				if err := _SLYWalletFactory.contract.UnpackLog(event, "SLYWalletCreated", log); err != nil {
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
