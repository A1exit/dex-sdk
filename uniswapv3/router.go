package uniswapv3

import (
	"bytes"
	"fmt"
	"math/big"
	"os"

	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const abiPath = "uniswapv3/abi/UniswapV3Router.abi.json"

var _ dex.Router = (*UniV3)(nil)

type UniV3 struct {
	routerAddress common.Address
	abi           abi.ABI
}

func New(routerAddress common.Address) *UniV3 {
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		panic(fmt.Errorf("failed to read ABI file: %w", err))
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		panic(fmt.Errorf("failed to parse ABI: %w", err))
	}

	return &UniV3{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}
}

func (u *UniV3) Name() string {
	return "uniswapv3"
}

func (u *UniV3) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	fee := uint32(3000) // default
	if params.Fee != nil {
		fee = *params.Fee
	}

	path := encodePath(params.TokenIn, params.TokenOut, fee)

	payload := map[string]interface{}{
		"path":             path,
		"recipient":        params.Recipient,
		"deadline":         params.Deadline,
		"amountIn":         params.AmountIn,
		"amountOutMinimum": calcAmountOutMin(params.AmountIn, params.Slippage),
	}

	input, err := u.abi.Pack("exactInput", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}

	return input, nil
}

func encodePath(tokenIn, tokenOut common.Address, fee uint32) []byte {
	var path []byte
	path = append(path, tokenIn.Bytes()...)
	path = append(path, uint24ToBytes(fee)...)
	path = append(path, tokenOut.Bytes()...)
	return path
}

func uint24ToBytes(v uint32) []byte {
	return []byte{byte(v >> 16), byte(v >> 8), byte(v)}
}

func calcAmountOutMin(amountIn *big.Int, slippage float64) *big.Int {
	slippageFactor := big.NewFloat(1.0 - slippage/100)
	amount := new(big.Float).SetInt(amountIn)
	amount.Mul(amount, slippageFactor)

	result := new(big.Int)
	amount.Int(result)
	return result
}
