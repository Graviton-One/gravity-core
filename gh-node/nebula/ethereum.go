package nebula

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

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

// NebulaABI is the input ABI used to generate the binding from.
const NebulaABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newName\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"newOracle\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentOracleId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"endEpochHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isPublseSubSent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"oracleInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isOnline\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"idInQueue\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"oracleQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"oracles\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"oraclesCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pulseQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pulses\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"confirmationCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"subscriptions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"minConfirmations\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"subscriptionsQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"confirmData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subscriptionId\",\"type\":\"bytes32\"}],\"name\":\"sendData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"minConfirmations\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Nebula is an auto generated Go binding around an Ethereum contract.
type Nebula struct {
	NebulaCaller     // Read-only binding to the contract
	NebulaTransactor // Write-only binding to the contract
	NebulaFilterer   // Log filterer for contract events
}

// NebulaCaller is an auto generated read-only Go binding around an Ethereum contract.
type NebulaCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NebulaTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NebulaTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NebulaFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NebulaFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NebulaSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NebulaSession struct {
	Contract     *Nebula           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NebulaCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NebulaCallerSession struct {
	Contract *NebulaCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// NebulaTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NebulaTransactorSession struct {
	Contract     *NebulaTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NebulaRaw is an auto generated low-level Go binding around an Ethereum contract.
type NebulaRaw struct {
	Contract *Nebula // Generic contract binding to access the raw methods on
}

// NebulaCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NebulaCallerRaw struct {
	Contract *NebulaCaller // Generic read-only contract binding to access the raw methods on
}

// NebulaTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NebulaTransactorRaw struct {
	Contract *NebulaTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNebula creates a new instance of Nebula, bound to a specific deployed contract.
func NewNebula(address common.Address, backend bind.ContractBackend) (*Nebula, error) {
	contract, err := bindNebula(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Nebula{NebulaCaller: NebulaCaller{contract: contract}, NebulaTransactor: NebulaTransactor{contract: contract}, NebulaFilterer: NebulaFilterer{contract: contract}}, nil
}

// NewNebulaCaller creates a new read-only instance of Nebula, bound to a specific deployed contract.
func NewNebulaCaller(address common.Address, caller bind.ContractCaller) (*NebulaCaller, error) {
	contract, err := bindNebula(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NebulaCaller{contract: contract}, nil
}

// NewNebulaTransactor creates a new write-only instance of Nebula, bound to a specific deployed contract.
func NewNebulaTransactor(address common.Address, transactor bind.ContractTransactor) (*NebulaTransactor, error) {
	contract, err := bindNebula(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NebulaTransactor{contract: contract}, nil
}

// NewNebulaFilterer creates a new log filterer instance of Nebula, bound to a specific deployed contract.
func NewNebulaFilterer(address common.Address, filterer bind.ContractFilterer) (*NebulaFilterer, error) {
	contract, err := bindNebula(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NebulaFilterer{contract: contract}, nil
}

// bindNebula binds a generic wrapper to an already deployed contract.
func bindNebula(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NebulaABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nebula *NebulaRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Nebula.Contract.NebulaCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nebula *NebulaRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nebula.Contract.NebulaTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nebula *NebulaRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nebula.Contract.NebulaTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nebula *NebulaCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Nebula.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nebula *NebulaTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nebula.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nebula *NebulaTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nebula.Contract.contract.Transact(opts, method, params...)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Nebula *NebulaCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "currentEpoch")
	return *ret0, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Nebula *NebulaSession) CurrentEpoch() (*big.Int, error) {
	return _Nebula.Contract.CurrentEpoch(&_Nebula.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Nebula *NebulaCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Nebula.Contract.CurrentEpoch(&_Nebula.CallOpts)
}

// CurrentOracleId is a free data retrieval call binding the contract method 0x524ce31f.
//
// Solidity: function currentOracleId() view returns(bytes32)
func (_Nebula *NebulaCaller) CurrentOracleId(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "currentOracleId")
	return *ret0, err
}

// CurrentOracleId is a free data retrieval call binding the contract method 0x524ce31f.
//
// Solidity: function currentOracleId() view returns(bytes32)
func (_Nebula *NebulaSession) CurrentOracleId() ([32]byte, error) {
	return _Nebula.Contract.CurrentOracleId(&_Nebula.CallOpts)
}

// CurrentOracleId is a free data retrieval call binding the contract method 0x524ce31f.
//
// Solidity: function currentOracleId() view returns(bytes32)
func (_Nebula *NebulaCallerSession) CurrentOracleId() ([32]byte, error) {
	return _Nebula.Contract.CurrentOracleId(&_Nebula.CallOpts)
}

// EndEpochHeight is a free data retrieval call binding the contract method 0x996211a5.
//
// Solidity: function endEpochHeight() view returns(uint256)
func (_Nebula *NebulaCaller) EndEpochHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "endEpochHeight")
	return *ret0, err
}

// EndEpochHeight is a free data retrieval call binding the contract method 0x996211a5.
//
// Solidity: function endEpochHeight() view returns(uint256)
func (_Nebula *NebulaSession) EndEpochHeight() (*big.Int, error) {
	return _Nebula.Contract.EndEpochHeight(&_Nebula.CallOpts)
}

// EndEpochHeight is a free data retrieval call binding the contract method 0x996211a5.
//
// Solidity: function endEpochHeight() view returns(uint256)
func (_Nebula *NebulaCallerSession) EndEpochHeight() (*big.Int, error) {
	return _Nebula.Contract.EndEpochHeight(&_Nebula.CallOpts)
}

// IsPublseSubSent is a free data retrieval call binding the contract method 0x6148d3f3.
//
// Solidity: function isPublseSubSent(uint256 , bytes32 ) view returns(bool)
func (_Nebula *NebulaCaller) IsPublseSubSent(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "isPublseSubSent", arg0, arg1)
	return *ret0, err
}

// IsPublseSubSent is a free data retrieval call binding the contract method 0x6148d3f3.
//
// Solidity: function isPublseSubSent(uint256 , bytes32 ) view returns(bool)
func (_Nebula *NebulaSession) IsPublseSubSent(arg0 *big.Int, arg1 [32]byte) (bool, error) {
	return _Nebula.Contract.IsPublseSubSent(&_Nebula.CallOpts, arg0, arg1)
}

// IsPublseSubSent is a free data retrieval call binding the contract method 0x6148d3f3.
//
// Solidity: function isPublseSubSent(uint256 , bytes32 ) view returns(bool)
func (_Nebula *NebulaCallerSession) IsPublseSubSent(arg0 *big.Int, arg1 [32]byte) (bool, error) {
	return _Nebula.Contract.IsPublseSubSent(&_Nebula.CallOpts, arg0, arg1)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nebula *NebulaCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nebula *NebulaSession) Name() (string, error) {
	return _Nebula.Contract.Name(&_Nebula.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nebula *NebulaCallerSession) Name() (string, error) {
	return _Nebula.Contract.Name(&_Nebula.CallOpts)
}

// OracleInfo is a free data retrieval call binding the contract method 0xbc623feb.
//
// Solidity: function oracleInfo(address ) view returns(address owner, bool isOnline, bytes32 idInQueue)
func (_Nebula *NebulaCaller) OracleInfo(opts *bind.CallOpts, arg0 common.Address) (struct {
	Owner     common.Address
	IsOnline  bool
	IdInQueue [32]byte
}, error) {
	ret := new(struct {
		Owner     common.Address
		IsOnline  bool
		IdInQueue [32]byte
	})
	out := ret
	err := _Nebula.contract.Call(opts, out, "oracleInfo", arg0)
	return *ret, err
}

// OracleInfo is a free data retrieval call binding the contract method 0xbc623feb.
//
// Solidity: function oracleInfo(address ) view returns(address owner, bool isOnline, bytes32 idInQueue)
func (_Nebula *NebulaSession) OracleInfo(arg0 common.Address) (struct {
	Owner     common.Address
	IsOnline  bool
	IdInQueue [32]byte
}, error) {
	return _Nebula.Contract.OracleInfo(&_Nebula.CallOpts, arg0)
}

// OracleInfo is a free data retrieval call binding the contract method 0xbc623feb.
//
// Solidity: function oracleInfo(address ) view returns(address owner, bool isOnline, bytes32 idInQueue)
func (_Nebula *NebulaCallerSession) OracleInfo(arg0 common.Address) (struct {
	Owner     common.Address
	IsOnline  bool
	IdInQueue [32]byte
}, error) {
	return _Nebula.Contract.OracleInfo(&_Nebula.CallOpts, arg0)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCaller) OracleQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	ret := new(struct {
		First [32]byte
		Last  [32]byte
	})
	out := ret
	err := _Nebula.contract.Call(opts, out, "oracleQueue")
	return *ret, err
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaSession) OracleQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Nebula.Contract.OracleQueue(&_Nebula.CallOpts)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCallerSession) OracleQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Nebula.Contract.OracleQueue(&_Nebula.CallOpts)
}

