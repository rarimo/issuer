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

// Reference imports to suppress errors if they are not otherwise used.
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

// Struct2 is an auto generated low-level Go binding around an user-defined struct.
type Struct2 struct {
	Root                *big.Int
	ReplacedByRoot      *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtTimestamp *big.Int
	CreatedAtBlock      *big.Int
	ReplacedAtBlock     *big.Int
}

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	Id                  *big.Int
	State               *big.Int
	ReplacedByState     *big.Int
	CreatedAtTimestamp  *big.Int
	ReplacedAtTimestamp *big.Int
	CreatedAtBlock      *big.Int
	ReplacedAtBlock     *big.Int
}

// Struct1 is an auto generated low-level Go binding around an user-defined struct.
type Struct1 struct {
	Root     *big.Int
	Siblings [32]*big.Int
	OldKey   *big.Int
	OldValue *big.Int
	IsOld0   bool
	Key      *big.Int
	Value    *big.Int
	Fnc      *big.Int
}

// StateStoreMetaData contains all meta data concerning the StateStore contract.
var StateStoreMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockN\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getAllStateInfosById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structStateInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getGISTProof\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getGISTProofByBlock\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_root\",\"type\":\"uint256\"}],\"name\":\"getGISTProofByRoot\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"getGISTProofByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256[32]\",\"name\":\"siblings\",\"type\":\"uint256[32]\"},{\"internalType\":\"uint256\",\"name\":\"oldKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isOld0\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fnc\",\"type\":\"uint256\"}],\"internalType\":\"structProof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGISTRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_end\",\"type\":\"uint256\"}],\"name\":\"getGISTRootHistory\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGISTRootHistoryLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_root\",\"type\":\"uint256\"}],\"name\":\"getGISTRootInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getGISTRootInfoByBlock\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"getGISTRootInfoByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structRootInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getStateInfoById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structStateInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_state\",\"type\":\"uint256\"}],\"name\":\"getStateInfoByState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedByState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedAtBlock\",\"type\":\"uint256\"}],\"internalType\":\"structStateInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"_verifierContractAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newVerifierAddr\",\"type\":\"address\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateEntries\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"replacedBy\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"statesHistories\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_oldState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newState\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_isOldStateGenesis\",\"type\":\"bool\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"name\":\"transitState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]\"",
}

// StateStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use StateStoreMetaData.ABI instead.
var StateStoreABI = StateStoreMetaData.ABI

// StateStore is an auto generated Go binding around an Ethereum contract.
type StateStore struct {
	StateStoreCaller     // Read-only binding to the contract
	StateStoreTransactor // Write-only binding to the contract
	StateStoreFilterer   // Log filterer for contract events
}

// StateStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateStoreSession struct {
	Contract     *StateStore       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateStoreCallerSession struct {
	Contract *StateStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// StateStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateStoreTransactorSession struct {
	Contract     *StateStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// StateStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateStoreRaw struct {
	Contract *StateStore // Generic contract binding to access the raw methods on
}

// StateStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateStoreCallerRaw struct {
	Contract *StateStoreCaller // Generic read-only contract binding to access the raw methods on
}

// StateStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateStoreTransactorRaw struct {
	Contract *StateStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStateStore creates a new instance of StateStore, bound to a specific deployed contract.
func NewStateStore(address common.Address, backend bind.ContractBackend) (*StateStore, error) {
	contract, err := bindStateStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StateStore{StateStoreCaller: StateStoreCaller{contract: contract}, StateStoreTransactor: StateStoreTransactor{contract: contract}, StateStoreFilterer: StateStoreFilterer{contract: contract}}, nil
}

// NewStateStoreCaller creates a new read-only instance of StateStore, bound to a specific deployed contract.
func NewStateStoreCaller(address common.Address, caller bind.ContractCaller) (*StateStoreCaller, error) {
	contract, err := bindStateStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateStoreCaller{contract: contract}, nil
}

// NewStateStoreTransactor creates a new write-only instance of StateStore, bound to a specific deployed contract.
func NewStateStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*StateStoreTransactor, error) {
	contract, err := bindStateStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateStoreTransactor{contract: contract}, nil
}

// NewStateStoreFilterer creates a new log filterer instance of StateStore, bound to a specific deployed contract.
func NewStateStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*StateStoreFilterer, error) {
	contract, err := bindStateStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateStoreFilterer{contract: contract}, nil
}

