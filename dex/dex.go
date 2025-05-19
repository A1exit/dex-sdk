package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var NativeTokenAddress = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")

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

type Router interface {
	BuildSwapCallData(params SwapParams) ([]byte, error)
}
