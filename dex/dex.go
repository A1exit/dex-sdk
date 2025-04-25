package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Dex interface {
	Name() string
	BuildSwapCallData(params SwapParams) ([]byte, error)
}

type SwapParams struct {
	TokenIn   common.Address
	TokenOut  common.Address
	AmountIn  *big.Int
	Slippage  float64
	Recipient common.Address
	Deadline  *big.Int
}