// bindStateStore binds a generic wrapper to an already deployed contract.
func bindStateStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StateStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateStore *StateStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateStore.Contract.StateStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateStore *StateStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateStore.Contract.StateStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateStore *StateStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateStore.Contract.StateStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateStore *StateStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateStore *StateStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateStore *StateStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateStore.Contract.contract.Transact(opts, method, params...)
}

// GetAllStateInfosById is a free data retrieval call binding the contract method 0x93485cee.
//
// Solidity: function getAllStateInfosById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_StateStore *StateStoreCaller) GetAllStateInfosById(opts *bind.CallOpts, _id *big.Int) ([]Struct0, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getAllStateInfosById", _id)

	if err != nil {
		return *new([]Struct0), err
	}

	out0 := *abi.ConvertType(out[0], new([]Struct0)).(*[]Struct0)

	return out0, err

}

// GetAllStateInfosById is a free data retrieval call binding the contract method 0x93485cee.
//
// Solidity: function getAllStateInfosById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_StateStore *StateStoreSession) GetAllStateInfosById(_id *big.Int) ([]Struct0, error) {
	return _StateStore.Contract.GetAllStateInfosById(&_StateStore.CallOpts, _id)
}

// GetAllStateInfosById is a free data retrieval call binding the contract method 0x93485cee.
//
// Solidity: function getAllStateInfosById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_StateStore *StateStoreCallerSession) GetAllStateInfosById(_id *big.Int) ([]Struct0, error) {
	return _StateStore.Contract.GetAllStateInfosById(&_StateStore.CallOpts, _id)
}

// GetGISTProof is a free data retrieval call binding the contract method 0x3025bb8c.
//
// Solidity: function getGISTProof(uint256 _id) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTProof(opts *bind.CallOpts, _id *big.Int) (Struct1, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTProof", _id)

	if err != nil {
		return *new(Struct1), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct1)).(*Struct1)

	return out0, err

}

// GetGISTProof is a free data retrieval call binding the contract method 0x3025bb8c.
//
// Solidity: function getGISTProof(uint256 _id) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTProof(_id *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProof(&_StateStore.CallOpts, _id)
}

// GetGISTProof is a free data retrieval call binding the contract method 0x3025bb8c.
//
// Solidity: function getGISTProof(uint256 _id) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTProof(_id *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProof(&_StateStore.CallOpts, _id)
}

// GetGISTProofByBlock is a free data retrieval call binding the contract method 0x046ff140.
//
// Solidity: function getGISTProofByBlock(uint256 _id, uint256 _block) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTProofByBlock(opts *bind.CallOpts, _id *big.Int, _block *big.Int) (Struct1, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTProofByBlock", _id, _block)

	if err != nil {
		return *new(Struct1), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct1)).(*Struct1)

	return out0, err

}

// GetGISTProofByBlock is a free data retrieval call binding the contract method 0x046ff140.
//
// Solidity: function getGISTProofByBlock(uint256 _id, uint256 _block) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTProofByBlock(_id *big.Int, _block *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProofByBlock(&_StateStore.CallOpts, _id, _block)
}

// GetGISTProofByBlock is a free data retrieval call binding the contract method 0x046ff140.
//
// Solidity: function getGISTProofByBlock(uint256 _id, uint256 _block) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTProofByBlock(_id *big.Int, _block *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProofByBlock(&_StateStore.CallOpts, _id, _block)
}

// GetGISTProofByRoot is a free data retrieval call binding the contract method 0xe12a36c0.
//
// Solidity: function getGISTProofByRoot(uint256 _id, uint256 _root) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTProofByRoot(opts *bind.CallOpts, _id *big.Int, _root *big.Int) (Struct1, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTProofByRoot", _id, _root)

	if err != nil {
		return *new(Struct1), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct1)).(*Struct1)

	return out0, err

}

// GetGISTProofByRoot is a free data retrieval call binding the contract method 0xe12a36c0.
//
// Solidity: function getGISTProofByRoot(uint256 _id, uint256 _root) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTProofByRoot(_id *big.Int, _root *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProofByRoot(&_StateStore.CallOpts, _id, _root)
}