// Oracles is a free data retrieval call binding the contract method 0xa81a2677.
//
// Solidity: function oracles(bytes32 ) view returns(address)
func (_Nebula *NebulaCaller) Oracles(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "oracles", arg0)
	return *ret0, err
}

// Oracles is a free data retrieval call binding the contract method 0xa81a2677.
//
// Solidity: function oracles(bytes32 ) view returns(address)
func (_Nebula *NebulaSession) Oracles(arg0 [32]byte) (common.Address, error) {
	return _Nebula.Contract.Oracles(&_Nebula.CallOpts, arg0)
}

// Oracles is a free data retrieval call binding the contract method 0xa81a2677.
//
// Solidity: function oracles(bytes32 ) view returns(address)
func (_Nebula *NebulaCallerSession) Oracles(arg0 [32]byte) (common.Address, error) {
	return _Nebula.Contract.Oracles(&_Nebula.CallOpts, arg0)
}

// OraclesCount is a free data retrieval call binding the contract method 0x76f70900.
//
// Solidity: function oraclesCount() view returns(uint256)
func (_Nebula *NebulaCaller) OraclesCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Nebula.contract.Call(opts, out, "oraclesCount")
	return *ret0, err
}

// OraclesCount is a free data retrieval call binding the contract method 0x76f70900.
//
// Solidity: function oraclesCount() view returns(uint256)
func (_Nebula *NebulaSession) OraclesCount() (*big.Int, error) {
	return _Nebula.Contract.OraclesCount(&_Nebula.CallOpts)
}

