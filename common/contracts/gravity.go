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

// GravityABI is the input ABI used to generate the binding from.
const GravityABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newConsuls\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"newBftValue\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"bftValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"consuls\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rounds\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getConsuls\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newConsuls\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"newLastRound\",\"type\":\"uint256\"}],\"name\":\"updateConsuls\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newConsuls\",\"type\":\"address[]\"}],\"name\":\"hashNewConsuls\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// Gravity is an auto generated Go binding around an Ethereum contract.
type Gravity struct {
	GravityCaller     // Read-only binding to the contract
	GravityTransactor // Write-only binding to the contract
	GravityFilterer   // Log filterer for contract events
}

// GravityCaller is an auto generated read-only Go binding around an Ethereum contract.
type GravityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GravityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GravityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GravityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GravityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GravitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GravitySession struct {
	Contract     *Gravity          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GravityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GravityCallerSession struct {
	Contract *GravityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// GravityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GravityTransactorSession struct {
	Contract     *GravityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// GravityRaw is an auto generated low-level Go binding around an Ethereum contract.
type GravityRaw struct {
	Contract *Gravity // Generic contract binding to access the raw methods on
}

// GravityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GravityCallerRaw struct {
	Contract *GravityCaller // Generic read-only contract binding to access the raw methods on
}

// GravityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GravityTransactorRaw struct {
	Contract *GravityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGravity creates a new instance of Gravity, bound to a specific deployed contract.
func NewGravity(address common.Address, backend bind.ContractBackend) (*Gravity, error) {
	contract, err := bindGravity(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Gravity{GravityCaller: GravityCaller{contract: contract}, GravityTransactor: GravityTransactor{contract: contract}, GravityFilterer: GravityFilterer{contract: contract}}, nil
}

// NewGravityCaller creates a new read-only instance of Gravity, bound to a specific deployed contract.
func NewGravityCaller(address common.Address, caller bind.ContractCaller) (*GravityCaller, error) {
	contract, err := bindGravity(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GravityCaller{contract: contract}, nil
}

// NewGravityTransactor creates a new write-only instance of Gravity, bound to a specific deployed contract.
func NewGravityTransactor(address common.Address, transactor bind.ContractTransactor) (*GravityTransactor, error) {
	contract, err := bindGravity(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GravityTransactor{contract: contract}, nil
}

// NewGravityFilterer creates a new log filterer instance of Gravity, bound to a specific deployed contract.
func NewGravityFilterer(address common.Address, filterer bind.ContractFilterer) (*GravityFilterer, error) {
	contract, err := bindGravity(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GravityFilterer{contract: contract}, nil
}

// bindGravity binds a generic wrapper to an already deployed contract.
func bindGravity(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GravityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gravity *GravityRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Gravity.Contract.GravityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gravity *GravityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gravity.Contract.GravityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gravity *GravityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gravity.Contract.GravityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gravity *GravityCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Gravity.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gravity *GravityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gravity.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gravity *GravityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gravity.Contract.contract.Transact(opts, method, params...)
}

// BftValue is a free data retrieval call binding the contract method 0x3cec1bdd.
//
// Solidity: function bftValue() view returns(uint256)
func (_Gravity *GravityCaller) BftValue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Gravity.contract.Call(opts, out, "bftValue")
	return *ret0, err
}

// BftValue is a free data retrieval call binding the contract method 0x3cec1bdd.
//
// Solidity: function bftValue() view returns(uint256)
func (_Gravity *GravitySession) BftValue() (*big.Int, error) {
	return _Gravity.Contract.BftValue(&_Gravity.CallOpts)
}

// BftValue is a free data retrieval call binding the contract method 0x3cec1bdd.
//
// Solidity: function bftValue() view returns(uint256)
func (_Gravity *GravityCallerSession) BftValue() (*big.Int, error) {
	return _Gravity.Contract.BftValue(&_Gravity.CallOpts)
}

// Consuls is a free data retrieval call binding the contract method 0xa2c2c617.
//
// Solidity: function consuls(uint256 ) view returns(address)
func (_Gravity *GravityCaller) Consuls(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Gravity.contract.Call(opts, out, "consuls", arg0)
	return *ret0, err
}

// Consuls is a free data retrieval call binding the contract method 0xa2c2c617.
//
// Solidity: function consuls(uint256 ) view returns(address)
func (_Gravity *GravitySession) Consuls(arg0 *big.Int) (common.Address, error) {
	return _Gravity.Contract.Consuls(&_Gravity.CallOpts, arg0)
}

// Consuls is a free data retrieval call binding the contract method 0xa2c2c617.
//
// Solidity: function consuls(uint256 ) view returns(address)
func (_Gravity *GravityCallerSession) Consuls(arg0 *big.Int) (common.Address, error) {
	return _Gravity.Contract.Consuls(&_Gravity.CallOpts, arg0)
}

// GetConsuls is a free data retrieval call binding the contract method 0xad595b1a.
//
// Solidity: function getConsuls() view returns(address[])
func (_Gravity *GravityCaller) GetConsuls(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Gravity.contract.Call(opts, out, "getConsuls")
	return *ret0, err
}

// GetConsuls is a free data retrieval call binding the contract method 0xad595b1a.
//
// Solidity: function getConsuls() view returns(address[])
func (_Gravity *GravitySession) GetConsuls() ([]common.Address, error) {
	return _Gravity.Contract.GetConsuls(&_Gravity.CallOpts)
}

// GetConsuls is a free data retrieval call binding the contract method 0xad595b1a.
//
// Solidity: function getConsuls() view returns(address[])
func (_Gravity *GravityCallerSession) GetConsuls() ([]common.Address, error) {
	return _Gravity.Contract.GetConsuls(&_Gravity.CallOpts)
}

// HashNewConsuls is a free data retrieval call binding the contract method 0x4dea5eba.
//
// Solidity: function hashNewConsuls(address[] newConsuls) pure returns(bytes32)
func (_Gravity *GravityCaller) HashNewConsuls(opts *bind.CallOpts, newConsuls []common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Gravity.contract.Call(opts, out, "hashNewConsuls", newConsuls)
	return *ret0, err
}

// HashNewConsuls is a free data retrieval call binding the contract method 0x4dea5eba.
//
// Solidity: function hashNewConsuls(address[] newConsuls) pure returns(bytes32)
func (_Gravity *GravitySession) HashNewConsuls(newConsuls []common.Address) ([32]byte, error) {
	return _Gravity.Contract.HashNewConsuls(&_Gravity.CallOpts, newConsuls)
}

// HashNewConsuls is a free data retrieval call binding the contract method 0x4dea5eba.
//
// Solidity: function hashNewConsuls(address[] newConsuls) pure returns(bytes32)
func (_Gravity *GravityCallerSession) HashNewConsuls(newConsuls []common.Address) ([32]byte, error) {
	return _Gravity.Contract.HashNewConsuls(&_Gravity.CallOpts, newConsuls)
}

// Rounds is a free data retrieval call binding the contract method 0x8c65c81f.
//
// Solidity: function rounds(uint256 ) view returns(bool)
func (_Gravity *GravityCaller) Rounds(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Gravity.contract.Call(opts, out, "rounds", arg0)
	return *ret0, err
}

// Rounds is a free data retrieval call binding the contract method 0x8c65c81f.
//
// Solidity: function rounds(uint256 ) view returns(bool)
func (_Gravity *GravitySession) Rounds(arg0 *big.Int) (bool, error) {
	return _Gravity.Contract.Rounds(&_Gravity.CallOpts, arg0)
}

// Rounds is a free data retrieval call binding the contract method 0x8c65c81f.
//
// Solidity: function rounds(uint256 ) view returns(bool)
func (_Gravity *GravityCallerSession) Rounds(arg0 *big.Int) (bool, error) {
	return _Gravity.Contract.Rounds(&_Gravity.CallOpts, arg0)
}

// UpdateConsuls is a paid mutator transaction binding the contract method 0x92c388ab.
//
// Solidity: function updateConsuls(address[] newConsuls, uint8[] v, bytes32[] r, bytes32[] s, uint256 newLastRound) returns()
func (_Gravity *GravityTransactor) UpdateConsuls(opts *bind.TransactOpts, newConsuls []common.Address, v []uint8, r [][32]byte, s [][32]byte, newLastRound *big.Int) (*types.Transaction, error) {
	return _Gravity.contract.Transact(opts, "updateConsuls", newConsuls, v, r, s, newLastRound)
}

// UpdateConsuls is a paid mutator transaction binding the contract method 0x92c388ab.
//
// Solidity: function updateConsuls(address[] newConsuls, uint8[] v, bytes32[] r, bytes32[] s, uint256 newLastRound) returns()
func (_Gravity *GravitySession) UpdateConsuls(newConsuls []common.Address, v []uint8, r [][32]byte, s [][32]byte, newLastRound *big.Int) (*types.Transaction, error) {
	return _Gravity.Contract.UpdateConsuls(&_Gravity.TransactOpts, newConsuls, v, r, s, newLastRound)
}

// UpdateConsuls is a paid mutator transaction binding the contract method 0x92c388ab.
//
// Solidity: function updateConsuls(address[] newConsuls, uint8[] v, bytes32[] r, bytes32[] s, uint256 newLastRound) returns()
func (_Gravity *GravityTransactorSession) UpdateConsuls(newConsuls []common.Address, v []uint8, r [][32]byte, s [][32]byte, newLastRound *big.Int) (*types.Transaction, error) {
	return _Gravity.Contract.UpdateConsuls(&_Gravity.TransactOpts, newConsuls, v, r, s, newLastRound)
}
