// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sender

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

// SubsSenderStringABI is the input ABI used to generate the binding from.
const SubsSenderStringABI = "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newNebulaAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"nebulaAddress\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subId\",\"type\":\"bytes32\"}],\"name\":\"sendValueToSub\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SubsSenderString is an auto generated Go binding around an Ethereum contract.
type SubsSenderString struct {
	SubsSenderStringCaller     // Read-only binding to the contract
	SubsSenderStringTransactor // Write-only binding to the contract
	SubsSenderStringFilterer   // Log filterer for contract events
}

// SubsSenderStringCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubsSenderStringCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderStringTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubsSenderStringTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderStringFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubsSenderStringFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderStringSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubsSenderStringSession struct {
	Contract     *SubsSenderString // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SubsSenderStringCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubsSenderStringCallerSession struct {
	Contract *SubsSenderStringCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SubsSenderStringTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubsSenderStringTransactorSession struct {
	Contract     *SubsSenderStringTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SubsSenderStringRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubsSenderStringRaw struct {
	Contract *SubsSenderString // Generic contract binding to access the raw methods on
}

// SubsSenderStringCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubsSenderStringCallerRaw struct {
	Contract *SubsSenderStringCaller // Generic read-only contract binding to access the raw methods on
}

// SubsSenderStringTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubsSenderStringTransactorRaw struct {
	Contract *SubsSenderStringTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubsSenderString creates a new instance of SubsSenderString, bound to a specific deployed contract.
func NewSubsSenderString(address common.Address, backend bind.ContractBackend) (*SubsSenderString, error) {
	contract, err := bindSubsSenderString(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubsSenderString{SubsSenderStringCaller: SubsSenderStringCaller{contract: contract}, SubsSenderStringTransactor: SubsSenderStringTransactor{contract: contract}, SubsSenderStringFilterer: SubsSenderStringFilterer{contract: contract}}, nil
}

// NewSubsSenderStringCaller creates a new read-only instance of SubsSenderString, bound to a specific deployed contract.
func NewSubsSenderStringCaller(address common.Address, caller bind.ContractCaller) (*SubsSenderStringCaller, error) {
	contract, err := bindSubsSenderString(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubsSenderStringCaller{contract: contract}, nil
}

// NewSubsSenderStringTransactor creates a new write-only instance of SubsSenderString, bound to a specific deployed contract.
func NewSubsSenderStringTransactor(address common.Address, transactor bind.ContractTransactor) (*SubsSenderStringTransactor, error) {
	contract, err := bindSubsSenderString(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubsSenderStringTransactor{contract: contract}, nil
}

// NewSubsSenderStringFilterer creates a new log filterer instance of SubsSenderString, bound to a specific deployed contract.
func NewSubsSenderStringFilterer(address common.Address, filterer bind.ContractFilterer) (*SubsSenderStringFilterer, error) {
	contract, err := bindSubsSenderString(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubsSenderStringFilterer{contract: contract}, nil
}

// bindSubsSenderString binds a generic wrapper to an already deployed contract.
func bindSubsSenderString(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SubsSenderStringABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubsSenderString *SubsSenderStringRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SubsSenderString.Contract.SubsSenderStringCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubsSenderString *SubsSenderStringRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubsSenderString.Contract.SubsSenderStringTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubsSenderString *SubsSenderStringRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubsSenderString.Contract.SubsSenderStringTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubsSenderString *SubsSenderStringCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SubsSenderString.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubsSenderString *SubsSenderStringTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubsSenderString.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubsSenderString *SubsSenderStringTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubsSenderString.Contract.contract.Transact(opts, method, params...)
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderString *SubsSenderStringCaller) NebulaAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SubsSenderString.contract.Call(opts, out, "nebulaAddress")
	return *ret0, err
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderString *SubsSenderStringSession) NebulaAddress() (common.Address, error) {
	return _SubsSenderString.Contract.NebulaAddress(&_SubsSenderString.CallOpts)
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderString *SubsSenderStringCallerSession) NebulaAddress() (common.Address, error) {
	return _SubsSenderString.Contract.NebulaAddress(&_SubsSenderString.CallOpts)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0x781be6d1.
//
// Solidity: function sendValueToSub(string value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderString *SubsSenderStringTransactor) SendValueToSub(opts *bind.TransactOpts, value string, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderString.contract.Transact(opts, "sendValueToSub", value, blockNumber, subId)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0x781be6d1.
//
// Solidity: function sendValueToSub(string value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderString *SubsSenderStringSession) SendValueToSub(value string, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderString.Contract.SendValueToSub(&_SubsSenderString.TransactOpts, value, blockNumber, subId)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0x781be6d1.
//
// Solidity: function sendValueToSub(string value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderString *SubsSenderStringTransactorSession) SendValueToSub(value string, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderString.Contract.SendValueToSub(&_SubsSenderString.TransactOpts, value, blockNumber, subId)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderString *SubsSenderStringTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _SubsSenderString.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderString *SubsSenderStringSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SubsSenderString.Contract.Fallback(&_SubsSenderString.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderString *SubsSenderStringTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SubsSenderString.Contract.Fallback(&_SubsSenderString.TransactOpts, calldata)
}