// OraclesCount is a free data retrieval call binding the contract method 0x76f70900.
//
// Solidity: function oraclesCount() view returns(uint256)
func (_Nebula *NebulaCallerSession) OraclesCount() (*big.Int, error) {
	return _Nebula.Contract.OraclesCount(&_Nebula.CallOpts)
}

// PulseQueue is a free data retrieval call binding the contract method 0x1d11f944.
//
// Solidity: function pulseQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCaller) PulseQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	ret := new(struct {
		First [32]byte
		Last  [32]byte
	})
	out := ret
	err := _Nebula.contract.Call(opts, out, "pulseQueue")
	return *ret, err
}

// PulseQueue is a free data retrieval call binding the contract method 0x1d11f944.
//
// Solidity: function pulseQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaSession) PulseQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Nebula.Contract.PulseQueue(&_Nebula.CallOpts)
}

// PulseQueue is a free data retrieval call binding the contract method 0x1d11f944.
//
// Solidity: function pulseQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCallerSession) PulseQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Nebula.Contract.PulseQueue(&_Nebula.CallOpts)
}

// Pulses is a free data retrieval call binding the contract method 0x0694fbb3.
//
// Solidity: function pulses(uint256 ) view returns(bytes32 dataHash, uint256 confirmationCount)
func (_Nebula *NebulaCaller) Pulses(opts *bind.CallOpts, arg0 *big.Int) (struct {
	DataHash          [32]byte
	ConfirmationCount *big.Int
}, error) {
	ret := new(struct {
		DataHash          [32]byte
		ConfirmationCount *big.Int
	})
	out := ret
	err := _Nebula.contract.Call(opts, out, "pulses", arg0)
	return *ret, err
}