// GetGISTProofByRoot is a free data retrieval call binding the contract method 0xe12a36c0.
//
// Solidity: function getGISTProofByRoot(uint256 _id, uint256 _root) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTProofByRoot(_id *big.Int, _root *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProofByRoot(&_StateStore.CallOpts, _id, _root)
}

// GetGISTProofByTime is a free data retrieval call binding the contract method 0xd51afebf.
//
// Solidity: function getGISTProofByTime(uint256 _id, uint256 _timestamp) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTProofByTime(opts *bind.CallOpts, _id *big.Int, _timestamp *big.Int) (Struct1, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTProofByTime", _id, _timestamp)

	if err != nil {
		return *new(Struct1), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct1)).(*Struct1)

	return out0, err

}

// GetGISTProofByTime is a free data retrieval call binding the contract method 0xd51afebf.
//
// Solidity: function getGISTProofByTime(uint256 _id, uint256 _timestamp) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTProofByTime(_id *big.Int, _timestamp *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProofByTime(&_StateStore.CallOpts, _id, _timestamp)
}

// GetGISTProofByTime is a free data retrieval call binding the contract method 0xd51afebf.
//
// Solidity: function getGISTProofByTime(uint256 _id, uint256 _timestamp) view returns((uint256,uint256[32],uint256,uint256,bool,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTProofByTime(_id *big.Int, _timestamp *big.Int) (Struct1, error) {
	return _StateStore.Contract.GetGISTProofByTime(&_StateStore.CallOpts, _id, _timestamp)
}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_StateStore *StateStoreCaller) GetGISTRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_StateStore *StateStoreSession) GetGISTRoot() (*big.Int, error) {
	return _StateStore.Contract.GetGISTRoot(&_StateStore.CallOpts)
}

// GetGISTRoot is a free data retrieval call binding the contract method 0x2439e3a6.
//
// Solidity: function getGISTRoot() view returns(uint256)
func (_StateStore *StateStoreCallerSession) GetGISTRoot() (*big.Int, error) {
	return _StateStore.Contract.GetGISTRoot(&_StateStore.CallOpts)
}

// GetGISTRootHistory is a free data retrieval call binding the contract method 0x2f7670e4.
//
// Solidity: function getGISTRootHistory(uint256 _start, uint256 _end) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_StateStore *StateStoreCaller) GetGISTRootHistory(opts *bind.CallOpts, _start *big.Int, _end *big.Int) ([]Struct2, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTRootHistory", _start, _end)

	if err != nil {
		return *new([]Struct2), err
	}

	out0 := *abi.ConvertType(out[0], new([]Struct2)).(*[]Struct2)

	return out0, err

}

// GetGISTRootHistory is a free data retrieval call binding the contract method 0x2f7670e4.
//
// Solidity: function getGISTRootHistory(uint256 _start, uint256 _end) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_StateStore *StateStoreSession) GetGISTRootHistory(_start *big.Int, _end *big.Int) ([]Struct2, error) {
	return _StateStore.Contract.GetGISTRootHistory(&_StateStore.CallOpts, _start, _end)
}

// GetGISTRootHistory is a free data retrieval call binding the contract method 0x2f7670e4.
//
// Solidity: function getGISTRootHistory(uint256 _start, uint256 _end) view returns((uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_StateStore *StateStoreCallerSession) GetGISTRootHistory(_start *big.Int, _end *big.Int) ([]Struct2, error) {
	return _StateStore.Contract.GetGISTRootHistory(&_StateStore.CallOpts, _start, _end)
}

// GetGISTRootHistoryLength is a free data retrieval call binding the contract method 0xdccbd57a.
//
// Solidity: function getGISTRootHistoryLength() view returns(uint256)
func (_StateStore *StateStoreCaller) GetGISTRootHistoryLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTRootHistoryLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGISTRootHistoryLength is a free data retrieval call binding the contract method 0xdccbd57a.
//
// Solidity: function getGISTRootHistoryLength() view returns(uint256)
func (_StateStore *StateStoreSession) GetGISTRootHistoryLength() (*big.Int, error) {
	return _StateStore.Contract.GetGISTRootHistoryLength(&_StateStore.CallOpts)
}

// GetGISTRootHistoryLength is a free data retrieval call binding the contract method 0xdccbd57a.
//
// Solidity: function getGISTRootHistoryLength() view returns(uint256)
func (_StateStore *StateStoreCallerSession) GetGISTRootHistoryLength() (*big.Int, error) {
	return _StateStore.Contract.GetGISTRootHistoryLength(&_StateStore.CallOpts)
}

