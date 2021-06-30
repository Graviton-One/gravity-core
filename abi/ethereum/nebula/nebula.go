// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nebula

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
const GravityABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"consuls\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"newBftValue\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"bftValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConsuls\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getConsulsByRoundId\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newConsuls\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"hashNewConsuls\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rounds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newConsuls\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"updateConsuls\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GravityFuncSigs maps the 4-byte function signature to its string representation.
var GravityFuncSigs = map[string]string{
	"3cec1bdd": "bftValue()",
	"ad595b1a": "getConsuls()",
	"fe7378bb": "getConsulsByRoundId(uint256)",
	"c85f8d33": "hashNewConsuls(address[],uint256)",
	"82bc07e6": "lastRound()",
	"e6da9213": "rounds(uint256,uint256)",
	"92c388ab": "updateConsuls(address[],uint8[],bytes32[],bytes32[],uint256)",
}

// GravityBin is the compiled bytecode used for deploying new contracts.
var GravityBin = "0x608060405234801561001057600080fd5b50604051610ae3380380610ae38339818101604052604081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825186602082028301116401000000008211171561008557600080fd5b82525081516020918201928201910280838360005b838110156100b257818101518382015260200161009a565b5050505091909101604052506020908101516000808052825284519093506100ff92507fad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb591850190610109565b506001555061018d565b82805482825590600052602060002090810192821561015e579160200282015b8281111561015e57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190610129565b5061016a92915061016e565b5090565b5b8082111561016a5780546001600160a01b031916815560010161016f565b6109478061019c6000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063ad595b1a1161005b578063ad595b1a146102cf578063c85f8d3314610327578063e6da9213146103ca578063fe7378bb146104095761007d565b80633cec1bdd1461008257806382bc07e61461009c57806392c388ab146100a4575b600080fd5b61008a610426565b60408051918252519081900360200190f35b61008a61042c565b6102cd600480360360a08110156100ba57600080fd5b810190602081018135600160201b8111156100d457600080fd5b8201836020820111156100e657600080fd5b803590602001918460208302840111600160201b8311171561010757600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561015657600080fd5b82018360208201111561016857600080fd5b803590602001918460208302840111600160201b8311171561018957600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156101d857600080fd5b8201836020820111156101ea57600080fd5b803590602001918460208302840111600160201b8311171561020b57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561025a57600080fd5b82018360208201111561026c57600080fd5b803590602001918460208302840111600160201b8311171561028d57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295505091359250610432915050565b005b6102d7610661565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156103135781810151838201526020016102fb565b505050509050019250505060405180910390f35b61008a6004803603604081101561033d57600080fd5b810190602081018135600160201b81111561035757600080fd5b82018360208201111561036957600080fd5b803590602001918460208302840111600160201b8311171561038a57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092955050913592506106cc915050565b6103ed600480360360408110156103e057600080fd5b50803590602001356107ee565b604080516001600160a01b039092168252519081900360200190f35b6102d76004803603602081101561041f57600080fd5b5035610823565b60015481565b60025481565b60006002548211610482576040805162461bcd60e51b81526020600482015260156024820152741c9bdd5b99081b195cdcc81b185cdd081c9bdd5b99605a1b604482015290519081900360640190fd5b600061048e87846106cc565b60025460009081526020818152604091829020805483518184028101840190945280845293945060609390918301828280156104f357602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116104d5575b5050505050905060005b81518110156105e85781818151811061051257fe5b60200260200101516001600160a01b03166001848a848151811061053257fe5b60200260200101518a858151811061054657fe5b60200260200101518a868151811061055a57fe5b602002602001015160405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156105b6573d6000803e3d6000fd5b505050602060405103516001600160a01b0316146105d55760006105d8565b60015b60ff1693909301926001016104fd565b50600154831015610634576040805162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a590818999d0818dbdd5b9d607a1b604482015290519081900360640190fd5b6000848152602081815260409091208951610651928b019061088d565b5050506002919091555050505050565b600254600090815260208181526040918290208054835181840281018401909452808452606093928301828280156106c257602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116106a4575b5050505050905090565b6000606060005b845181101561077957818582815181106106e957fe5b60200260200101516040516020018083805190602001908083835b602083106107235780518252601f199092019160209182019101610704565b6001836020036101000a038019825116818451168082178552505050505050905001826001600160a01b031660601b815260140192505050604051602081830303815290604052915080806001019150506106d3565b5080836040516020018083805190602001908083835b602083106107ae5780518252601f19909201916020918201910161078f565b51815160209384036101000a6000190180199092169116179052920193845250604080518085038152938201905282519201919091209695505050505050565b6000602052816000526040600020818154811061080757fe5b6000918252602090912001546001600160a01b03169150829050565b6000818152602081815260409182902080548351818402810184019094528084526060939283018282801561088157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610863575b50505050509050919050565b8280548282559060005260206000209081019282156108e2579160200282015b828111156108e257825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906108ad565b506108ee9291506108f2565b5090565b5b808211156108ee5780546001600160a01b03191681556001016108f356fea26469706673582212203484156838341c4d242265bccfaeab7043b9a6614e2eb9051f76e6983a1059d964736f6c63430007000033"

