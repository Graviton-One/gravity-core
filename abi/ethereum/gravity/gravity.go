// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gravity

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