// GetGISTRootInfo is a free data retrieval call binding the contract method 0x7c1a66de.
//
// Solidity: function getGISTRootInfo(uint256 _root) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTRootInfo(opts *bind.CallOpts, _root *big.Int) (Struct2, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTRootInfo", _root)

	if err != nil {
		return *new(Struct2), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct2)).(*Struct2)

	return out0, err

}

// GetGISTRootInfo is a free data retrieval call binding the contract method 0x7c1a66de.
//
// Solidity: function getGISTRootInfo(uint256 _root) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTRootInfo(_root *big.Int) (Struct2, error) {
	return _StateStore.Contract.GetGISTRootInfo(&_StateStore.CallOpts, _root)
}

// GetGISTRootInfo is a free data retrieval call binding the contract method 0x7c1a66de.
//
// Solidity: function getGISTRootInfo(uint256 _root) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTRootInfo(_root *big.Int) (Struct2, error) {
	return _StateStore.Contract.GetGISTRootInfo(&_StateStore.CallOpts, _root)
}

// GetGISTRootInfoByBlock is a free data retrieval call binding the contract method 0x5845e530.
//
// Solidity: function getGISTRootInfoByBlock(uint256 _block) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTRootInfoByBlock(opts *bind.CallOpts, _block *big.Int) (Struct2, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTRootInfoByBlock", _block)

	if err != nil {
		return *new(Struct2), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct2)).(*Struct2)

	return out0, err

}

// GetGISTRootInfoByBlock is a free data retrieval call binding the contract method 0x5845e530.
//
// Solidity: function getGISTRootInfoByBlock(uint256 _block) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTRootInfoByBlock(_block *big.Int) (Struct2, error) {
	return _StateStore.Contract.GetGISTRootInfoByBlock(&_StateStore.CallOpts, _block)
}

// GetGISTRootInfoByBlock is a free data retrieval call binding the contract method 0x5845e530.
//
// Solidity: function getGISTRootInfoByBlock(uint256 _block) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTRootInfoByBlock(_block *big.Int) (Struct2, error) {
	return _StateStore.Contract.GetGISTRootInfoByBlock(&_StateStore.CallOpts, _block)
}

// GetGISTRootInfoByTime is a free data retrieval call binding the contract method 0x0ef6e65b.
//
// Solidity: function getGISTRootInfoByTime(uint256 _timestamp) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetGISTRootInfoByTime(opts *bind.CallOpts, _timestamp *big.Int) (Struct2, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getGISTRootInfoByTime", _timestamp)

	if err != nil {
		return *new(Struct2), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct2)).(*Struct2)

	return out0, err

}

// GetGISTRootInfoByTime is a free data retrieval call binding the contract method 0x0ef6e65b.
//
// Solidity: function getGISTRootInfoByTime(uint256 _timestamp) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetGISTRootInfoByTime(_timestamp *big.Int) (Struct2, error) {
	return _StateStore.Contract.GetGISTRootInfoByTime(&_StateStore.CallOpts, _timestamp)
}

// GetGISTRootInfoByTime is a free data retrieval call binding the contract method 0x0ef6e65b.
//
// Solidity: function getGISTRootInfoByTime(uint256 _timestamp) view returns((uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetGISTRootInfoByTime(_timestamp *big.Int) (Struct2, error) {
	return _StateStore.Contract.GetGISTRootInfoByTime(&_StateStore.CallOpts, _timestamp)
}

// GetStateInfoById is a free data retrieval call binding the contract method 0xb4bdea55.
//
// Solidity: function getStateInfoById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetStateInfoById(opts *bind.CallOpts, _id *big.Int) (Struct0, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getStateInfoById", _id)

	if err != nil {
		return *new(Struct0), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct0)).(*Struct0)

	return out0, err

}

// GetStateInfoById is a free data retrieval call binding the contract method 0xb4bdea55.
//
// Solidity: function getStateInfoById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetStateInfoById(_id *big.Int) (Struct0, error) {
	return _StateStore.Contract.GetStateInfoById(&_StateStore.CallOpts, _id)
}

// GetStateInfoById is a free data retrieval call binding the contract method 0xb4bdea55.
//
// Solidity: function getStateInfoById(uint256 _id) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetStateInfoById(_id *big.Int) (Struct0, error) {
	return _StateStore.Contract.GetStateInfoById(&_StateStore.CallOpts, _id)
}