// DeployGravity deploys a new Ethereum contract, binding an instance of Gravity to it.
func DeployGravity(auth *bind.TransactOpts, backend bind.ContractBackend, consuls []common.Address, newBftValue *big.Int) (common.Address, *types.Transaction, *Gravity, error) {
	parsed, err := abi.JSON(strings.NewReader(GravityABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GravityBin), backend, consuls, newBftValue)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Gravity{GravityCaller: GravityCaller{contract: contract}, GravityTransactor: GravityTransactor{contract: contract}, GravityFilterer: GravityFilterer{contract: contract}}, nil
}

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
func (_Gravity *GravityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_Gravity *GravityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
	var out []interface{}
	err := _Gravity.contract.Call(opts, &out, "bftValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

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

// GetConsuls is a free data retrieval call binding the contract method 0xad595b1a.
//
// Solidity: function getConsuls() view returns(address[])
func (_Gravity *GravityCaller) GetConsuls(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Gravity.contract.Call(opts, &out, "getConsuls")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

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

// GetConsulsByRoundId is a free data retrieval call binding the contract method 0xfe7378bb.
//
// Solidity: function getConsulsByRoundId(uint256 roundId) view returns(address[])
func (_Gravity *GravityCaller) GetConsulsByRoundId(opts *bind.CallOpts, roundId *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _Gravity.contract.Call(opts, &out, "getConsulsByRoundId", roundId)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetConsulsByRoundId is a free data retrieval call binding the contract method 0xfe7378bb.
//
// Solidity: function getConsulsByRoundId(uint256 roundId) view returns(address[])
func (_Gravity *GravitySession) GetConsulsByRoundId(roundId *big.Int) ([]common.Address, error) {
	return _Gravity.Contract.GetConsulsByRoundId(&_Gravity.CallOpts, roundId)
}

// GetConsulsByRoundId is a free data retrieval call binding the contract method 0xfe7378bb.
//
// Solidity: function getConsulsByRoundId(uint256 roundId) view returns(address[])
func (_Gravity *GravityCallerSession) GetConsulsByRoundId(roundId *big.Int) ([]common.Address, error) {
	return _Gravity.Contract.GetConsulsByRoundId(&_Gravity.CallOpts, roundId)
}

// HashNewConsuls is a free data retrieval call binding the contract method 0xc85f8d33.
//
// Solidity: function hashNewConsuls(address[] newConsuls, uint256 roundId) pure returns(bytes32)
func (_Gravity *GravityCaller) HashNewConsuls(opts *bind.CallOpts, newConsuls []common.Address, roundId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Gravity.contract.Call(opts, &out, "hashNewConsuls", newConsuls, roundId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashNewConsuls is a free data retrieval call binding the contract method 0xc85f8d33.
//
// Solidity: function hashNewConsuls(address[] newConsuls, uint256 roundId) pure returns(bytes32)
func (_Gravity *GravitySession) HashNewConsuls(newConsuls []common.Address, roundId *big.Int) ([32]byte, error) {
	return _Gravity.Contract.HashNewConsuls(&_Gravity.CallOpts, newConsuls, roundId)
}

// HashNewConsuls is a free data retrieval call binding the contract method 0xc85f8d33.
//
// Solidity: function hashNewConsuls(address[] newConsuls, uint256 roundId) pure returns(bytes32)
func (_Gravity *GravityCallerSession) HashNewConsuls(newConsuls []common.Address, roundId *big.Int) ([32]byte, error) {
	return _Gravity.Contract.HashNewConsuls(&_Gravity.CallOpts, newConsuls, roundId)
}

// LastRound is a free data retrieval call binding the contract method 0x82bc07e6.
//
// Solidity: function lastRound() view returns(uint256)
func (_Gravity *GravityCaller) LastRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Gravity.contract.Call(opts, &out, "lastRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastRound is a free data retrieval call binding the contract method 0x82bc07e6.
//
// Solidity: function lastRound() view returns(uint256)
func (_Gravity *GravitySession) LastRound() (*big.Int, error) {
	return _Gravity.Contract.LastRound(&_Gravity.CallOpts)
}

// LastRound is a free data retrieval call binding the contract method 0x82bc07e6.
//
// Solidity: function lastRound() view returns(uint256)
func (_Gravity *GravityCallerSession) LastRound() (*big.Int, error) {
	return _Gravity.Contract.LastRound(&_Gravity.CallOpts)
}

// Rounds is a free data retrieval call binding the contract method 0xe6da9213.
//
// Solidity: function rounds(uint256 , uint256 ) view returns(address)
func (_Gravity *GravityCaller) Rounds(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Gravity.contract.Call(opts, &out, "rounds", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Rounds is a free data retrieval call binding the contract method 0xe6da9213.
//
// Solidity: function rounds(uint256 , uint256 ) view returns(address)
func (_Gravity *GravitySession) Rounds(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _Gravity.Contract.Rounds(&_Gravity.CallOpts, arg0, arg1)
}

// Rounds is a free data retrieval call binding the contract method 0xe6da9213.
//
// Solidity: function rounds(uint256 , uint256 ) view returns(address)
func (_Gravity *GravityCallerSession) Rounds(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _Gravity.Contract.Rounds(&_Gravity.CallOpts, arg0, arg1)
}

// UpdateConsuls is a paid mutator transaction binding the contract method 0x92c388ab.
//
// Solidity: function updateConsuls(address[] newConsuls, uint8[] v, bytes32[] r, bytes32[] s, uint256 roundId) returns()
func (_Gravity *GravityTransactor) UpdateConsuls(opts *bind.TransactOpts, newConsuls []common.Address, v []uint8, r [][32]byte, s [][32]byte, roundId *big.Int) (*types.Transaction, error) {
	return _Gravity.contract.Transact(opts, "updateConsuls", newConsuls, v, r, s, roundId)
}

// UpdateConsuls is a paid mutator transaction binding the contract method 0x92c388ab.
//
// Solidity: function updateConsuls(address[] newConsuls, uint8[] v, bytes32[] r, bytes32[] s, uint256 roundId) returns()
func (_Gravity *GravitySession) UpdateConsuls(newConsuls []common.Address, v []uint8, r [][32]byte, s [][32]byte, roundId *big.Int) (*types.Transaction, error) {
	return _Gravity.Contract.UpdateConsuls(&_Gravity.TransactOpts, newConsuls, v, r, s, roundId)
}

// UpdateConsuls is a paid mutator transaction binding the contract method 0x92c388ab.
//
// Solidity: function updateConsuls(address[] newConsuls, uint8[] v, bytes32[] r, bytes32[] s, uint256 roundId) returns()
func (_Gravity *GravityTransactorSession) UpdateConsuls(newConsuls []common.Address, v []uint8, r [][32]byte, s [][32]byte, roundId *big.Int) (*types.Transaction, error) {
	return _Gravity.Contract.UpdateConsuls(&_Gravity.TransactOpts, newConsuls, v, r, s, roundId)
}

// ISubscriberBytesABI is the input ABI used to generate the binding from.
const ISubscriberBytesABI = "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"attachValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ISubscriberBytesFuncSigs maps the 4-byte function signature to its string representation.
var ISubscriberBytesFuncSigs = map[string]string{
	"cc32a151": "attachValue(bytes)",
}

// ISubscriberBytes is an auto generated Go binding around an Ethereum contract.
type ISubscriberBytes struct {
	ISubscriberBytesCaller     // Read-only binding to the contract
	ISubscriberBytesTransactor // Write-only binding to the contract
	ISubscriberBytesFilterer   // Log filterer for contract events
}

// ISubscriberBytesCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISubscriberBytesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberBytesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISubscriberBytesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberBytesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISubscriberBytesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberBytesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISubscriberBytesSession struct {
	Contract     *ISubscriberBytes // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISubscriberBytesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISubscriberBytesCallerSession struct {
	Contract *ISubscriberBytesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ISubscriberBytesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISubscriberBytesTransactorSession struct {
	Contract     *ISubscriberBytesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ISubscriberBytesRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISubscriberBytesRaw struct {
	Contract *ISubscriberBytes // Generic contract binding to access the raw methods on
}

// ISubscriberBytesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISubscriberBytesCallerRaw struct {
	Contract *ISubscriberBytesCaller // Generic read-only contract binding to access the raw methods on
}

// ISubscriberBytesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISubscriberBytesTransactorRaw struct {
	Contract *ISubscriberBytesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISubscriberBytes creates a new instance of ISubscriberBytes, bound to a specific deployed contract.
func NewISubscriberBytes(address common.Address, backend bind.ContractBackend) (*ISubscriberBytes, error) {
	contract, err := bindISubscriberBytes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISubscriberBytes{ISubscriberBytesCaller: ISubscriberBytesCaller{contract: contract}, ISubscriberBytesTransactor: ISubscriberBytesTransactor{contract: contract}, ISubscriberBytesFilterer: ISubscriberBytesFilterer{contract: contract}}, nil
}

// NewISubscriberBytesCaller creates a new read-only instance of ISubscriberBytes, bound to a specific deployed contract.
func NewISubscriberBytesCaller(address common.Address, caller bind.ContractCaller) (*ISubscriberBytesCaller, error) {
	contract, err := bindISubscriberBytes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISubscriberBytesCaller{contract: contract}, nil
}

// NewISubscriberBytesTransactor creates a new write-only instance of ISubscriberBytes, bound to a specific deployed contract.
func NewISubscriberBytesTransactor(address common.Address, transactor bind.ContractTransactor) (*ISubscriberBytesTransactor, error) {
	contract, err := bindISubscriberBytes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISubscriberBytesTransactor{contract: contract}, nil
}

// NewISubscriberBytesFilterer creates a new log filterer instance of ISubscriberBytes, bound to a specific deployed contract.
func NewISubscriberBytesFilterer(address common.Address, filterer bind.ContractFilterer) (*ISubscriberBytesFilterer, error) {
	contract, err := bindISubscriberBytes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISubscriberBytesFilterer{contract: contract}, nil
}

// bindISubscriberBytes binds a generic wrapper to an already deployed contract.
func bindISubscriberBytes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISubscriberBytesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISubscriberBytes *ISubscriberBytesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISubscriberBytes.Contract.ISubscriberBytesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISubscriberBytes *ISubscriberBytesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISubscriberBytes.Contract.ISubscriberBytesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISubscriberBytes *ISubscriberBytesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISubscriberBytes.Contract.ISubscriberBytesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISubscriberBytes *ISubscriberBytesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISubscriberBytes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISubscriberBytes *ISubscriberBytesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISubscriberBytes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISubscriberBytes *ISubscriberBytesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISubscriberBytes.Contract.contract.Transact(opts, method, params...)
}

// AttachValue is a paid mutator transaction binding the contract method 0xcc32a151.
//
// Solidity: function attachValue(bytes value) returns()
func (_ISubscriberBytes *ISubscriberBytesTransactor) AttachValue(opts *bind.TransactOpts, value []byte) (*types.Transaction, error) {
	return _ISubscriberBytes.contract.Transact(opts, "attachValue", value)
}

// AttachValue is a paid mutator transaction binding the contract method 0xcc32a151.
//
// Solidity: function attachValue(bytes value) returns()
func (_ISubscriberBytes *ISubscriberBytesSession) AttachValue(value []byte) (*types.Transaction, error) {
	return _ISubscriberBytes.Contract.AttachValue(&_ISubscriberBytes.TransactOpts, value)
}

// AttachValue is a paid mutator transaction binding the contract method 0xcc32a151.
//
// Solidity: function attachValue(bytes value) returns()
func (_ISubscriberBytes *ISubscriberBytesTransactorSession) AttachValue(value []byte) (*types.Transaction, error) {
	return _ISubscriberBytes.Contract.AttachValue(&_ISubscriberBytes.TransactOpts, value)
}

// ISubscriberIntABI is the input ABI used to generate the binding from.
const ISubscriberIntABI = "[{\"inputs\":[{\"internalType\":\"int64\",\"name\":\"value\",\"type\":\"int64\"}],\"name\":\"attachValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ISubscriberIntFuncSigs maps the 4-byte function signature to its string representation.
var ISubscriberIntFuncSigs = map[string]string{
	"7dc7c1b0": "attachValue(int64)",
}

// ISubscriberInt is an auto generated Go binding around an Ethereum contract.
type ISubscriberInt struct {
	ISubscriberIntCaller     // Read-only binding to the contract
	ISubscriberIntTransactor // Write-only binding to the contract
	ISubscriberIntFilterer   // Log filterer for contract events
}

// ISubscriberIntCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISubscriberIntCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberIntTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISubscriberIntTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberIntFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISubscriberIntFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberIntSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISubscriberIntSession struct {
	Contract     *ISubscriberInt   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISubscriberIntCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISubscriberIntCallerSession struct {
	Contract *ISubscriberIntCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ISubscriberIntTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISubscriberIntTransactorSession struct {
	Contract     *ISubscriberIntTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ISubscriberIntRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISubscriberIntRaw struct {
	Contract *ISubscriberInt // Generic contract binding to access the raw methods on
}

// ISubscriberIntCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISubscriberIntCallerRaw struct {
	Contract *ISubscriberIntCaller // Generic read-only contract binding to access the raw methods on
}

// ISubscriberIntTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISubscriberIntTransactorRaw struct {
	Contract *ISubscriberIntTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISubscriberInt creates a new instance of ISubscriberInt, bound to a specific deployed contract.
func NewISubscriberInt(address common.Address, backend bind.ContractBackend) (*ISubscriberInt, error) {
	contract, err := bindISubscriberInt(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISubscriberInt{ISubscriberIntCaller: ISubscriberIntCaller{contract: contract}, ISubscriberIntTransactor: ISubscriberIntTransactor{contract: contract}, ISubscriberIntFilterer: ISubscriberIntFilterer{contract: contract}}, nil
}

// NewISubscriberIntCaller creates a new read-only instance of ISubscriberInt, bound to a specific deployed contract.
func NewISubscriberIntCaller(address common.Address, caller bind.ContractCaller) (*ISubscriberIntCaller, error) {
	contract, err := bindISubscriberInt(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISubscriberIntCaller{contract: contract}, nil
}

// NewISubscriberIntTransactor creates a new write-only instance of ISubscriberInt, bound to a specific deployed contract.
func NewISubscriberIntTransactor(address common.Address, transactor bind.ContractTransactor) (*ISubscriberIntTransactor, error) {
	contract, err := bindISubscriberInt(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISubscriberIntTransactor{contract: contract}, nil
}

// NewISubscriberIntFilterer creates a new log filterer instance of ISubscriberInt, bound to a specific deployed contract.
func NewISubscriberIntFilterer(address common.Address, filterer bind.ContractFilterer) (*ISubscriberIntFilterer, error) {
	contract, err := bindISubscriberInt(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISubscriberIntFilterer{contract: contract}, nil
}

// bindISubscriberInt binds a generic wrapper to an already deployed contract.
func bindISubscriberInt(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISubscriberIntABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISubscriberInt *ISubscriberIntRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISubscriberInt.Contract.ISubscriberIntCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISubscriberInt *ISubscriberIntRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISubscriberInt.Contract.ISubscriberIntTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISubscriberInt *ISubscriberIntRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISubscriberInt.Contract.ISubscriberIntTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISubscriberInt *ISubscriberIntCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISubscriberInt.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISubscriberInt *ISubscriberIntTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISubscriberInt.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISubscriberInt *ISubscriberIntTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISubscriberInt.Contract.contract.Transact(opts, method, params...)
}

// AttachValue is a paid mutator transaction binding the contract method 0x7dc7c1b0.
//
// Solidity: function attachValue(int64 value) returns()
func (_ISubscriberInt *ISubscriberIntTransactor) AttachValue(opts *bind.TransactOpts, value int64) (*types.Transaction, error) {
	return _ISubscriberInt.contract.Transact(opts, "attachValue", value)
}

// AttachValue is a paid mutator transaction binding the contract method 0x7dc7c1b0.
//
// Solidity: function attachValue(int64 value) returns()
func (_ISubscriberInt *ISubscriberIntSession) AttachValue(value int64) (*types.Transaction, error) {
	return _ISubscriberInt.Contract.AttachValue(&_ISubscriberInt.TransactOpts, value)
}

// AttachValue is a paid mutator transaction binding the contract method 0x7dc7c1b0.
//
// Solidity: function attachValue(int64 value) returns()
func (_ISubscriberInt *ISubscriberIntTransactorSession) AttachValue(value int64) (*types.Transaction, error) {
	return _ISubscriberInt.Contract.AttachValue(&_ISubscriberInt.TransactOpts, value)
}

// ISubscriberStringABI is the input ABI used to generate the binding from.
const ISubscriberStringABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"name\":\"attachValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ISubscriberStringFuncSigs maps the 4-byte function signature to its string representation.
var ISubscriberStringFuncSigs = map[string]string{
	"bb327823": "attachValue(string)",
}

// ISubscriberString is an auto generated Go binding around an Ethereum contract.
type ISubscriberString struct {
	ISubscriberStringCaller     // Read-only binding to the contract
	ISubscriberStringTransactor // Write-only binding to the contract
	ISubscriberStringFilterer   // Log filterer for contract events
}

// ISubscriberStringCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISubscriberStringCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberStringTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISubscriberStringTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberStringFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISubscriberStringFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISubscriberStringSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISubscriberStringSession struct {
	Contract     *ISubscriberString // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ISubscriberStringCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISubscriberStringCallerSession struct {
	Contract *ISubscriberStringCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ISubscriberStringTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISubscriberStringTransactorSession struct {
	Contract     *ISubscriberStringTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ISubscriberStringRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISubscriberStringRaw struct {
	Contract *ISubscriberString // Generic contract binding to access the raw methods on
}

// ISubscriberStringCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISubscriberStringCallerRaw struct {
	Contract *ISubscriberStringCaller // Generic read-only contract binding to access the raw methods on
}

// ISubscriberStringTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISubscriberStringTransactorRaw struct {
	Contract *ISubscriberStringTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISubscriberString creates a new instance of ISubscriberString, bound to a specific deployed contract.
func NewISubscriberString(address common.Address, backend bind.ContractBackend) (*ISubscriberString, error) {
	contract, err := bindISubscriberString(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISubscriberString{ISubscriberStringCaller: ISubscriberStringCaller{contract: contract}, ISubscriberStringTransactor: ISubscriberStringTransactor{contract: contract}, ISubscriberStringFilterer: ISubscriberStringFilterer{contract: contract}}, nil
}

// NewISubscriberStringCaller creates a new read-only instance of ISubscriberString, bound to a specific deployed contract.
func NewISubscriberStringCaller(address common.Address, caller bind.ContractCaller) (*ISubscriberStringCaller, error) {
	contract, err := bindISubscriberString(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISubscriberStringCaller{contract: contract}, nil
}

// NewISubscriberStringTransactor creates a new write-only instance of ISubscriberString, bound to a specific deployed contract.
func NewISubscriberStringTransactor(address common.Address, transactor bind.ContractTransactor) (*ISubscriberStringTransactor, error) {
	contract, err := bindISubscriberString(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISubscriberStringTransactor{contract: contract}, nil
}

// NewISubscriberStringFilterer creates a new log filterer instance of ISubscriberString, bound to a specific deployed contract.
func NewISubscriberStringFilterer(address common.Address, filterer bind.ContractFilterer) (*ISubscriberStringFilterer, error) {
	contract, err := bindISubscriberString(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISubscriberStringFilterer{contract: contract}, nil
}

// bindISubscriberString binds a generic wrapper to an already deployed contract.
func bindISubscriberString(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISubscriberStringABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISubscriberString *ISubscriberStringRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISubscriberString.Contract.ISubscriberStringCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISubscriberString *ISubscriberStringRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISubscriberString.Contract.ISubscriberStringTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISubscriberString *ISubscriberStringRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISubscriberString.Contract.ISubscriberStringTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISubscriberString *ISubscriberStringCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISubscriberString.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISubscriberString *ISubscriberStringTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISubscriberString.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISubscriberString *ISubscriberStringTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISubscriberString.Contract.contract.Transact(opts, method, params...)
}

// AttachValue is a paid mutator transaction binding the contract method 0xbb327823.
//
// Solidity: function attachValue(string value) returns()
func (_ISubscriberString *ISubscriberStringTransactor) AttachValue(opts *bind.TransactOpts, value string) (*types.Transaction, error) {
	return _ISubscriberString.contract.Transact(opts, "attachValue", value)
}

// AttachValue is a paid mutator transaction binding the contract method 0xbb327823.
//
// Solidity: function attachValue(string value) returns()
func (_ISubscriberString *ISubscriberStringSession) AttachValue(value string) (*types.Transaction, error) {
	return _ISubscriberString.Contract.AttachValue(&_ISubscriberString.TransactOpts, value)
}

// AttachValue is a paid mutator transaction binding the contract method 0xbb327823.
//
// Solidity: function attachValue(string value) returns()
func (_ISubscriberString *ISubscriberStringTransactorSession) AttachValue(value string) (*types.Transaction, error) {
	return _ISubscriberString.Contract.AttachValue(&_ISubscriberString.TransactOpts, value)
}

// NModelsABI is the input ABI used to generate the binding from.
const NModelsABI = "[]"

// NModelsBin is the compiled bytecode used for deploying new contracts.
var NModelsBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212200c28306d968239ea8ce242901359bba17cef784a9c073c0bd7ec4d561426918c64736f6c63430007000033"

// DeployNModels deploys a new Ethereum contract, binding an instance of NModels to it.
func DeployNModels(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NModels, error) {
	parsed, err := abi.JSON(strings.NewReader(NModelsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NModelsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NModels{NModelsCaller: NModelsCaller{contract: contract}, NModelsTransactor: NModelsTransactor{contract: contract}, NModelsFilterer: NModelsFilterer{contract: contract}}, nil
}

// NModels is an auto generated Go binding around an Ethereum contract.
type NModels struct {
	NModelsCaller     // Read-only binding to the contract
	NModelsTransactor // Write-only binding to the contract
	NModelsFilterer   // Log filterer for contract events
}

// NModelsCaller is an auto generated read-only Go binding around an Ethereum contract.
type NModelsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NModelsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NModelsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NModelsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NModelsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NModelsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NModelsSession struct {
	Contract     *NModels          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NModelsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NModelsCallerSession struct {
	Contract *NModelsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// NModelsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NModelsTransactorSession struct {
	Contract     *NModelsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// NModelsRaw is an auto generated low-level Go binding around an Ethereum contract.
type NModelsRaw struct {
	Contract *NModels // Generic contract binding to access the raw methods on
}

// NModelsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NModelsCallerRaw struct {
	Contract *NModelsCaller // Generic read-only contract binding to access the raw methods on
}

// NModelsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NModelsTransactorRaw struct {
	Contract *NModelsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNModels creates a new instance of NModels, bound to a specific deployed contract.
func NewNModels(address common.Address, backend bind.ContractBackend) (*NModels, error) {
	contract, err := bindNModels(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NModels{NModelsCaller: NModelsCaller{contract: contract}, NModelsTransactor: NModelsTransactor{contract: contract}, NModelsFilterer: NModelsFilterer{contract: contract}}, nil
}

// NewNModelsCaller creates a new read-only instance of NModels, bound to a specific deployed contract.
func NewNModelsCaller(address common.Address, caller bind.ContractCaller) (*NModelsCaller, error) {
	contract, err := bindNModels(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NModelsCaller{contract: contract}, nil
}

// NewNModelsTransactor creates a new write-only instance of NModels, bound to a specific deployed contract.
func NewNModelsTransactor(address common.Address, transactor bind.ContractTransactor) (*NModelsTransactor, error) {
	contract, err := bindNModels(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NModelsTransactor{contract: contract}, nil
}

// NewNModelsFilterer creates a new log filterer instance of NModels, bound to a specific deployed contract.
func NewNModelsFilterer(address common.Address, filterer bind.ContractFilterer) (*NModelsFilterer, error) {
	contract, err := bindNModels(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NModelsFilterer{contract: contract}, nil
}

// bindNModels binds a generic wrapper to an already deployed contract.
func bindNModels(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NModelsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NModels *NModelsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NModels.Contract.NModelsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NModels *NModelsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NModels.Contract.NModelsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NModels *NModelsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NModels.Contract.NModelsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NModels *NModelsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NModels.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NModels *NModelsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NModels.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NModels *NModelsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NModels.Contract.contract.Transact(opts, method, params...)
}

// NebulaABI is the input ABI used to generate the binding from.
const NebulaABI = "[{\"inputs\":[{\"internalType\":\"enumNModels.DataType\",\"name\":\"newDataType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"newGravityContract\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"newOracle\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"newBftValue\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pulseId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"NewPulse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"NewSubscriber\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bftValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dataType\",\"outputs\":[{\"internalType\":\"enumNModels.DataType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscribersIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gravityContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newOracles\",\"type\":\"address[]\"}],\"name\":\"hashNewOracles\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isPulseSubSent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastPulseId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracles\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pulseQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pulses\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rounds\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"sendHashValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"pulseId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subId\",\"type\":\"bytes32\"}],\"name\":\"sendValueToSubByte\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int64\",\"name\":\"value\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"pulseId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subId\",\"type\":\"bytes32\"}],\"name\":\"sendValueToSubInt\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"pulseId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subId\",\"type\":\"bytes32\"}],\"name\":\"sendValueToSubString\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"minConfirmations\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"subscriptionIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"subscriptions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"minConfirmations\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subscriptionsQueue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"first\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"last\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newOracles\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"v\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"r\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"s\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"newRound\",\"type\":\"uint256\"}],\"name\":\"updateOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// NebulaFuncSigs maps the 4-byte function signature to its string representation.
var NebulaFuncSigs = map[string]string{
	"3cec1bdd": "bftValue()",
	"6175ff00": "dataType()",
	"40884c52": "getOracles()",
	"9505f6d4": "getSubscribersIds()",
	"770e58d5": "gravityContract()",
	"8bec345f": "hashNewOracles(address[])",
	"3ce7516d": "isPulseSubSent(uint256,bytes32)",
	"8d00662b": "lastPulseId()",
	"69a4246d": "oracleQueue()",
	"5b69a7d8": "oracles(uint256)",
	"1d11f944": "pulseQueue()",
	"0694fbb3": "pulses(uint256)",
	"8c65c81f": "rounds(uint256)",
	"bf2c0c42": "sendHashValue(bytes32,uint8[],bytes32[],bytes32[])",
	"ac557141": "sendValueToSubByte(bytes,uint256,bytes32)",
	"ff51063b": "sendValueToSubInt(int64,uint256,bytes32)",
	"9f95e525": "sendValueToSubString(string,uint256,bytes32)",
	"3527715d": "subscribe(address,uint8,uint256)",
	"8cafc358": "subscriptionIds(uint256)",
	"94259c6c": "subscriptions(bytes32)",
	"b48a9c9b": "subscriptionsQueue()",
	"febae9ea": "updateOracles(address[],uint8[],bytes32[],bytes32[],uint256)",
}

// NebulaBin is the compiled bytecode used for deploying new contracts.
var NebulaBin = "0x60806040523480156200001157600080fd5b5060405162001a3f38038062001a3f833981810160405260808110156200003757600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200006357600080fd5b9083019060208201858111156200007957600080fd5b82518660208202830111640100000000821117156200009757600080fd5b82525081516020918201928201910280838360005b83811015620000c6578181015183820152602001620000ac565b50505050919091016040525060200151600f80549193508692509060ff60a01b1916600160a01b836002811115620000fa57fe5b021790555081516200011490600d90602085019062000140565b50600e5550600f80546001600160a01b0319166001600160a01b039290921691909117905550620001cb565b82805482825590600052602060002090810192821562000198579160200282015b828111156200019857825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019062000161565b50620001a6929150620001aa565b5090565b5b80821115620001a65780546001600160a01b0319168155600101620001ab565b61186480620001db6000396000f3fe6080604052600436106101395760003560e01c80638c65c81f116100ab5780639f95e5251161006f5780639f95e525146104e1578063ac55714114610597578063b48a9c9b1461064d578063bf2c0c4214610662578063febae9ea1461081b578063ff51063b14610a5157610140565b80638c65c81f146104055780638cafc3581461042f5780638d00662b1461045957806394259c6c1461046e5780639505f6d4146104cc57610140565b806340884c52116100fd57806340884c521461024c5780635b69a7d8146102b15780636175ff00146102f757806369a4246d1461032d578063770e58d5146103425780638bec345f1461035757610140565b80630694fbb3146101455780631d11f944146101885780633527715d1461019d5780633ce7516d146101e15780633cec1bdd1461022557610140565b3661014057005b600080fd5b34801561015157600080fd5b5061016f6004803603602081101561016857600080fd5b5035610a8a565b6040805192835260208301919091528051918290030190f35b34801561019457600080fd5b5061016f610aa3565b3480156101a957600080fd5b506101df600480360360608110156101c057600080fd5b506001600160a01b038135169060ff6020820135169060400135610aac565b005b3480156101ed57600080fd5b506102116004803603604081101561020457600080fd5b5080359060200135610d38565b604080519115158252519081900360200190f35b34801561023157600080fd5b5061023a610d58565b60408051918252519081900360200190f35b34801561025857600080fd5b50610261610d5e565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561029d578181015183820152602001610285565b505050509050019250505060405180910390f35b3480156102bd57600080fd5b506102db600480360360208110156102d457600080fd5b5035610dc1565b604080516001600160a01b039092168252519081900360200190f35b34801561030357600080fd5b5061030c610de8565b6040518082600281111561031c57fe5b815260200191505060405180910390f35b34801561033957600080fd5b5061016f610df8565b34801561034e57600080fd5b506102db610e01565b34801561036357600080fd5b5061023a6004803603602081101561037a57600080fd5b810190602081018135600160201b81111561039457600080fd5b8201836020820111156103a657600080fd5b803590602001918460208302840111600160201b831117156103c757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550610e10945050505050565b34801561041157600080fd5b506102116004803603602081101561042857600080fd5b5035610ecc565b34801561043b57600080fd5b5061023a6004803603602081101561045257600080fd5b5035610ee1565b34801561046557600080fd5b5061023a610eff565b34801561047a57600080fd5b506104986004803603602081101561049157600080fd5b5035610f05565b604080516001600160a01b03958616815293909416602084015260ff90911682840152606082015290519081900360800190f35b3480156104d857600080fd5b50610261610f3e565b3480156104ed57600080fd5b506101df6004803603606081101561050457600080fd5b810190602081018135600160201b81111561051e57600080fd5b82018360208201111561053057600080fd5b803590602001918460018302840111600160201b8311171561055157600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505082359350505060200135610f95565b3480156105a357600080fd5b506101df600480360360608110156105ba57600080fd5b810190602081018135600160201b8111156105d457600080fd5b8201836020820111156105e657600080fd5b803590602001918460018302840111600160201b8311171561060757600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050823593505050602001356110d3565b34801561065957600080fd5b5061016f6111a5565b34801561066e57600080fd5b506101df6004803603608081101561068557600080fd5b81359190810190604081016020820135600160201b8111156106a657600080fd5b8201836020820111156106b857600080fd5b803590602001918460208302840111600160201b831117156106d957600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561072857600080fd5b82018360208201111561073a57600080fd5b803590602001918460208302840111600160201b8311171561075b57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156107aa57600080fd5b8201836020820111156107bc57600080fd5b803590602001918460208302840111600160201b831117156107dd57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295506111ae945050505050565b34801561082757600080fd5b506101df600480360360a081101561083e57600080fd5b810190602081018135600160201b81111561085857600080fd5b82018360208201111561086a57600080fd5b803590602001918460208302840111600160201b8311171561088b57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156108da57600080fd5b8201836020820111156108ec57600080fd5b803590602001918460208302840111600160201b8311171561090d57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561095c57600080fd5b82018360208201111561096e57600080fd5b803590602001918460208302840111600160201b8311171561098f57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156109de57600080fd5b8201836020820111156109f057600080fd5b803590602001918460208302840111600160201b83111715610a1157600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929550509135925061137d915050565b348015610a5d57600080fd5b506101df60048036036060811015610a7457600080fd5b50803560070b9060208101359060400135611614565b6013602052600090815260409020805460019091015482565b600954600a5482565b60408051600080356001600160e01b03191660208084019190915233606090811b602485015287901b6bffffffffffffffffffffffff1916603884015260f886901b6001600160f81b031916604c8401528351808403602d018152604d84019094528351919392606d019182918401908083835b60208310610b3f5780518252601f199092019160209182019101610b20565b51815160209384036101000a60001901801990921691161790526040805192909401828103601f1901835284528151918101919091206000818152601290925292902054919450506001600160a01b0316159150610bd29050576040805162461bcd60e51b815260206004820152600960248201526872712065786973747360b81b604482015290519081900360640190fd5b604080516080810182523381526001600160a01b03868116602080840191825260ff8881168587019081526060860189815260008981526012909452878420965187546001600160a01b0319908116918816919091178855945160018801805493519390961696169590951760ff60a01b1916600160a01b91909216021790915590516002909201919091558151632941b65560e21b81526005600482015260248101849052915173__$965488a2eb4d2c0ee0435ee6139709d0a6$__9263a506d954926044808301939192829003018186803b158015610cb257600080fd5b505af4158015610cc6573d6000803e3d6000fd5b5050601080546001810182556000919091527f1b6847dc741a1b0cd08d278845f9d819d87b734759afb55fe2de5cb82a9ae6720183905550506040805182815290517fce45703d589912fa7e0caf69c208323a21733441d675118cd08e8b4479c61daf9181900360200190a150505050565b601460209081526000928352604080842090915290825290205460ff1681565b600e5481565b6060600d805480602002602001604051908101604052809291908181526020018280548015610db657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610d98575b505050505090505b90565b600d8181548110610dce57fe5b6000918252602090912001546001600160a01b0316905081565b600f54600160a01b900460ff1681565b60015460025482565b600f546001600160a01b031681565b6000606060005b8351811015610ebd5781848281518110610e2d57fe5b60200260200101516040516020018083805190602001908083835b60208310610e675780518252601f199092019160209182019101610e48565b6001836020036101000a038019825116818451168082178552505050505050905001826001600160a01b031660601b81526014019250505060405160208183030381529060405291508080600101915050610e17565b50805160209091012092915050565b60006020819052908152604090205460ff1681565b60108181548110610eee57fe5b600091825260209091200154905081565b60115481565b6012602052600090815260409020805460018201546002909201546001600160a01b039182169291821691600160a01b900460ff169084565b60606010805480602002602001604051908101604052809291908181526020018280548015610db657602002820191906000526020600020905b815481526020019060010190808311610f78575050505050905090565b610f9d6116d0565b610fee576040805162461bcd60e51b815260206004820181905260248201527f63616c6c6572206973206e6f74206f6e65206f6620746865206f7261636c6573604482015290519081900360640190fd5b610ff8828261172a565b600081815260126020908152604080832060010154905163bb32782360e01b8152600481018381528751602483015287516001600160a01b039093169463bb32782394899492938493604490910192918601918190849084905b8381101561106a578181015183820152602001611052565b50505050905090810190601f1680156110975780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b1580156110b657600080fd5b505af11580156110ca573d6000803e3d6000fd5b50505050505050565b6110db6116d0565b61112c576040805162461bcd60e51b815260206004820181905260248201527f63616c6c6572206973206e6f74206f6e65206f6620746865206f7261636c6573604482015290519081900360640190fd5b611136828261172a565b600081815260126020908152604080832060010154905163cc32a15160e01b8152600481018381528751602483015287516001600160a01b039093169463cc32a1519489949293849360449091019291860191819084908490831561106a578181015183820152602001611052565b60055460065482565b6000805b600d548110156112b557600d81815481106111c957fe5b9060005260206000200160009054906101000a90046001600160a01b03166001600160a01b03166001878784815181106111ff57fe5b602002602001015187858151811061121357fe5b602002602001015187868151811061122757fe5b602002602001015160405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611283573d6000803e3d6000fd5b505050602060405103516001600160a01b0316146112a25760006112a5565b60015b60ff1691909101906001016111b2565b50600e54811015611301576040805162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a590818999d0818dbdd5b9d607a1b604482015290519081900360640190fd5b6011546040805180820182528781524360208083018281526001958601600081815260138452869020945185559051939095019290925582518481529182015280820188905290517f59327fa227d3fcc72b61eb5e7083cd72b448c22d7ea96849f4655db8de22d6399181900360600190a16011555050505050565b60008061138987610e10565b90506060600f60009054906101000a90046001600160a01b03166001600160a01b031663ad595b1a6040518163ffffffff1660e01b815260040160006040518083038186803b1580156113db57600080fd5b505afa1580156113ef573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052602081101561141857600080fd5b8101908080516040519392919084600160201b82111561143757600080fd5b90830190602082018581111561144c57600080fd5b82518660208202830111600160201b8211171561146857600080fd5b82525081516020918201928201910280838360005b8381101561149557818101518382015260200161147d565b50505050905001604052505050905060005b8151811015611592578181815181106114bc57fe5b60200260200101516001600160a01b03166001848a84815181106114dc57fe5b60200260200101518a85815181106114f057fe5b60200260200101518a868151811061150457fe5b602002602001015160405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015611560573d6000803e3d6000fd5b505050602060405103516001600160a01b03161461157f576000611582565b60015b60ff1693909301926001016114a7565b50600e548310156115de576040805162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a590818999d0818dbdd5b9d607a1b604482015290519081900360640190fd5b87516115f190600d9060208b01906117aa565b505050600091825250602081905260409020805460ff1916600117905550505050565b61161c6116d0565b61166d576040805162461bcd60e51b815260206004820181905260248201527f63616c6c6572206973206e6f74206f6e65206f6620746865206f7261636c6573604482015290519081900360640190fd5b611677828261172a565b6000818152601260205260408082206001015481516307dc7c1b60e41b8152600787900b600482015291516001600160a01b0390911692637dc7c1b0926024808201939182900301818387803b1580156110b657600080fd5b6000805b600d5481101561172257336001600160a01b0316600d82815481106116f557fe5b6000918252602090912001546001600160a01b0316141561171a576001915050610dbe565b6001016116d4565b506000905090565b600082815260146020908152604080832084845290915290205460ff1615611784576040805162461bcd60e51b81526020600482015260086024820152671cdd58881cd95b9d60c21b604482015290519081900360640190fd5b60009182526014602090815260408084209284529190529020805460ff19166001179055565b8280548282559060005260206000209081019282156117ff579160200282015b828111156117ff57825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906117ca565b5061180b92915061180f565b5090565b5b8082111561180b5780546001600160a01b031916815560010161181056fea26469706673582212208101116915832e931c81c88d6256ec95d42984872d972551a58f3a94f5b153bb64736f6c63430007000033"

// DeployNebula deploys a new Ethereum contract, binding an instance of Nebula to it.
func DeployNebula(auth *bind.TransactOpts, backend bind.ContractBackend, newDataType uint8, newGravityContract common.Address, newOracle []common.Address, newBftValue *big.Int) (common.Address, *types.Transaction, *Nebula, error) {
	parsed, err := abi.JSON(strings.NewReader(NebulaABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	queueLibAddr, _, _, _ := DeployQueueLib(auth, backend)
	NebulaBin = strings.Replace(NebulaBin, "__$965488a2eb4d2c0ee0435ee6139709d0a6$__", queueLibAddr.String()[2:], -1)

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NebulaBin), backend, newDataType, newGravityContract, newOracle, newBftValue)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Nebula{NebulaCaller: NebulaCaller{contract: contract}, NebulaTransactor: NebulaTransactor{contract: contract}, NebulaFilterer: NebulaFilterer{contract: contract}}, nil
}

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
func (_Nebula *NebulaRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_Nebula *NebulaCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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

// BftValue is a free data retrieval call binding the contract method 0x3cec1bdd.
//
// Solidity: function bftValue() view returns(uint256)
func (_Nebula *NebulaCaller) BftValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "bftValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BftValue is a free data retrieval call binding the contract method 0x3cec1bdd.
//
// Solidity: function bftValue() view returns(uint256)
func (_Nebula *NebulaSession) BftValue() (*big.Int, error) {
	return _Nebula.Contract.BftValue(&_Nebula.CallOpts)
}

// BftValue is a free data retrieval call binding the contract method 0x3cec1bdd.
//
// Solidity: function bftValue() view returns(uint256)
func (_Nebula *NebulaCallerSession) BftValue() (*big.Int, error) {
	return _Nebula.Contract.BftValue(&_Nebula.CallOpts)
}

// DataType is a free data retrieval call binding the contract method 0x6175ff00.
//
// Solidity: function dataType() view returns(uint8)
func (_Nebula *NebulaCaller) DataType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "dataType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// DataType is a free data retrieval call binding the contract method 0x6175ff00.
//
// Solidity: function dataType() view returns(uint8)
func (_Nebula *NebulaSession) DataType() (uint8, error) {
	return _Nebula.Contract.DataType(&_Nebula.CallOpts)
}

// DataType is a free data retrieval call binding the contract method 0x6175ff00.
//
// Solidity: function dataType() view returns(uint8)
func (_Nebula *NebulaCallerSession) DataType() (uint8, error) {
	return _Nebula.Contract.DataType(&_Nebula.CallOpts)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_Nebula *NebulaCaller) GetOracles(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "getOracles")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_Nebula *NebulaSession) GetOracles() ([]common.Address, error) {
	return _Nebula.Contract.GetOracles(&_Nebula.CallOpts)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_Nebula *NebulaCallerSession) GetOracles() ([]common.Address, error) {
	return _Nebula.Contract.GetOracles(&_Nebula.CallOpts)
}

// GetSubscribersIds is a free data retrieval call binding the contract method 0x9505f6d4.
//
// Solidity: function getSubscribersIds() view returns(bytes32[])
func (_Nebula *NebulaCaller) GetSubscribersIds(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "getSubscribersIds")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetSubscribersIds is a free data retrieval call binding the contract method 0x9505f6d4.
//
// Solidity: function getSubscribersIds() view returns(bytes32[])
func (_Nebula *NebulaSession) GetSubscribersIds() ([][32]byte, error) {
	return _Nebula.Contract.GetSubscribersIds(&_Nebula.CallOpts)
}

// GetSubscribersIds is a free data retrieval call binding the contract method 0x9505f6d4.
//
// Solidity: function getSubscribersIds() view returns(bytes32[])
func (_Nebula *NebulaCallerSession) GetSubscribersIds() ([][32]byte, error) {
	return _Nebula.Contract.GetSubscribersIds(&_Nebula.CallOpts)
}

// GravityContract is a free data retrieval call binding the contract method 0x770e58d5.
//
// Solidity: function gravityContract() view returns(address)
func (_Nebula *NebulaCaller) GravityContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "gravityContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GravityContract is a free data retrieval call binding the contract method 0x770e58d5.
//
// Solidity: function gravityContract() view returns(address)
func (_Nebula *NebulaSession) GravityContract() (common.Address, error) {
	return _Nebula.Contract.GravityContract(&_Nebula.CallOpts)
}

// GravityContract is a free data retrieval call binding the contract method 0x770e58d5.
//
// Solidity: function gravityContract() view returns(address)
func (_Nebula *NebulaCallerSession) GravityContract() (common.Address, error) {
	return _Nebula.Contract.GravityContract(&_Nebula.CallOpts)
}

// HashNewOracles is a free data retrieval call binding the contract method 0x8bec345f.
//
// Solidity: function hashNewOracles(address[] newOracles) pure returns(bytes32)
func (_Nebula *NebulaCaller) HashNewOracles(opts *bind.CallOpts, newOracles []common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "hashNewOracles", newOracles)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashNewOracles is a free data retrieval call binding the contract method 0x8bec345f.
//
// Solidity: function hashNewOracles(address[] newOracles) pure returns(bytes32)
func (_Nebula *NebulaSession) HashNewOracles(newOracles []common.Address) ([32]byte, error) {
	return _Nebula.Contract.HashNewOracles(&_Nebula.CallOpts, newOracles)
}

// HashNewOracles is a free data retrieval call binding the contract method 0x8bec345f.
//
// Solidity: function hashNewOracles(address[] newOracles) pure returns(bytes32)
func (_Nebula *NebulaCallerSession) HashNewOracles(newOracles []common.Address) ([32]byte, error) {
	return _Nebula.Contract.HashNewOracles(&_Nebula.CallOpts, newOracles)
}

// IsPulseSubSent is a free data retrieval call binding the contract method 0x3ce7516d.
//
// Solidity: function isPulseSubSent(uint256 , bytes32 ) view returns(bool)
func (_Nebula *NebulaCaller) IsPulseSubSent(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) (bool, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "isPulseSubSent", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPulseSubSent is a free data retrieval call binding the contract method 0x3ce7516d.
//
// Solidity: function isPulseSubSent(uint256 , bytes32 ) view returns(bool)
func (_Nebula *NebulaSession) IsPulseSubSent(arg0 *big.Int, arg1 [32]byte) (bool, error) {
	return _Nebula.Contract.IsPulseSubSent(&_Nebula.CallOpts, arg0, arg1)
}

// IsPulseSubSent is a free data retrieval call binding the contract method 0x3ce7516d.
//
// Solidity: function isPulseSubSent(uint256 , bytes32 ) view returns(bool)
func (_Nebula *NebulaCallerSession) IsPulseSubSent(arg0 *big.Int, arg1 [32]byte) (bool, error) {
	return _Nebula.Contract.IsPulseSubSent(&_Nebula.CallOpts, arg0, arg1)
}

// LastPulseId is a free data retrieval call binding the contract method 0x8d00662b.
//
// Solidity: function lastPulseId() view returns(uint256)
func (_Nebula *NebulaCaller) LastPulseId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "lastPulseId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastPulseId is a free data retrieval call binding the contract method 0x8d00662b.
//
// Solidity: function lastPulseId() view returns(uint256)
func (_Nebula *NebulaSession) LastPulseId() (*big.Int, error) {
	return _Nebula.Contract.LastPulseId(&_Nebula.CallOpts)
}

// LastPulseId is a free data retrieval call binding the contract method 0x8d00662b.
//
// Solidity: function lastPulseId() view returns(uint256)
func (_Nebula *NebulaCallerSession) LastPulseId() (*big.Int, error) {
	return _Nebula.Contract.LastPulseId(&_Nebula.CallOpts)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCaller) OracleQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "oracleQueue")

	outstruct := new(struct {
		First [32]byte
		Last  [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.First = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Last = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

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

// Oracles is a free data retrieval call binding the contract method 0x5b69a7d8.
//
// Solidity: function oracles(uint256 ) view returns(address)
func (_Nebula *NebulaCaller) Oracles(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "oracles", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracles is a free data retrieval call binding the contract method 0x5b69a7d8.
//
// Solidity: function oracles(uint256 ) view returns(address)
func (_Nebula *NebulaSession) Oracles(arg0 *big.Int) (common.Address, error) {
	return _Nebula.Contract.Oracles(&_Nebula.CallOpts, arg0)
}

// Oracles is a free data retrieval call binding the contract method 0x5b69a7d8.
//
// Solidity: function oracles(uint256 ) view returns(address)
func (_Nebula *NebulaCallerSession) Oracles(arg0 *big.Int) (common.Address, error) {
	return _Nebula.Contract.Oracles(&_Nebula.CallOpts, arg0)
}

// PulseQueue is a free data retrieval call binding the contract method 0x1d11f944.
//
// Solidity: function pulseQueue() view returns(bytes32 first, bytes32 last)
func (_Nebula *NebulaCaller) PulseQueue(opts *bind.CallOpts) (struct {
	First [32]byte
	Last  [32]byte
}, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "pulseQueue")

	outstruct := new(struct {
		First [32]byte
		Last  [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.First = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Last = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

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
// Solidity: function pulses(uint256 ) view returns(bytes32 dataHash, uint256 height)
func (_Nebula *NebulaCaller) Pulses(opts *bind.CallOpts, arg0 *big.Int) (struct {
	DataHash [32]byte
	Height   *big.Int
}, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "pulses", arg0)

	outstruct := new(struct {
		DataHash [32]byte
		Height   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DataHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Pulses is a free data retrieval call binding the contract method 0x0694fbb3.
//
// Solidity: function pulses(uint256 ) view returns(bytes32 dataHash, uint256 height)
func (_Nebula *NebulaSession) Pulses(arg0 *big.Int) (struct {
	DataHash [32]byte
	Height   *big.Int
}, error) {
	return _Nebula.Contract.Pulses(&_Nebula.CallOpts, arg0)
}

// Pulses is a free data retrieval call binding the contract method 0x0694fbb3.
//
// Solidity: function pulses(uint256 ) view returns(bytes32 dataHash, uint256 height)
func (_Nebula *NebulaCallerSession) Pulses(arg0 *big.Int) (struct {
	DataHash [32]byte
	Height   *big.Int
}, error) {
	return _Nebula.Contract.Pulses(&_Nebula.CallOpts, arg0)
}

// Rounds is a free data retrieval call binding the contract method 0x8c65c81f.
//
// Solidity: function rounds(uint256 ) view returns(bool)
func (_Nebula *NebulaCaller) Rounds(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "rounds", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Rounds is a free data retrieval call binding the contract method 0x8c65c81f.
//
// Solidity: function rounds(uint256 ) view returns(bool)
func (_Nebula *NebulaSession) Rounds(arg0 *big.Int) (bool, error) {
	return _Nebula.Contract.Rounds(&_Nebula.CallOpts, arg0)
}

// Rounds is a free data retrieval call binding the contract method 0x8c65c81f.
//
// Solidity: function rounds(uint256 ) view returns(bool)
func (_Nebula *NebulaCallerSession) Rounds(arg0 *big.Int) (bool, error) {
	return _Nebula.Contract.Rounds(&_Nebula.CallOpts, arg0)
}

// SubscriptionIds is a free data retrieval call binding the contract method 0x8cafc358.
//
// Solidity: function subscriptionIds(uint256 ) view returns(bytes32)
func (_Nebula *NebulaCaller) SubscriptionIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "subscriptionIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SubscriptionIds is a free data retrieval call binding the contract method 0x8cafc358.
//
// Solidity: function subscriptionIds(uint256 ) view returns(bytes32)
func (_Nebula *NebulaSession) SubscriptionIds(arg0 *big.Int) ([32]byte, error) {
	return _Nebula.Contract.SubscriptionIds(&_Nebula.CallOpts, arg0)
}

// SubscriptionIds is a free data retrieval call binding the contract method 0x8cafc358.
//
// Solidity: function subscriptionIds(uint256 ) view returns(bytes32)
func (_Nebula *NebulaCallerSession) SubscriptionIds(arg0 *big.Int) ([32]byte, error) {
	return _Nebula.Contract.SubscriptionIds(&_Nebula.CallOpts, arg0)
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
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "subscriptions", arg0)

	outstruct := new(struct {
		Owner            common.Address
		ContractAddress  common.Address
		MinConfirmations uint8
		Reward           *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ContractAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.MinConfirmations = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.Reward = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

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
	var out []interface{}
	err := _Nebula.contract.Call(opts, &out, "subscriptionsQueue")

	outstruct := new(struct {
		First [32]byte
		Last  [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.First = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Last = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

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

// SendHashValue is a paid mutator transaction binding the contract method 0xbf2c0c42.
//
// Solidity: function sendHashValue(bytes32 dataHash, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Nebula *NebulaTransactor) SendHashValue(opts *bind.TransactOpts, dataHash [32]byte, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "sendHashValue", dataHash, v, r, s)
}

// SendHashValue is a paid mutator transaction binding the contract method 0xbf2c0c42.
//
// Solidity: function sendHashValue(bytes32 dataHash, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Nebula *NebulaSession) SendHashValue(dataHash [32]byte, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendHashValue(&_Nebula.TransactOpts, dataHash, v, r, s)
}

// SendHashValue is a paid mutator transaction binding the contract method 0xbf2c0c42.
//
// Solidity: function sendHashValue(bytes32 dataHash, uint8[] v, bytes32[] r, bytes32[] s) returns()
func (_Nebula *NebulaTransactorSession) SendHashValue(dataHash [32]byte, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendHashValue(&_Nebula.TransactOpts, dataHash, v, r, s)
}

// SendValueToSubByte is a paid mutator transaction binding the contract method 0xac557141.
//
// Solidity: function sendValueToSubByte(bytes value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaTransactor) SendValueToSubByte(opts *bind.TransactOpts, value []byte, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "sendValueToSubByte", value, pulseId, subId)
}

// SendValueToSubByte is a paid mutator transaction binding the contract method 0xac557141.
//
// Solidity: function sendValueToSubByte(bytes value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaSession) SendValueToSubByte(value []byte, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendValueToSubByte(&_Nebula.TransactOpts, value, pulseId, subId)
}

// SendValueToSubByte is a paid mutator transaction binding the contract method 0xac557141.
//
// Solidity: function sendValueToSubByte(bytes value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaTransactorSession) SendValueToSubByte(value []byte, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendValueToSubByte(&_Nebula.TransactOpts, value, pulseId, subId)
}

// SendValueToSubInt is a paid mutator transaction binding the contract method 0xff51063b.
//
// Solidity: function sendValueToSubInt(int64 value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaTransactor) SendValueToSubInt(opts *bind.TransactOpts, value int64, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "sendValueToSubInt", value, pulseId, subId)
}

// SendValueToSubInt is a paid mutator transaction binding the contract method 0xff51063b.
//
// Solidity: function sendValueToSubInt(int64 value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaSession) SendValueToSubInt(value int64, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendValueToSubInt(&_Nebula.TransactOpts, value, pulseId, subId)
}

// SendValueToSubInt is a paid mutator transaction binding the contract method 0xff51063b.
//
// Solidity: function sendValueToSubInt(int64 value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaTransactorSession) SendValueToSubInt(value int64, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendValueToSubInt(&_Nebula.TransactOpts, value, pulseId, subId)
}

// SendValueToSubString is a paid mutator transaction binding the contract method 0x9f95e525.
//
// Solidity: function sendValueToSubString(string value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaTransactor) SendValueToSubString(opts *bind.TransactOpts, value string, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "sendValueToSubString", value, pulseId, subId)
}

// SendValueToSubString is a paid mutator transaction binding the contract method 0x9f95e525.
//
// Solidity: function sendValueToSubString(string value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaSession) SendValueToSubString(value string, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendValueToSubString(&_Nebula.TransactOpts, value, pulseId, subId)
}

// SendValueToSubString is a paid mutator transaction binding the contract method 0x9f95e525.
//
// Solidity: function sendValueToSubString(string value, uint256 pulseId, bytes32 subId) returns()
func (_Nebula *NebulaTransactorSession) SendValueToSubString(value string, pulseId *big.Int, subId [32]byte) (*types.Transaction, error) {
	return _Nebula.Contract.SendValueToSubString(&_Nebula.TransactOpts, value, pulseId, subId)
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

// UpdateOracles is a paid mutator transaction binding the contract method 0xfebae9ea.
//
// Solidity: function updateOracles(address[] newOracles, uint8[] v, bytes32[] r, bytes32[] s, uint256 newRound) returns()
func (_Nebula *NebulaTransactor) UpdateOracles(opts *bind.TransactOpts, newOracles []common.Address, v []uint8, r [][32]byte, s [][32]byte, newRound *big.Int) (*types.Transaction, error) {
	return _Nebula.contract.Transact(opts, "updateOracles", newOracles, v, r, s, newRound)
}

// UpdateOracles is a paid mutator transaction binding the contract method 0xfebae9ea.
//
// Solidity: function updateOracles(address[] newOracles, uint8[] v, bytes32[] r, bytes32[] s, uint256 newRound) returns()
func (_Nebula *NebulaSession) UpdateOracles(newOracles []common.Address, v []uint8, r [][32]byte, s [][32]byte, newRound *big.Int) (*types.Transaction, error) {
	return _Nebula.Contract.UpdateOracles(&_Nebula.TransactOpts, newOracles, v, r, s, newRound)
}

// UpdateOracles is a paid mutator transaction binding the contract method 0xfebae9ea.
//
// Solidity: function updateOracles(address[] newOracles, uint8[] v, bytes32[] r, bytes32[] s, uint256 newRound) returns()
func (_Nebula *NebulaTransactorSession) UpdateOracles(newOracles []common.Address, v []uint8, r [][32]byte, s [][32]byte, newRound *big.Int) (*types.Transaction, error) {
	return _Nebula.Contract.UpdateOracles(&_Nebula.TransactOpts, newOracles, v, r, s, newRound)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Nebula *NebulaTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nebula.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Nebula *NebulaSession) Receive() (*types.Transaction, error) {
	return _Nebula.Contract.Receive(&_Nebula.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Nebula *NebulaTransactorSession) Receive() (*types.Transaction, error) {
	return _Nebula.Contract.Receive(&_Nebula.TransactOpts)
}

// NebulaNewPulseIterator is returned from FilterNewPulse and is used to iterate over the raw logs and unpacked data for NewPulse events raised by the Nebula contract.
type NebulaNewPulseIterator struct {
	Event *NebulaNewPulse // Event containing the contract specifics and raw log

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
func (it *NebulaNewPulseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NebulaNewPulse)
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
		it.Event = new(NebulaNewPulse)
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
func (it *NebulaNewPulseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NebulaNewPulseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NebulaNewPulse represents a NewPulse event raised by the Nebula contract.
type NebulaNewPulse struct {
	PulseId  *big.Int
	Height   *big.Int
	DataHash [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewPulse is a free log retrieval operation binding the contract event 0x59327fa227d3fcc72b61eb5e7083cd72b448c22d7ea96849f4655db8de22d639.
//
// Solidity: event NewPulse(uint256 pulseId, uint256 height, bytes32 dataHash)
func (_Nebula *NebulaFilterer) FilterNewPulse(opts *bind.FilterOpts) (*NebulaNewPulseIterator, error) {

	logs, sub, err := _Nebula.contract.FilterLogs(opts, "NewPulse")
	if err != nil {
		return nil, err
	}
	return &NebulaNewPulseIterator{contract: _Nebula.contract, event: "NewPulse", logs: logs, sub: sub}, nil
}

// WatchNewPulse is a free log subscription operation binding the contract event 0x59327fa227d3fcc72b61eb5e7083cd72b448c22d7ea96849f4655db8de22d639.
//
// Solidity: event NewPulse(uint256 pulseId, uint256 height, bytes32 dataHash)
func (_Nebula *NebulaFilterer) WatchNewPulse(opts *bind.WatchOpts, sink chan<- *NebulaNewPulse) (event.Subscription, error) {

	logs, sub, err := _Nebula.contract.WatchLogs(opts, "NewPulse")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NebulaNewPulse)
				if err := _Nebula.contract.UnpackLog(event, "NewPulse", log); err != nil {
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

// ParseNewPulse is a log parse operation binding the contract event 0x59327fa227d3fcc72b61eb5e7083cd72b448c22d7ea96849f4655db8de22d639.
//
// Solidity: event NewPulse(uint256 pulseId, uint256 height, bytes32 dataHash)
func (_Nebula *NebulaFilterer) ParseNewPulse(log types.Log) (*NebulaNewPulse, error) {
	event := new(NebulaNewPulse)
	if err := _Nebula.contract.UnpackLog(event, "NewPulse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NebulaNewSubscriberIterator is returned from FilterNewSubscriber and is used to iterate over the raw logs and unpacked data for NewSubscriber events raised by the Nebula contract.
type NebulaNewSubscriberIterator struct {
	Event *NebulaNewSubscriber // Event containing the contract specifics and raw log

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
func (it *NebulaNewSubscriberIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NebulaNewSubscriber)
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
		it.Event = new(NebulaNewSubscriber)
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
func (it *NebulaNewSubscriberIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NebulaNewSubscriberIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NebulaNewSubscriber represents a NewSubscriber event raised by the Nebula contract.
type NebulaNewSubscriber struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterNewSubscriber is a free log retrieval operation binding the contract event 0xce45703d589912fa7e0caf69c208323a21733441d675118cd08e8b4479c61daf.
//
// Solidity: event NewSubscriber(bytes32 id)
func (_Nebula *NebulaFilterer) FilterNewSubscriber(opts *bind.FilterOpts) (*NebulaNewSubscriberIterator, error) {

	logs, sub, err := _Nebula.contract.FilterLogs(opts, "NewSubscriber")
	if err != nil {
		return nil, err
	}
	return &NebulaNewSubscriberIterator{contract: _Nebula.contract, event: "NewSubscriber", logs: logs, sub: sub}, nil
}

// WatchNewSubscriber is a free log subscription operation binding the contract event 0xce45703d589912fa7e0caf69c208323a21733441d675118cd08e8b4479c61daf.
//
// Solidity: event NewSubscriber(bytes32 id)
func (_Nebula *NebulaFilterer) WatchNewSubscriber(opts *bind.WatchOpts, sink chan<- *NebulaNewSubscriber) (event.Subscription, error) {

	logs, sub, err := _Nebula.contract.WatchLogs(opts, "NewSubscriber")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NebulaNewSubscriber)
				if err := _Nebula.contract.UnpackLog(event, "NewSubscriber", log); err != nil {
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

// ParseNewSubscriber is a log parse operation binding the contract event 0xce45703d589912fa7e0caf69c208323a21733441d675118cd08e8b4479c61daf.
//
// Solidity: event NewSubscriber(bytes32 id)
func (_Nebula *NebulaFilterer) ParseNewSubscriber(log types.Log) (*NebulaNewSubscriber, error) {
	event := new(NebulaNewSubscriber)
	if err := _Nebula.contract.UnpackLog(event, "NewSubscriber", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QueueLibABI is the input ABI used to generate the binding from.
const QueueLibABI = "[]"

// QueueLibFuncSigs maps the 4-byte function signature to its string representation.
var QueueLibFuncSigs = map[string]string{
	"9d6ad84b": "drop(QueueLib.Queue storage,bytes32)",
	"a506d954": "push(QueueLib.Queue storage,bytes32)",
}

// QueueLibBin is the compiled bytecode used for deploying new contracts.
var QueueLibBin = "0x610198610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100405760003560e01c80639d6ad84b14610045578063a506d95414610077575b600080fd5b81801561005157600080fd5b506100756004803603604081101561006857600080fd5b50803590602001356100a7565b005b81801561008357600080fd5b506100756004803603604081101561009a57600080fd5b5080359060200135610114565b6000818152600383016020908152604080832054600286019092529091205481156100e457600082815260028501602052604090208190556100e8565b8084555b8015610106576000818152600385016020526040902082905561010e565b600184018290555b50505050565b8154610129578082556001820181905561015e565b600182018054600090815260028401602081815260408084208690558454868552600388018352818520559190528120558190555b505056fea264697066735822122079a1d72e0215cebd639cf1a1ab795199cbc591edebb16f097af722bb00951f1d64736f6c63430007000033"

// DeployQueueLib deploys a new Ethereum contract, binding an instance of QueueLib to it.
func DeployQueueLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *QueueLib, error) {
	parsed, err := abi.JSON(strings.NewReader(QueueLibABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(QueueLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &QueueLib{QueueLibCaller: QueueLibCaller{contract: contract}, QueueLibTransactor: QueueLibTransactor{contract: contract}, QueueLibFilterer: QueueLibFilterer{contract: contract}}, nil
}

// QueueLib is an auto generated Go binding around an Ethereum contract.
type QueueLib struct {
	QueueLibCaller     // Read-only binding to the contract
	QueueLibTransactor // Write-only binding to the contract
	QueueLibFilterer   // Log filterer for contract events
}

// QueueLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type QueueLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueueLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QueueLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueueLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QueueLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QueueLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QueueLibSession struct {
	Contract     *QueueLib         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QueueLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QueueLibCallerSession struct {
	Contract *QueueLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// QueueLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QueueLibTransactorSession struct {
	Contract     *QueueLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// QueueLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type QueueLibRaw struct {
	Contract *QueueLib // Generic contract binding to access the raw methods on
}

// QueueLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QueueLibCallerRaw struct {
	Contract *QueueLibCaller // Generic read-only contract binding to access the raw methods on
}

// QueueLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QueueLibTransactorRaw struct {
	Contract *QueueLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQueueLib creates a new instance of QueueLib, bound to a specific deployed contract.
func NewQueueLib(address common.Address, backend bind.ContractBackend) (*QueueLib, error) {
	contract, err := bindQueueLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QueueLib{QueueLibCaller: QueueLibCaller{contract: contract}, QueueLibTransactor: QueueLibTransactor{contract: contract}, QueueLibFilterer: QueueLibFilterer{contract: contract}}, nil
}

// NewQueueLibCaller creates a new read-only instance of QueueLib, bound to a specific deployed contract.
func NewQueueLibCaller(address common.Address, caller bind.ContractCaller) (*QueueLibCaller, error) {
	contract, err := bindQueueLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QueueLibCaller{contract: contract}, nil
}

// NewQueueLibTransactor creates a new write-only instance of QueueLib, bound to a specific deployed contract.
func NewQueueLibTransactor(address common.Address, transactor bind.ContractTransactor) (*QueueLibTransactor, error) {
	contract, err := bindQueueLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QueueLibTransactor{contract: contract}, nil
}

// NewQueueLibFilterer creates a new log filterer instance of QueueLib, bound to a specific deployed contract.
func NewQueueLibFilterer(address common.Address, filterer bind.ContractFilterer) (*QueueLibFilterer, error) {
	contract, err := bindQueueLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QueueLibFilterer{contract: contract}, nil
}

// bindQueueLib binds a generic wrapper to an already deployed contract.
func bindQueueLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(QueueLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QueueLib *QueueLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QueueLib.Contract.QueueLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QueueLib *QueueLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QueueLib.Contract.QueueLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QueueLib *QueueLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QueueLib.Contract.QueueLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QueueLib *QueueLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QueueLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QueueLib *QueueLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QueueLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QueueLib *QueueLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QueueLib.Contract.contract.Transact(opts, method, params...)
}
