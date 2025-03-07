// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SLYWalletABI is the input ABI used to generate the binding from.
const SLYWalletABI = "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AssetWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"destinations\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"BatchExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"key\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enum ISLYWallet.Role\",\"name\":\"role\",\"type\":\"uint8\"}],\"name\":\"KeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"key\",\"type\":\"address\"}],\"name\":\"KeyRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"relayer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"functionSelector\",\"type\":\"bytes4\"}],\"name\":\"MetaTransactionExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"NonceUsed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"key\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enum ISLYWallet.Permission\",\"name\":\"permission\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"PermissionChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"}],\"name\":\"SignatureValidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransactionExecuted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"},{\"internalType\":\"enum ISLYWallet.Role\",\"name\":\"_role\",\"type\":\"uint8\"}],\"name\":\"addKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_to\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_value\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"_data\",\"type\":\"bytes[]\"}],\"name\":\"executeBatch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_to\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_value\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"_data\",\"type\":\"bytes[]\"},{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"executeBatchWithSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"executeWithSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"getKeyRole\",\"outputs\":[{\"internalType\":\"enum ISLYWallet.Role\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enum ISLYWallet.Role\",\"name\":\"_role\",\"type\":\"uint8\"}],\"name\":\"getKeysByRole\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"}],\"name\":\"getNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"},{\"internalType\":\"enum ISLYWallet.Permission\",\"name\":\"_permission\",\"type\":\"uint8\"}],\"name\":\"hasPermission\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"magicValue\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"keyExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"removeKey\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address payable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawETH\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SLYWallet is an auto generated Go binding around an Ethereum contract.
type SLYWallet struct {
	SLYWalletCaller     // Read-only binding to the contract
	SLYWalletTransactor // Write-only binding to the contract
	SLYWalletFilterer   // Log filterer for contract events
}

// SLYWalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type SLYWalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLYWalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SLYWalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLYWalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SLYWalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLYWalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SLYWalletSession struct {
	Contract     *SLYWallet        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SLYWalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SLYWalletCallerSession struct {
	Contract *SLYWalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SLYWalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SLYWalletTransactorSession struct {
	Contract     *SLYWalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// NewSLYWallet creates a new instance of SLYWallet, bound to a specific deployed contract.
func NewSLYWallet(address common.Address, backend bind.ContractBackend) (*SLYWallet, error) {
	contract, err := bindSLYWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SLYWallet{SLYWalletCaller: SLYWalletCaller{contract: contract}, SLYWalletTransactor: SLYWalletTransactor{contract: contract}, SLYWalletFilterer: SLYWalletFilterer{contract: contract}}, nil
}

// NewSLYWalletCaller creates a new read-only instance of SLYWallet, bound to a specific deployed contract.
func NewSLYWalletCaller(address common.Address, caller bind.ContractCaller) (*SLYWalletCaller, error) {
	contract, err := bindSLYWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SLYWalletCaller{contract: contract}, nil
}

// NewSLYWalletTransactor creates a new write-only instance of SLYWallet, bound to a specific deployed contract.
func NewSLYWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*SLYWalletTransactor, error) {
	contract, err := bindSLYWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SLYWalletTransactor{contract: contract}, nil
}

// NewSLYWalletFilterer creates a new log filterer instance of SLYWallet, bound to a specific deployed contract.
func NewSLYWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*SLYWalletFilterer, error) {
	contract, err := bindSLYWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SLYWalletFilterer{contract: contract}, nil
}

// bindSLYWallet binds a generic wrapper to an already deployed contract.
func bindSLYWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SLYWalletABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Role represents the role type in the SLYWallet contract
type Role uint8

const (
	RoleNone          Role = 0
	RoleOwner         Role = 1
	RoleAdmin         Role = 2
	RoleAuthenticator Role = 3
	RoleRoleCount     Role = 4
)

// Permission represents the permission type in the SLYWallet contract
type Permission uint8

const (
	PermissionAddKey           Permission = 0
	PermissionRemoveKey        Permission = 1
	PermissionExecute          Permission = 2
	PermissionExecuteBatch     Permission = 3
	PermissionValidateSignature Permission = 4
	PermissionWithdrawAssets   Permission = 5
	PermissionDiamondCut       Permission = 6
)

// AddKey is a paid mutator transaction binding the contract method 0xd8f4b32e.
func (_SLYWallet *SLYWalletTransactor) AddKey(opts *bind.TransactOpts, _key common.Address, _role uint8) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "addKey", _key, _role)
}

