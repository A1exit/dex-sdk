package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// NativeTokenAddress represents the native token address (ETH, MATIC, BNB)
var NativeTokenAddress = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")

// SwapParams contains all parameters needed for a swap
type SwapParams struct {
	TokenIn       common.Address
	TokenOut      common.Address
	AmountIn      *big.Int
	Slippage      float64
	Recipient     common.Address
	Deadline      *big.Int
	WrappedNative common.Address
	Fee           *uint32 // Optional fee for V3 swaps
}

// Router interface defines methods for building swap calldata
type Router interface {
	// BuildSwapCallData builds the calldata for a swap transaction
	BuildSwapCallData(params SwapParams) ([]byte, error)
}
