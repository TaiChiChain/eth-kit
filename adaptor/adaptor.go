package adaptor

import (
	"errors"
	"math/big"

	"github.com/axiomesh/axiom-kit/types"
	vm "github.com/axiomesh/eth-kit/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func CallArgsToMessage(args *types.CallArgs, globalGasCap uint64, baseFee *big.Int) (*vm.Message, error) {
	// Reject invalid combinations of pre- and post-1559 fee styles
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// Set sender address or use zero address if none specified.
	addr := args.GetFrom()

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = uint64(math.MaxUint64 / 2)
	}
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		// log.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}
	var (
		gasPrice  *big.Int
		gasFeeCap *big.Int
		gasTipCap *big.Int
	)
	if baseFee == nil {
		// If there's no basefee, then it must be a non-1559 execution
		gasPrice = new(big.Int)
		if args.GasPrice != nil {
			gasPrice = args.GasPrice.ToInt()
		}
		gasFeeCap, gasTipCap = gasPrice, gasPrice
	} else {
		// A basefee is provided, necessitating 1559-type execution
		if args.GasPrice != nil {
			// User specified the legacy gas field, convert to 1559 gas typing
			gasPrice = args.GasPrice.ToInt()
			gasFeeCap, gasTipCap = gasPrice, gasPrice
		} else {
			// User specified 1559 gas fields (or none), use those
			gasFeeCap = new(big.Int)
			if args.MaxFeePerGas != nil {
				gasFeeCap = args.MaxFeePerGas.ToInt()
			}
			gasTipCap = new(big.Int)
			if args.MaxPriorityFeePerGas != nil {
				gasTipCap = args.MaxPriorityFeePerGas.ToInt()
			}
			// Backfill the legacy gasPrice for EVM execution, unless we're all zeroes
			gasPrice = new(big.Int)
			if gasFeeCap.BitLen() > 0 || gasTipCap.BitLen() > 0 {
				gasPrice = math.BigMin(new(big.Int).Add(gasTipCap, baseFee), gasFeeCap)
			}
		}
	}
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}
	data := args.GetData()
	var accessList ethtypes.AccessList
	if args.AccessList != nil {
		accessList = *args.AccessList
	}
	msg := &vm.Message{
		From:              addr,
		To:                args.To,
		Value:             value,
		GasLimit:          gas,
		GasPrice:          gasPrice,
		GasFeeCap:         gasFeeCap,
		GasTipCap:         gasTipCap,
		Data:              data,
		AccessList:        accessList,
		SkipAccountChecks: true,
	}
	return msg, nil
}

func TransactionToMessage(tx *types.Transaction) *vm.Message {
	from := common.BytesToAddress(tx.GetFrom().Bytes())
	var to *common.Address
	if tx.GetTo() != nil {
		toAddr := common.BytesToAddress(tx.GetTo().Bytes())
		to = &toAddr
	}

	isFake := false
	if v, _, _ := tx.GetRawSignature(); v == nil {
		isFake = true
	}

	msg := &vm.Message{
		Nonce:             tx.GetNonce(),
		GasLimit:          tx.GetGas(),
		GasPrice:          new(big.Int).Set(tx.GetGasPrice()),
		GasFeeCap:         new(big.Int).Set(tx.GetGasFeeCap()),
		GasTipCap:         new(big.Int).Set(tx.GetGasTipCap()),
		From:              from,
		To:                to,
		Value:             tx.GetValue(),
		Data:              tx.GetPayload(),
		AccessList:        tx.GetInner().GetAccessList(),
		SkipAccountChecks: isFake,
	}
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	// if baseFee != nil {
	// 	msg.GasPrice = cmath.BigMin(msg.GasPrice.Add(msg.GasTipCap, baseFee), msg.GasFeeCap)
	// }
	return msg
}