// Execute is a paid mutator transaction binding the contract method 0x1cff79cd.
func (_SLYWallet *SLYWalletTransactor) Execute(opts *bind.TransactOpts, _to common.Address, _value *big.Int, _data []byte) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "execute", _to, _value, _data)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x34eecd8c.
func (_SLYWallet *SLYWalletTransactor) ExecuteBatch(opts *bind.TransactOpts, _to []common.Address, _value []*big.Int, _data [][]byte) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "executeBatch", _to, _value, _data)
}

// ExecuteBatchWithSignature is a paid mutator transaction binding the contract method 0x8c7f75a0.
func (_SLYWallet *SLYWalletTransactor) ExecuteBatchWithSignature(opts *bind.TransactOpts, _to []common.Address, _value []*big.Int, _data [][]byte, _signer common.Address, _nonce *big.Int, _signature []byte) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "executeBatchWithSignature", _to, _value, _data, _signer, _nonce, _signature)
}

// ExecuteWithSignature is a paid mutator transaction binding the contract method 0x03e7ad81.
func (_SLYWallet *SLYWalletTransactor) ExecuteWithSignature(opts *bind.TransactOpts, _to common.Address, _value *big.Int, _data []byte, _signer common.Address, _nonce *big.Int, _signature []byte) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "executeWithSignature", _to, _value, _data, _signer, _nonce, _signature)
}

// GetKeyRole is a free data retrieval call binding the contract method 0xdd8c132d.
func (_SLYWallet *SLYWalletCaller) GetKeyRole(opts *bind.CallOpts, _key common.Address) (uint8, error) {
	var out []interface{}
	err := _SLYWallet.contract.Call(opts, &out, "getKeyRole", _key)
	if err != nil {
		return 0, err
	}
	return *abi.ConvertType(out[0], new(uint8)).(*uint8), nil
}

// GetKeysByRole is a free data retrieval call binding the contract method 0x7d32e7bd.
func (_SLYWallet *SLYWalletCaller) GetKeysByRole(opts *bind.CallOpts, _role uint8) ([]common.Address, error) {
	var out []interface{}
	err := _SLYWallet.contract.Call(opts, &out, "getKeysByRole", _role)
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address), nil
}

// GetNonce is a free data retrieval call binding the contract method 0x2d0335ab.
func (_SLYWallet *SLYWalletCaller) GetNonce(opts *bind.CallOpts, _signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SLYWallet.contract.Call(opts, &out, "getNonce", _signer)
	if err != nil {
		return nil, err
	}
	return *abi.ConvertType(out[0], new(*big.Int)).(**big.Int), nil
}

