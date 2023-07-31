package ledger

import (
	"github.com/meshplus/bitxhub-kit/types"
)

// BlockchainLedger handles block, transaction and receipt data.
//
//go:generate mockgen -destination mock_ledger/mock_chain_ledger.go -package mock_ledger -source chain_ledger.go
type ChainLedger interface {
	// GetBlock get block with height
	GetBlock(height uint64) (*types.Block, error)

	// GetBlockSign get the signature of block
	GetBlockSign(height uint64) ([]byte, error)

	// GetBlockByHash get the block using block hash
	GetBlockByHash(hash *types.Hash) (*types.Block, error)

	// GetTransaction get the transaction using transaction hash
	GetTransaction(hash *types.Hash) (*types.Transaction, error)

	// GetTransactionMeta get the transaction meta data
	GetTransactionMeta(hash *types.Hash) (*types.TransactionMeta, error)

	// GetReceipt get the transaction receipt
	GetReceipt(hash *types.Hash) (*types.Receipt, error)

	// PersistExecutionResult persist the execution result
	PersistExecutionResult(block *types.Block, receipts []*types.Receipt) error

	// GetChainMeta get chain meta data
	GetChainMeta() *types.ChainMeta

	// UpdateChainMeta update the chain meta data
	UpdateChainMeta(*types.ChainMeta)

	// LoadChainMeta get chain meta data
	LoadChainMeta() (*types.ChainMeta, error)

	// GetTxCountInBlock get the transaction count in a block
	GetTransactionCount(height uint64) (uint64, error)

	RollbackBlockChain(height uint64) error

	GetBlockHash(height uint64) *types.Hash

	Close()
}
