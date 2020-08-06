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

// SubsSenderBytesABI is the input ABI used to generate the binding from.
const SubsSenderBytesABI = "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newNebulaAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"nebulaAddress\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subId\",\"type\":\"bytes32\"}],\"name\":\"sendValueToSub\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SubsSenderBytes is an auto generated Go binding around an Ethereum contract.
type SubsSenderBytes struct {
	SubsSenderBytesCaller     // Read-only binding to the contract
	SubsSenderBytesTransactor // Write-only binding to the contract
	SubsSenderBytesFilterer   // Log filterer for contract events
}

// SubsSenderBytesCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubsSenderBytesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderBytesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubsSenderBytesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderBytesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubsSenderBytesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderBytesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubsSenderBytesSession struct {
	Contract     *SubsSenderBytes  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SubsSenderBytesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubsSenderBytesCallerSession struct {
	Contract *SubsSenderBytesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SubsSenderBytesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubsSenderBytesTransactorSession struct {
	Contract     *SubsSenderBytesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SubsSenderBytesRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubsSenderBytesRaw struct {
	Contract *SubsSenderBytes // Generic contract binding to access the raw methods on
}

// SubsSenderBytesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubsSenderBytesCallerRaw struct {
	Contract *SubsSenderBytesCaller // Generic read-only contract binding to access the raw methods on
}

// SubsSenderBytesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubsSenderBytesTransactorRaw struct {
	Contract *SubsSenderBytesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubsSenderBytes creates a new instance of SubsSenderBytes, bound to a specific deployed contract.
func NewSubsSenderBytes(address common.Address, backend bind.ContractBackend) (*SubsSenderBytes, error) {
	contract, err := bindSubsSenderBytes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubsSenderBytes{SubsSenderBytesCaller: SubsSenderBytesCaller{contract: contract}, SubsSenderBytesTransactor: SubsSenderBytesTransactor{contract: contract}, SubsSenderBytesFilterer: SubsSenderBytesFilterer{contract: contract}}, nil
}

// NewSubsSenderBytesCaller creates a new read-only instance of SubsSenderBytes, bound to a specific deployed contract.
func NewSubsSenderBytesCaller(address common.Address, caller bind.ContractCaller) (*SubsSenderBytesCaller, error) {
	contract, err := bindSubsSenderBytes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubsSenderBytesCaller{contract: contract}, nil
}

// NewSubsSenderBytesTransactor creates a new write-only instance of SubsSenderBytes, bound to a specific deployed contract.
func NewSubsSenderBytesTransactor(address common.Address, transactor bind.ContractTransactor) (*SubsSenderBytesTransactor, error) {
	contract, err := bindSubsSenderBytes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubsSenderBytesTransactor{contract: contract}, nil
}

// NewSubsSenderBytesFilterer creates a new log filterer instance of SubsSenderBytes, bound to a specific deployed contract.
func NewSubsSenderBytesFilterer(address common.Address, filterer bind.ContractFilterer) (*SubsSenderBytesFilterer, error) {
	contract, err := bindSubsSenderBytes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubsSenderBytesFilterer{contract: contract}, nil
}

// bindSubsSenderBytes binds a generic wrapper to an already deployed contract.
func bindSubsSenderBytes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SubsSenderBytesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubsSenderBytes *SubsSenderBytesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SubsSenderBytes.Contract.SubsSenderBytesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubsSenderBytes *SubsSenderBytesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.SubsSenderBytesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubsSenderBytes *SubsSenderBytesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.SubsSenderBytesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubsSenderBytes *SubsSenderBytesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SubsSenderBytes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubsSenderBytes *SubsSenderBytesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubsSenderBytes *SubsSenderBytesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.contract.Transact(opts, method, params...)
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderBytes *SubsSenderBytesCaller) NebulaAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SubsSenderBytes.contract.Call(opts, out, "nebulaAddress")
	return *ret0, err
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderBytes *SubsSenderBytesSession) NebulaAddress() (common.Address, error) {
	return _SubsSenderBytes.Contract.NebulaAddress(&_SubsSenderBytes.CallOpts)
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderBytes *SubsSenderBytesCallerSession) NebulaAddress() (common.Address, error) {
	return _SubsSenderBytes.Contract.NebulaAddress(&_SubsSenderBytes.CallOpts)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0xa1e9a02a.
//
// Solidity: function sendValueToSub(bytes value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderBytes *SubsSenderBytesTransactor) SendValueToSub(opts *bind.TransactOpts, value []byte, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderBytes.contract.Transact(opts, "sendValueToSub", value, blockNumber, subId)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0xa1e9a02a.
//
// Solidity: function sendValueToSub(bytes value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderBytes *SubsSenderBytesSession) SendValueToSub(value []byte, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.SendValueToSub(&_SubsSenderBytes.TransactOpts, value, blockNumber, subId)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0xa1e9a02a.
//
// Solidity: function sendValueToSub(bytes value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderBytes *SubsSenderBytesTransactorSession) SendValueToSub(value []byte, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.SendValueToSub(&_SubsSenderBytes.TransactOpts, value, blockNumber, subId)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderBytes *SubsSenderBytesTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _SubsSenderBytes.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderBytes *SubsSenderBytesSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.Fallback(&_SubsSenderBytes.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderBytes *SubsSenderBytesTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SubsSenderBytes.Contract.Fallback(&_SubsSenderBytes.TransactOpts, calldata)
}
