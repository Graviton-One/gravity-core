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

// SubsSenderIntABI is the input ABI used to generate the binding from.
const SubsSenderIntABI = "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"newNebulaAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"nebulaAddress\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"int64\",\"name\":\"value\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subId\",\"type\":\"bytes32\"}],\"name\":\"sendValueToSub\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SubsSenderInt is an auto generated Go binding around an Ethereum contract.
type SubsSenderInt struct {
	SubsSenderIntCaller     // Read-only binding to the contract
	SubsSenderIntTransactor // Write-only binding to the contract
	SubsSenderIntFilterer   // Log filterer for contract events
}

// SubsSenderIntCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubsSenderIntCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderIntTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubsSenderIntTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderIntFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubsSenderIntFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubsSenderIntSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubsSenderIntSession struct {
	Contract     *SubsSenderInt    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SubsSenderIntCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubsSenderIntCallerSession struct {
	Contract *SubsSenderIntCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SubsSenderIntTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubsSenderIntTransactorSession struct {
	Contract     *SubsSenderIntTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SubsSenderIntRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubsSenderIntRaw struct {
	Contract *SubsSenderInt // Generic contract binding to access the raw methods on
}

// SubsSenderIntCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubsSenderIntCallerRaw struct {
	Contract *SubsSenderIntCaller // Generic read-only contract binding to access the raw methods on
}

// SubsSenderIntTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubsSenderIntTransactorRaw struct {
	Contract *SubsSenderIntTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubsSenderInt creates a new instance of SubsSenderInt, bound to a specific deployed contract.
func NewSubsSenderInt(address common.Address, backend bind.ContractBackend) (*SubsSenderInt, error) {
	contract, err := bindSubsSenderInt(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubsSenderInt{SubsSenderIntCaller: SubsSenderIntCaller{contract: contract}, SubsSenderIntTransactor: SubsSenderIntTransactor{contract: contract}, SubsSenderIntFilterer: SubsSenderIntFilterer{contract: contract}}, nil
}

// NewSubsSenderIntCaller creates a new read-only instance of SubsSenderInt, bound to a specific deployed contract.
func NewSubsSenderIntCaller(address common.Address, caller bind.ContractCaller) (*SubsSenderIntCaller, error) {
	contract, err := bindSubsSenderInt(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubsSenderIntCaller{contract: contract}, nil
}

// NewSubsSenderIntTransactor creates a new write-only instance of SubsSenderInt, bound to a specific deployed contract.
func NewSubsSenderIntTransactor(address common.Address, transactor bind.ContractTransactor) (*SubsSenderIntTransactor, error) {
	contract, err := bindSubsSenderInt(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubsSenderIntTransactor{contract: contract}, nil
}

// NewSubsSenderIntFilterer creates a new log filterer instance of SubsSenderInt, bound to a specific deployed contract.
func NewSubsSenderIntFilterer(address common.Address, filterer bind.ContractFilterer) (*SubsSenderIntFilterer, error) {
	contract, err := bindSubsSenderInt(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubsSenderIntFilterer{contract: contract}, nil
}

// bindSubsSenderInt binds a generic wrapper to an already deployed contract.
func bindSubsSenderInt(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SubsSenderIntABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubsSenderInt *SubsSenderIntRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SubsSenderInt.Contract.SubsSenderIntCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubsSenderInt *SubsSenderIntRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.SubsSenderIntTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubsSenderInt *SubsSenderIntRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.SubsSenderIntTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubsSenderInt *SubsSenderIntCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SubsSenderInt.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubsSenderInt *SubsSenderIntTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubsSenderInt *SubsSenderIntTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.contract.Transact(opts, method, params...)
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderInt *SubsSenderIntCaller) NebulaAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SubsSenderInt.contract.Call(opts, out, "nebulaAddress")
	return *ret0, err
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderInt *SubsSenderIntSession) NebulaAddress() (common.Address, error) {
	return _SubsSenderInt.Contract.NebulaAddress(&_SubsSenderInt.CallOpts)
}

// NebulaAddress is a free data retrieval call binding the contract method 0x18f20d63.
//
// Solidity: function nebulaAddress() view returns(address)
func (_SubsSenderInt *SubsSenderIntCallerSession) NebulaAddress() (common.Address, error) {
	return _SubsSenderInt.Contract.NebulaAddress(&_SubsSenderInt.CallOpts)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0x1fd8ce7a.
//
// Solidity: function sendValueToSub(int64 value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderInt *SubsSenderIntTransactor) SendValueToSub(opts *bind.TransactOpts, value int64, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderInt.contract.Transact(opts, "sendValueToSub", value, blockNumber, subId)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0x1fd8ce7a.
//
// Solidity: function sendValueToSub(int64 value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderInt *SubsSenderIntSession) SendValueToSub(value int64, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.SendValueToSub(&_SubsSenderInt.TransactOpts, value, blockNumber, subId)
}

// SendValueToSub is a paid mutator transaction binding the contract method 0x1fd8ce7a.
//
// Solidity: function sendValueToSub(int64 value, uint256 blockNumber, bytes32 subId) returns()
func (_SubsSenderInt *SubsSenderIntTransactorSession) SendValueToSub(value int64, blockNumber *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.SendValueToSub(&_SubsSenderInt.TransactOpts, value, blockNumber, subId)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderInt *SubsSenderIntTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _SubsSenderInt.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderInt *SubsSenderIntSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.Fallback(&_SubsSenderInt.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SubsSenderInt *SubsSenderIntTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SubsSenderInt.Contract.Fallback(&_SubsSenderInt.TransactOpts, calldata)
}