// GetStateInfoByState is a free data retrieval call binding the contract method 0x3622b0bc.
//
// Solidity: function getStateInfoByState(uint256 _state) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCaller) GetStateInfoByState(opts *bind.CallOpts, _state *big.Int) (Struct0, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getStateInfoByState", _state)

	if err != nil {
		return *new(Struct0), err
	}

	out0 := *abi.ConvertType(out[0], new(Struct0)).(*Struct0)

	return out0, err

}

// GetStateInfoByState is a free data retrieval call binding the contract method 0x3622b0bc.
//
// Solidity: function getStateInfoByState(uint256 _state) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreSession) GetStateInfoByState(_state *big.Int) (Struct0, error) {
	return _StateStore.Contract.GetStateInfoByState(&_StateStore.CallOpts, _state)
}

// GetStateInfoByState is a free data retrieval call binding the contract method 0x3622b0bc.
//
// Solidity: function getStateInfoByState(uint256 _state) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_StateStore *StateStoreCallerSession) GetStateInfoByState(_state *big.Int) (Struct0, error) {
	return _StateStore.Contract.GetStateInfoByState(&_StateStore.CallOpts, _state)
}

// GetVerifier is a free data retrieval call binding the contract method 0x46657fe9.
//
// Solidity: function getVerifier() view returns(address)
func (_StateStore *StateStoreCaller) GetVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "getVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVerifier is a free data retrieval call binding the contract method 0x46657fe9.
//
// Solidity: function getVerifier() view returns(address)
func (_StateStore *StateStoreSession) GetVerifier() (common.Address, error) {
	return _StateStore.Contract.GetVerifier(&_StateStore.CallOpts)
}