// HasPermission is a free data retrieval call binding the contract method 0x30f01a13.
func (_SLYWallet *SLYWalletCaller) HasPermission(opts *bind.CallOpts, _key common.Address, _permission uint8) (bool, error) {
	var out []interface{}
	err := _SLYWallet.contract.Call(opts, &out, "hasPermission", _key, _permission)
	if err != nil {
		return false, err
	}
	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

// IsValidSignature is a free data retrieval call binding the contract method 0x1626ba7e.
func (_SLYWallet *SLYWalletCaller) IsValidSignature(opts *bind.CallOpts, _hash [32]byte, _signature []byte) ([4]byte, error) {
	var out []interface{}
	err := _SLYWallet.contract.Call(opts, &out, "isValidSignature", _hash, _signature)
	if err != nil {
		return [4]byte{}, err
	}
	return *abi.ConvertType(out[0], new([4]byte)).(*[4]byte), nil
}

// KeyExists is a free data retrieval call binding the contract method 0x83ae4806.
func (_SLYWallet *SLYWalletCaller) KeyExists(opts *bind.CallOpts, _key common.Address) (bool, error) {
	var out []interface{}
	err := _SLYWallet.contract.Call(opts, &out, "keyExists", _key)
	if err != nil {
		return false, err
	}
	return *abi.ConvertType(out[0], new(bool)).(*bool), nil
}

// RemoveKey is a paid mutator transaction binding the contract method 0x5c81979d.
func (_SLYWallet *SLYWalletTransactor) RemoveKey(opts *bind.TransactOpts, _key common.Address) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "removeKey", _key)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0x7cc19e70.
func (_SLYWallet *SLYWalletTransactor) WithdrawERC20(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "withdrawERC20", _token, _to, _amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x594591a4.
func (_SLYWallet *SLYWalletTransactor) WithdrawETH(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _SLYWallet.contract.Transact(opts, "withdrawETH", _to, _amount)
}

// AssetWithdrawnIterator is returned from FilterAssetWithdrawn and is used to iterate over the raw logs and unpacked data for AssetWithdrawn events raised by the SLYWallet contract.
type AssetWithdrawnIterator struct {
	Event *AssetWithdrawnEvent // Event containing the contract specifics and raw log

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
func (it *AssetWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AssetWithdrawnEvent)
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
		it.Event = new(AssetWithdrawnEvent)
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
func (it *AssetWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AssetWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AssetWithdrawnEvent represents a AssetWithdrawn event raised by the SLYWallet contract.
type AssetWithdrawnEvent struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAssetWithdrawn is a free log retrieval operation binding the contract event 0x0e13b9b38f18a10a3b231f62177f184c5809cb7ea54431297f9370ea53e4c9f3.
//
// Solidity: event AssetWithdrawn(address indexed token, address indexed to, uint256 amount)
func (_SLYWallet *SLYWalletFilterer) FilterAssetWithdrawn(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*AssetWithdrawnIterator, error) {
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SLYWallet.contract.FilterLogs(opts, "AssetWithdrawn", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AssetWithdrawnIterator{contract: _SLYWallet.contract, event: "AssetWithdrawn", logs: logs, sub: sub}, nil
}

// WatchAssetWithdrawn is a free log subscription operation binding the contract event 0x0e13b9b38f18a10a3b231f62177f184c5809cb7ea54431297f9370ea53e4c9f3.
//
// Solidity: event AssetWithdrawn(address indexed token, address indexed to, uint256 amount)
func (_SLYWallet *SLYWalletFilterer) WatchAssetWithdrawn(opts *bind.WatchOpts, sink chan<- *AssetWithdrawnEvent, token []common.Address, to []common.Address) (ethereum.Subscription, error) {
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SLYWallet.contract.WatchLogs(opts, "AssetWithdrawn", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AssetWithdrawnEvent)
				if err := _SLYWallet.contract.UnpackLog(event, "AssetWithdrawn", log); err != nil {
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

// BatchExecutedEvent represents a BatchExecuted event raised by the SLYWallet contract.
type BatchExecutedEvent struct {
	Destinations []common.Address
	Values       []*big.Int
	Data         [][]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBatchExecuted is a free log retrieval operation binding the contract event 0x53b7f7f02c34ed73e56bc3b29dfe2f3a37447ce307c0d29626df25bb8ce9c66d.
//
// Solidity: event BatchExecuted(address[] destinations, uint256[] values, bytes[] data)
func (_SLYWallet *SLYWalletFilterer) FilterBatchExecuted(opts *bind.FilterOpts) (*BatchExecutedIterator, error) {
	logs, sub, err := _SLYWallet.contract.FilterLogs(opts, "BatchExecuted")
	if err != nil {
		return nil, err
	}
	return &BatchExecutedIterator{contract: _SLYWallet.contract, event: "BatchExecuted", logs: logs, sub: sub}, nil
}

// BatchExecutedIterator is returned from FilterBatchExecuted and is used to iterate over the raw logs and unpacked data for BatchExecuted events raised by the SLYWallet contract.
type BatchExecutedIterator struct {
	Event *BatchExecutedEvent // Event containing the contract specifics and raw log

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
func (it *BatchExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchExecutedEvent)
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
		it.Event = new(BatchExecutedEvent)
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
func (it *BatchExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WatchBatchExecuted is a free log subscription operation binding the contract event 0x53b7f7f02c34ed73e56bc3b29dfe2f3a37447ce307c0d29626df25bb8ce9c66d.
//
// Solidity: event BatchExecuted(address[] destinations, uint256[] values, bytes[] data)
func (_SLYWallet *SLYWalletFilterer) WatchBatchExecuted(opts *bind.WatchOpts, sink chan<- *BatchExecutedEvent) (ethereum.Subscription, error) {
	logs, sub, err := _SLYWallet.contract.WatchLogs(opts, "BatchExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchExecutedEvent)
				if err := _SLYWallet.contract.UnpackLog(event, "BatchExecuted", log); err != nil {
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

// KeyAddedEvent represents a KeyAdded event raised by the SLYWallet contract.
type KeyAddedEvent struct {
	Key  common.Address
	Role uint8
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterKeyAdded is a free log retrieval operation binding the contract event 0x3aedc386ed11227c7fab686c9352ed414fd187d3058783e99951f0f87ef2d41d.
//
// Solidity: event KeyAdded(address indexed key, uint8 role)
func (_SLYWallet *SLYWalletFilterer) FilterKeyAdded(opts *bind.FilterOpts, key []common.Address) (*KeyAddedIterator, error) {
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _SLYWallet.contract.FilterLogs(opts, "KeyAdded", keyRule)
	if err != nil {
		return nil, err
	}
	return &KeyAddedIterator{contract: _SLYWallet.contract, event: "KeyAdded", logs: logs, sub: sub}, nil
}

// KeyAddedIterator is returned from FilterKeyAdded and is used to iterate over the raw logs and unpacked data for KeyAdded events raised by the SLYWallet contract.
type KeyAddedIterator struct {
	Event *KeyAddedEvent // Event containing the contract specifics and raw log

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
func (it *KeyAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeyAddedEvent)
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
		it.Event = new(KeyAddedEvent)
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
func (it *KeyAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WatchKeyAdded is a free log subscription operation binding the contract event 0x3aedc386ed11227c7fab686c9352ed414fd187d3058783e99951f0f87ef2d41d.
//
// Solidity: event KeyAdded(address indexed key, uint8 role)
func (_SLYWallet *SLYWalletFilterer) WatchKeyAdded(opts *bind.WatchOpts, sink chan<- *KeyAddedEvent, key []common.Address) (ethereum.Subscription, error) {
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _SLYWallet.contract.WatchLogs(opts, "KeyAdded", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeyAddedEvent)
				if err := _SLYWallet.contract.UnpackLog(event, "KeyAdded", log); err != nil {
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

// KeyRemovedEvent represents a KeyRemoved event raised by the SLYWallet contract.
type KeyRemovedEvent struct {
	Key common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterKeyRemoved is a free log retrieval operation binding the contract event 0x3d2d9335cf2adf5b43003476a2376ce36f969f84d2cffdd599b6a0c1f6e32ebb.
//
// Solidity: event KeyRemoved(address indexed key)
func (_SLYWallet *SLYWalletFilterer) FilterKeyRemoved(opts *bind.FilterOpts, key []common.Address) (*KeyRemovedIterator, error) {
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _SLYWallet.contract.FilterLogs(opts, "KeyRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return &KeyRemovedIterator{contract: _SLYWallet.contract, event: "KeyRemoved", logs: logs, sub: sub}, nil
}

// KeyRemovedIterator is returned from FilterKeyRemoved and is used to iterate over the raw logs and unpacked data for KeyRemoved events raised by the SLYWallet contract.
type KeyRemovedIterator struct {
	Event *KeyRemovedEvent // Event containing the contract specifics and raw log

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
func (it *KeyRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeyRemovedEvent)
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
		it.Event = new(KeyRemovedEvent)
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
func (it *KeyRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeyRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WatchKeyRemoved is a free log subscription operation binding the contract event 0x3d2d9335cf2adf5b43003476a2376ce36f969f84d2cffdd599b6a0c1f6e32ebb.
//
// Solidity: event KeyRemoved(address indexed key)
func (_SLYWallet *SLYWalletFilterer) WatchKeyRemoved(opts *bind.WatchOpts, sink chan<- *KeyRemovedEvent, key []common.Address) (ethereum.Subscription, error) {
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _SLYWallet.contract.WatchLogs(opts, "KeyRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeyRemovedEvent)
				if err := _SLYWallet.contract.UnpackLog(event, "KeyRemoved", log); err != nil {
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

// MetaTransactionExecutedEvent represents a MetaTransactionExecuted event raised by the SLYWallet contract.
type MetaTransactionExecutedEvent struct {
	Signer          common.Address
	Relayer         common.Address
	FunctionSelector [4]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMetaTransactionExecuted is a free log retrieval operation binding the contract event 0x5845892132946850460bff5a0083f71031bc5bf9aadcd40f1de79423eac9b10b.
//
// Solidity: event MetaTransactionExecuted(address indexed signer, address indexed relayer, bytes4 functionSelector)
func (_SLYWallet *SLYWalletFilterer) FilterMetaTransactionExecuted(opts *bind.FilterOpts, signer []common.Address, relayer []common.Address) (*MetaTransactionExecutedIterator, error) {
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var relayerRule []interface{}
	for _, relayerItem := range relayer {
		relayerRule = append(relayerRule, relayerItem)
	}

	logs, sub, err := _SLYWallet.contract.FilterLogs(opts, "MetaTransactionExecuted", signerRule, relayerRule)
	if err != nil {
		return nil, err
	}
	return &MetaTransactionExecutedIterator{contract: _SLYWallet.contract, event: "MetaTransactionExecuted", logs: logs, sub: sub}, nil
}

// MetaTransactionExecutedIterator is returned from FilterMetaTransactionExecuted and is used to iterate over the raw logs and unpacked data for MetaTransactionExecuted events raised by the SLYWallet contract.
type MetaTransactionExecutedIterator struct {
	Event *MetaTransactionExecutedEvent // Event containing the contract specifics and raw log

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
func (it *MetaTransactionExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetaTransactionExecutedEvent)
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
		it.Event = new(MetaTransactionExecutedEvent)
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
func (it *MetaTransactionExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetaTransactionExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WatchMetaTransactionExecuted is a free log subscription operation binding the contract event 0x5845892132946850460bff5a0083f71031bc5bf9aadcd40f1de79423eac9b10b.
//
// Solidity: event MetaTransactionExecuted(address indexed signer, address indexed relayer, bytes4 functionSelector)
func (_SLYWallet *SLYWalletFilterer) WatchMetaTransactionExecuted(opts *bind.WatchOpts, sink chan<- *MetaTransactionExecutedEvent, signer []common.Address, relayer []common.Address) (ethereum.Subscription, error) {
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var relayerRule []interface{}
	for _, relayerItem := range relayer {
		relayerRule = append(relayerRule, relayerItem)
	}

	logs, sub, err := _SLYWallet.contract.WatchLogs(opts, "MetaTransactionExecuted", signerRule, relayerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetaTransactionExecutedEvent)
				if err := _SLYWallet.contract.UnpackLog(event, "MetaTransactionExecuted", log); err != nil {
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

// NonceUsedEvent represents a NonceUsed event raised by the SLYWallet contract.
type NonceUsedEvent struct {
	Signer common.Address
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNonceUsed is a free log retrieval operation binding the contract event 0x1ede735f8f4baf2693e43546c3c67549a799c2f8573aeb8cc58917a0d7df9a9a.
//
// Solidity: event NonceUsed(address indexed signer, uint256 nonce)
func (_SLYWallet *SLYWalletFilterer) FilterNonceUsed(opts *bind.FilterOpts, signer []common.Address) (*NonceUsedIterator, error) {
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _SLYWallet.contract.FilterLogs(opts, "NonceUsed", signerRule)
	if err != nil {
		return nil, err
	}
	return &NonceUsedIterator{contract: _SLYWallet.contract, event: "NonceUsed", logs: logs, sub: sub}, nil
}

// NonceUsedIterator is returned from FilterNonceUsed and is used to iterate over the raw logs and unpacked data for NonceUsed events raised by the SLYWallet contract.
type NonceUsedIterator struct {
	Event *NonceUsedEvent // Event containing the contract specifics and raw log

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
func (it *NonceUsedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NonceUsedEvent)
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
		it.Event = new(NonceUsedEvent)
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
func (it *NonceUsedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NonceUsedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WatchNonceUsed is a free log subscription operation binding the contract event 0x1ede735f8f4baf2693e43546c3c67549a799c2f8573aeb8cc58917a0d7df9a9a.
//
// Solidity: event NonceUsed(address indexed signer, uint256 nonce)
func (_SLYWallet *SLYWalletFilterer) WatchNonceUsed(opts *bind.WatchOpts, sink chan<- *NonceUsedEvent, signer []common.Address) (ethereum.Subscription, error) {
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _SLYWallet.contract.WatchLogs(opts, "NonceUsed", signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NonceUsedEvent)
				if err := _SLYWallet.contract.UnpackLog(event, "NonceUsed", log); err != nil {
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