// Pulses is a free data retrieval call binding the contract method 0x0694fbb3.
//
// Solidity: function pulses(uint256 ) view returns(bytes32 dataHash, uint256 confirmationCount)
func (_Nebula *NebulaSession) Pulses(arg0 *big.Int) (struct {
	DataHash          [32]byte
	ConfirmationCount *big.Int
}, error) {
	return _Nebula.Contract.Pulses(&_Nebula.CallOpts, arg0)
}

// Pulses is a free data retrieval call binding the contract method 0x0694fbb3.
//
// Solidity: function pulses(uint256 ) view returns(bytes32 dataHash, uint256 confirmationCount)
func (_Nebula *NebulaCallerSession) Pulses(arg0 *big.Int) (struct {
	DataHash          [32]byte
	ConfirmationCount *big.Int
}, error) {
	return _Nebula.Contract.Pulses(&_Nebula.CallOpts, arg0)
}

// Subscriptions is a free data retrieval call binding the contract method 0x94259c6c.
//
// Solidity: function subscriptions(bytes32 ) view returns(address owner, address contractAddress, uint8 minConfirmations, uint256 reward)
func (_Nebula *NebulaCaller) Subscriptions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Owner            common.Address
	ContractAddress  common.Address
	MinConfirmations uint8
	Reward           *big.Int
}, error) {
	ret := new(struct {
		Owner            common.Address
		ContractAddress  common.Address
		MinConfirmations uint8
		Reward           *big.Int
	})
	out := ret
	err := _Nebula.contract.Call(opts, out, "subscriptions", arg0)
	return *ret, err
}

// Subscriptions is a free data retrieval call binding the contract method 0x94259c6c.
//
// Solidity: function subscriptions(bytes32 ) view returns(address owner, address contractAddress, uint8 minConfirmations, uint256 reward)
func (_Nebula *NebulaSession) Subscriptions(arg0 [32]byte) (struct {
	Owner            common.Address
	ContractAddress  common.Address
	MinConfirmations uint8
	Reward           *big.Int
}, error) {
	return _Nebula.Contract.Subscriptions(&_Nebula.CallOpts, arg0)
}

// Subscriptions is a free data retrieval call binding the contract method 0x94259c6c.
//
// Solidity: function subscriptions(bytes32 ) view returns(address owner, address contractAddress, uint8 minConfirmations, uint256 reward)
func (_Nebula *NebulaCallerSession) Subscriptions(arg0 [32]byte) (struct {
	Owner            common.Address
	ContractAddress  common.Address
	MinConfirmations uint8
	Reward           *big.Int
}, error) {
	return _Nebula.Contract.Subscriptions(&_Nebula.CallOpts, arg0)
}

// SubscriptionsQueue is a free data retrieval call binding the contract method 0xb48a9c9b.
//
// Solidity: function subscriptionsQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCaller) SubscriptionsQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	ret := new(struct {
		First [32]byte
		Last  [32]byte
	})
	out := ret
	err := _Nebula.contract.Call(opts, out, "subscriptionsQueue")
	return *ret, err
}

// SubscriptionsQueue is a free data retrieval call binding the contract method 0xb48a9c9b.
//
// Solidity: function subscriptionsQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaSession) SubscriptionsQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Nebula.Contract.SubscriptionsQueue(&_Nebula.CallOpts)
}

// SubscriptionsQueue is a free data retrieval call binding the contract method 0xb48a9c9b.
//
// Solidity: function subscriptionsQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCallerSession) SubscriptionsQueue() (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	return _Nebula.Contract.SubscriptionsQueue(&_Nebula.CallOpts)
}

// ConfirmData is a paid mutator transaction binding the contract method 0x60be26d9.
//
// Solidity: function confirmData(bytes32 dataHash, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Nebula *NebulaTransactor) ConfirmData(opts *bind.TransactOpts, dataHash [32]byte, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "confirmData", dataHash, v, r, s)
}

