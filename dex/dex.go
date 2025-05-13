package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Router interface {
	BuildSwapCallData(params SwapParams) ([]byte, error)
}

type SwapParams struct {
	TokenIn   common.Address
	TokenOut  common.Address
	AmountIn  *big.Int
	Slippage  float64
	Fee       *uint32
	Recipient common.Address
	Deadline  *big.Int
	WrappedNative common.Address
}