// GetVerifier is a free data retrieval call binding the contract method 0x46657fe9.
//
// Solidity: function getVerifier() view returns(address)
func (_StateStore *StateStoreCallerSession) GetVerifier() (common.Address, error) {
	return _StateStore.Contract.GetVerifier(&_StateStore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StateStore *StateStoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StateStore *StateStoreSession) Owner() (common.Address, error) {
	return _StateStore.Contract.Owner(&_StateStore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StateStore *StateStoreCallerSession) Owner() (common.Address, error) {
	return _StateStore.Contract.Owner(&_StateStore.CallOpts)
}

// StateEntries is a free data retrieval call binding the contract method 0x3d8c1445.
//
// Solidity: function stateEntries(uint256 ) view returns(uint256 id, uint256 timestamp, uint256 block, uint256 replacedBy)
func (_StateStore *StateStoreCaller) StateEntries(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id         *big.Int
	Timestamp  *big.Int
	Block      *big.Int
	ReplacedBy *big.Int
}, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "stateEntries", arg0)

	outstruct := new(struct {
		Id         *big.Int
		Timestamp  *big.Int
		Block      *big.Int
		ReplacedBy *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Block = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.ReplacedBy = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StateEntries is a free data retrieval call binding the contract method 0x3d8c1445.
//
// Solidity: function stateEntries(uint256 ) view returns(uint256 id, uint256 timestamp, uint256 block, uint256 replacedBy)
func (_StateStore *StateStoreSession) StateEntries(arg0 *big.Int) (struct {
	Id         *big.Int
	Timestamp  *big.Int
	Block      *big.Int
	ReplacedBy *big.Int
}, error) {
	return _StateStore.Contract.StateEntries(&_StateStore.CallOpts, arg0)
}

// StateEntries is a free data retrieval call binding the contract method 0x3d8c1445.
//
// Solidity: function stateEntries(uint256 ) view returns(uint256 id, uint256 timestamp, uint256 block, uint256 replacedBy)
func (_StateStore *StateStoreCallerSession) StateEntries(arg0 *big.Int) (struct {
	Id         *big.Int
	Timestamp  *big.Int
	Block      *big.Int
	ReplacedBy *big.Int
}, error) {
	return _StateStore.Contract.StateEntries(&_StateStore.CallOpts, arg0)
}

// StatesHistories is a free data retrieval call binding the contract method 0xb9617362.
//
// Solidity: function statesHistories(uint256 , uint256 ) view returns(uint256)
func (_StateStore *StateStoreCaller) StatesHistories(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "statesHistories", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StatesHistories is a free data retrieval call binding the contract method 0xb9617362.
//
// Solidity: function statesHistories(uint256 , uint256 ) view returns(uint256)
func (_StateStore *StateStoreSession) StatesHistories(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _StateStore.Contract.StatesHistories(&_StateStore.CallOpts, arg0, arg1)
}

// StatesHistories is a free data retrieval call binding the contract method 0xb9617362.
//
// Solidity: function statesHistories(uint256 , uint256 ) view returns(uint256)
func (_StateStore *StateStoreCallerSession) StatesHistories(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _StateStore.Contract.StatesHistories(&_StateStore.CallOpts, arg0, arg1)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_StateStore *StateStoreCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StateStore.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_StateStore *StateStoreSession) Verifier() (common.Address, error) {
	return _StateStore.Contract.Verifier(&_StateStore.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_StateStore *StateStoreCallerSession) Verifier() (common.Address, error) {
	return _StateStore.Contract.Verifier(&_StateStore.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_StateStore *StateStoreTransactor) Initialize(opts *bind.TransactOpts, _verifierContractAddr common.Address) (*types.Transaction, error) {
	return _StateStore.contract.Transact(opts, "initialize", _verifierContractAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_StateStore *StateStoreSession) Initialize(_verifierContractAddr common.Address) (*types.Transaction, error) {
	return _StateStore.Contract.Initialize(&_StateStore.TransactOpts, _verifierContractAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _verifierContractAddr) returns()
func (_StateStore *StateStoreTransactorSession) Initialize(_verifierContractAddr common.Address) (*types.Transaction, error) {
	return _StateStore.Contract.Initialize(&_StateStore.TransactOpts, _verifierContractAddr)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StateStore *StateStoreTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateStore.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StateStore *StateStoreSession) RenounceOwnership() (*types.Transaction, error) {
	return _StateStore.Contract.RenounceOwnership(&_StateStore.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StateStore *StateStoreTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _StateStore.Contract.RenounceOwnership(&_StateStore.TransactOpts)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address _newVerifierAddr) returns()
func (_StateStore *StateStoreTransactor) SetVerifier(opts *bind.TransactOpts, _newVerifierAddr common.Address) (*types.Transaction, error) {
	return _StateStore.contract.Transact(opts, "setVerifier", _newVerifierAddr)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address _newVerifierAddr) returns()
func (_StateStore *StateStoreSession) SetVerifier(_newVerifierAddr common.Address) (*types.Transaction, error) {
	return _StateStore.Contract.SetVerifier(&_StateStore.TransactOpts, _newVerifierAddr)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x5437988d.
//
// Solidity: function setVerifier(address _newVerifierAddr) returns()
func (_StateStore *StateStoreTransactorSession) SetVerifier(_newVerifierAddr common.Address) (*types.Transaction, error) {
	return _StateStore.Contract.SetVerifier(&_StateStore.TransactOpts, _newVerifierAddr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StateStore *StateStoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _StateStore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StateStore *StateStoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StateStore.Contract.TransferOwnership(&_StateStore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StateStore *StateStoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StateStore.Contract.TransferOwnership(&_StateStore.TransactOpts, newOwner)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 _id, uint256 _oldState, uint256 _newState, bool _isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_StateStore *StateStoreTransactor) TransitState(opts *bind.TransactOpts, _id *big.Int, _oldState *big.Int, _newState *big.Int, _isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _StateStore.contract.Transact(opts, "transitState", _id, _oldState, _newState, _isOldStateGenesis, a, b, c)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 _id, uint256 _oldState, uint256 _newState, bool _isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_StateStore *StateStoreSession) TransitState(_id *big.Int, _oldState *big.Int, _newState *big.Int, _isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _StateStore.Contract.TransitState(&_StateStore.TransactOpts, _id, _oldState, _newState, _isOldStateGenesis, a, b, c)
}

// TransitState is a paid mutator transaction binding the contract method 0x28f88a65.
//
// Solidity: function transitState(uint256 _id, uint256 _oldState, uint256 _newState, bool _isOldStateGenesis, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_StateStore *StateStoreTransactorSession) TransitState(_id *big.Int, _oldState *big.Int, _newState *big.Int, _isOldStateGenesis bool, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _StateStore.Contract.TransitState(&_StateStore.TransactOpts, _id, _oldState, _newState, _isOldStateGenesis, a, b, c)
}

// StateStoreInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the StateStore contract.
type StateStoreInitializedIterator struct {
	Event *StateStoreInitialized // Event containing the contract specifics and raw log

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
func (it *StateStoreInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateStoreInitialized)
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
		it.Event = new(StateStoreInitialized)
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
func (it *StateStoreInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateStoreInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateStoreInitialized represents a Initialized event raised by the StateStore contract.
type StateStoreInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_StateStore *StateStoreFilterer) FilterInitialized(opts *bind.FilterOpts) (*StateStoreInitializedIterator, error) {

	logs, sub, err := _StateStore.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StateStoreInitializedIterator{contract: _StateStore.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_StateStore *StateStoreFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StateStoreInitialized) (event.Subscription, error) {

	logs, sub, err := _StateStore.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateStoreInitialized)
				if err := _StateStore.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_StateStore *StateStoreFilterer) ParseInitialized(log types.Log) (*StateStoreInitialized, error) {
	event := new(StateStoreInitialized)
	if err := _StateStore.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateStoreOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the StateStore contract.
type StateStoreOwnershipTransferredIterator struct {
	Event *StateStoreOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StateStoreOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateStoreOwnershipTransferred)
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
		it.Event = new(StateStoreOwnershipTransferred)
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
func (it *StateStoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateStoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateStoreOwnershipTransferred represents a OwnershipTransferred event raised by the StateStore contract.
type StateStoreOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StateStore *StateStoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StateStoreOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StateStore.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StateStoreOwnershipTransferredIterator{contract: _StateStore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StateStore *StateStoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StateStoreOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StateStore.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateStoreOwnershipTransferred)
				if err := _StateStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StateStore *StateStoreFilterer) ParseOwnershipTransferred(log types.Log) (*StateStoreOwnershipTransferred, error) {
	event := new(StateStoreOwnershipTransferred)
	if err := _StateStore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateStoreStateUpdatedIterator is returned from FilterStateUpdated and is used to iterate over the raw logs and unpacked data for StateUpdated events raised by the StateStore contract.
type StateStoreStateUpdatedIterator struct {
	Event *StateStoreStateUpdated // Event containing the contract specifics and raw log

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
func (it *StateStoreStateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateStoreStateUpdated)
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
		it.Event = new(StateStoreStateUpdated)
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
func (it *StateStoreStateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateStoreStateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateStoreStateUpdated represents a StateUpdated event raised by the StateStore contract.
type StateStoreStateUpdated struct {
	Id        *big.Int
	BlockN    *big.Int
	Timestamp *big.Int
	State     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateUpdated is a free log retrieval operation binding the contract event 0x88aef4d78ad30d12a12a98e96007f5b09c1610b5364b2b99960b7d07e00a8838.
//
// Solidity: event StateUpdated(uint256 id, uint256 blockN, uint256 timestamp, uint256 state)
func (_StateStore *StateStoreFilterer) FilterStateUpdated(opts *bind.FilterOpts) (*StateStoreStateUpdatedIterator, error) {

	logs, sub, err := _StateStore.contract.FilterLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return &StateStoreStateUpdatedIterator{contract: _StateStore.contract, event: "StateUpdated", logs: logs, sub: sub}, nil
}

// WatchStateUpdated is a free log subscription operation binding the contract event 0x88aef4d78ad30d12a12a98e96007f5b09c1610b5364b2b99960b7d07e00a8838.
//
// Solidity: event StateUpdated(uint256 id, uint256 blockN, uint256 timestamp, uint256 state)
func (_StateStore *StateStoreFilterer) WatchStateUpdated(opts *bind.WatchOpts, sink chan<- *StateStoreStateUpdated) (event.Subscription, error) {

	logs, sub, err := _StateStore.contract.WatchLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateStoreStateUpdated)
				if err := _StateStore.contract.UnpackLog(event, "StateUpdated", log); err != nil {
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

// ParseStateUpdated is a log parse operation binding the contract event 0x88aef4d78ad30d12a12a98e96007f5b09c1610b5364b2b99960b7d07e00a8838.
//
// Solidity: event StateUpdated(uint256 id, uint256 blockN, uint256 timestamp, uint256 state)
func (_StateStore *StateStoreFilterer) ParseStateUpdated(log types.Log) (*StateStoreStateUpdated, error) {
	event := new(StateStoreStateUpdated)
	if err := _StateStore.contract.UnpackLog(event, "StateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
