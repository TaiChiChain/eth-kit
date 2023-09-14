package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateEVMAccount(common.Address)

	SubEVMBalance(common.Address, *big.Int)
	AddEVMBalance(common.Address, *big.Int)
	GetEVMBalance(common.Address) *big.Int

	GetEVMNonce(common.Address) uint64
	SetEVMNonce(common.Address, uint64)

	GetEVMCodeHash(common.Address) common.Hash
	GetEVMCode(common.Address) []byte
	SetEVMCode(common.Address, []byte)
	GetEVMCodeSize(common.Address) int

	AddEVMRefund(uint64)
	SubEVMRefund(uint64)
	GetEVMRefund() uint64

	GetEVMCommittedState(common.Address, common.Hash) common.Hash
	GetEVMState(common.Address, common.Hash) common.Hash
	SetEVMState(common.Address, common.Hash, common.Hash)

	GetEVMTransientState(addr common.Address, key common.Hash) common.Hash
	SetEVMTransientState(addr common.Address, key, value common.Hash)

	SuicideEVM(common.Address) bool
	HasSuicideEVM(common.Address) bool

	// Exist reports whether the given Account exists in state.
	// Notably this should also return true for suicided accounts.
	ExistEVM(common.Address) bool
	// Empty returns whether the given Account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	EmptyEVM(common.Address) bool

	// PrepareEVMAccessList(sender common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)
	AddressInEVMAccessList(addr common.Address) bool
	SlotInEVMAceessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
	// AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddAddressToEVMAccessList(addr common.Address)
	// AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddSlotToEVMAccessList(addr common.Address, slot common.Hash)
	PrepareEVM(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)

	RevertToSnapshot(int)
	Snapshot() int

	AddEVMLog(log *types.Log)
	AddEVMPreimage(common.Hash, []byte)
}
