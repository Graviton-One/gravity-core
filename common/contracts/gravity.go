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
const GravityABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newConsuls\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"oracleQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"newScores\",\"type\":\"uint256[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"updateScores\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"scores\",\"type\":\"uint256[]\"}],\"name\":\"hashScores\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

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

// HashScores is a free data retrieval call binding the contract method 0xd0b6c1d4.
//
// Solidity: function hashScores(address[] newOracles, uint256[] scores) pure returns(bytes32)
func (_Gravity *GravityCaller) HashScores(opts *bind.CallOpts, newOracles []common.Address, scores []*big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Gravity.contract.Call(opts, out, "hashScores", newOracles, scores)
	return *ret0, err
}

// HashScores is a free data retrieval call binding the contract method 0xd0b6c1d4.
//
// Solidity: function hashScores(address[] newOracles, uint256[] scores) pure returns(bytes32)
func (_Gravity *GravitySession) HashScores(newOracles []common.Address, scores []*big.Int) ([32]byte, error) {
	return _Gravity.Contract.HashScores(&_Gravity.CallOpts, newOracles, scores)
}

// HashScores is a free data retrieval call binding the contract method 0xd0b6c1d4.
//
// Solidity: function hashScores(address[] newOracles, uint256[] scores) pure returns(bytes32)
func (_Gravity *GravityCallerSession) HashScores(newOracles []common.Address, scores []*big.Int) ([32]byte, error) {
	return _Gravity.Contract.HashScores(&_Gravity.CallOpts, newOracles, scores)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Gravity *GravityCaller) OracleQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	ret := new(struct {
		First [32]byte
		Last  [32]byte
	})
	out := ret
	err := _Gravity.contract.Call(opts, out, "oracleQueue")
	return *ret, err
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Gravity *GravitySession) OracleQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Gravity.Contract.OracleQueue(&_Gravity.CallOpts)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Gravity *GravityCallerSession) OracleQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Gravity.Contract.OracleQueue(&_Gravity.CallOpts)
}

// UpdateScores is a paid mutator transaction binding the contract method 0x16b2d85e.
//
// Solidity: function updateScores(address[] newOracles, uint256[] newScores, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Gravity *GravityTransactor) UpdateScores(opts *bind.TransactOpts, newOracles []common.Address, newScores []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Gravity.contract.Transact(opts, "updateScores", newOracles, newScores, v, r, s)
}

// UpdateScores is a paid mutator transaction binding the contract method 0x16b2d85e.
//
// Solidity: function updateScores(address[] newOracles, uint256[] newScores, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Gravity *GravitySession) UpdateScores(newOracles []common.Address, newScores []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Gravity.Contract.UpdateScores(&_Gravity.TransactOpts, newOracles, newScores, v, r, s)
}

// UpdateScores is a paid mutator transaction binding the contract method 0x16b2d85e.
//
// Solidity: function updateScores(address[] newOracles, uint256[] newScores, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Gravity *GravityTransactorSession) UpdateScores(newOracles []common.Address, newScores []*big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Gravity.Contract.UpdateScores(&_Gravity.TransactOpts, newOracles, newScores, v, r, s)
}
