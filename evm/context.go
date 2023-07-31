// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"math/big"

	"github.com/axiomesh/eth-kit/ledger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the header corresponding to the hash/number argument pair.
	GetHeader(common.Hash, uint64) *types.Header
}

// NewEVMBlockContext creates a new context for use in the EVM.
func NewEVMBlockContext(number uint64, timestamp uint64, db ledger.StateDB, ledger ledger.ChainLedger, admin string) BlockContext {
	return BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(ledger),
		Coinbase:    common.HexToAddress(admin),
		BlockNumber: new(big.Int).SetUint64(number),
		Time:        timestamp,
		Difficulty:  big.NewInt(0x2000),
		BaseFee:     big.NewInt(0),
		GasLimit:    0x2fefd8,
		Random:      &common.Hash{},
	}
}

// NewEVMTxContext creates a new transaction context for a single transaction.
func NewEVMTxContext(msg *Message) TxContext {
	return TxContext{
		Origin:   msg.From,
		GasPrice: new(big.Int).Set(msg.GasPrice),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ledger ledger.ChainLedger) func(n uint64) common.Hash {
	return func(n uint64) common.Hash {
		hash := ledger.GetBlockHash(n)
		if hash == nil {
			return common.Hash{}
		}
		return common.BytesToHash(hash.Bytes())
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db ledger.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetEVMBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db ledger.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubEVMBalance(sender, amount)
	db.AddEVMBalance(recipient, amount)
}