// ConfirmData is a paid mutator transaction binding the contract method 0x60be26d9.
//
// Solidity: function confirmData(bytes32 dataHash, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Nebula *NebulaSession) ConfirmData(dataHash [32]byte, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.ConfirmData(&_Nebula.TransactOpts, dataHash, v, r, s)
}

// ConfirmData is a paid mutator transaction binding the contract method 0x60be26d9.
//
// Solidity: function confirmData(bytes32 dataHash, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Nebula *NebulaTransactorSession) ConfirmData(dataHash [32]byte, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.ConfirmData(&_Nebula.TransactOpts, dataHash, v, r, s)
}

// SendData is a paid mutator transaction binding the contract method 0x04dd915d.
//
// Solidity: function sendData(uint256 value, uint256 blockNumber, bytes32 subscriptionId) returns()
func (_Nebula *NebulaTransactor) SendData(opts *bind.TransactOpts, value *big.Int, blockNumber *big.Int, subscriptionId [32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "sendData", value, blockNumber, subscriptionId)
}

// SendData is a paid mutator transaction binding the contract method 0x04dd915d.
//
// Solidity: function sendData(uint256 value, uint256 blockNumber, bytes32 subscriptionId) returns()
func (_Nebula *NebulaSession) SendData(value *big.Int, blockNumber *big.Int, subscriptionId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendData(&_Nebula.TransactOpts, value, blockNumber, subscriptionId)
}

// SendData is a paid mutator transaction binding the contract method 0x04dd915d.
//
// Solidity: function sendData(uint256 value, uint256 blockNumber, bytes32 subscriptionId) returns()
func (_Nebula *NebulaTransactorSession) SendData(value *big.Int, blockNumber *big.Int, subscriptionId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendData(&_Nebula.TransactOpts, value, blockNumber, subscriptionId)
}

// Subscribe is a paid mutator transaction binding the contract method 0x3527715d.
//
// Solidity: function subscribe(address contractAddress, uint8 minConfirmations, uint256 reward) returns()
func (_Nebula *NebulaTransactor) Subscribe(opts *bind.TransactOpts, contractAddress common.Address, minConfirmations uint8, reward *big.Int) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "subscribe", contractAddress, minConfirmations, reward)
}

// Subscribe is a paid mutator transaction binding the contract method 0x3527715d.
//
// Solidity: function subscribe(address contractAddress, uint8 minConfirmations, uint256 reward) returns()
func (_Nebula *NebulaSession) Subscribe(contractAddress common.Address, minConfirmations uint8, reward *big.Int) (*types.Transaction, error) {
	return _Nebula.Contract.Subscribe(&_Nebula.TransactOpts, contractAddress, minConfirmations, reward)
}

// Subscribe is a paid mutator transaction binding the contract method 0x3527715d.
//
// Solidity: function subscribe(address contractAddress, uint8 minConfirmations, uint256 reward) returns()
func (_Nebula *NebulaTransactorSession) Subscribe(contractAddress common.Address, minConfirmations uint8, reward *big.Int) (*types.Transaction, error) {
	return _Nebula.Contract.Subscribe(&_Nebula.TransactOpts, contractAddress, minConfirmations, reward)
}

// Unsubscribe is a paid mutator transaction binding the contract method 0xf8c60146.
//
// Solidity: function unsubscribe(bytes32 id) returns()
func (_Nebula *NebulaTransactor) Unsubscribe(opts *bind.TransactOpts, id [32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "unsubscribe", id)
}

// Unsubscribe is a paid mutator transaction binding the contract method 0xf8c60146.
//
// Solidity: function unsubscribe(bytes32 id) returns()
func (_Nebula *NebulaSession) Unsubscribe(id [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.Unsubscribe(&_Nebula.TransactOpts, id)
}

// Unsubscribe is a paid mutator transaction binding the contract method 0xf8c60146.
//
// Solidity: function unsubscribe(bytes32 id) returns()
func (_Nebula *NebulaTransactorSession) Unsubscribe(id [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.Unsubscribe(&_Nebula.TransactOpts, id)
}
