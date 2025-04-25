package uniswapv2

import (
	"bytes"
	"fmt"
	"github.com/A1exit/dex-sdk/dex"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type UniV2 struct {
	routerAddress common.Address
	abi           abi.ABI
}

func New(routerAddress common.Address, abiPath string) *UniV2 {
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		panic(fmt.Errorf("failed to read ABI file: %w", err))
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		panic(fmt.Errorf("failed to parse ABI: %w", err))
	}

	return &UniV2{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}
}

func (u *UniV2) Name() string {
	return "uniswapv2"
}

func (u *UniV2) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	path := []common.Address{params.TokenIn, params.TokenOut}

	slippage := big.NewInt(int64(params.Slippage * 10000))
	one := big.NewInt(10000)
	amountOutMin := new(big.Int).Mul(params.AmountIn, new(big.Int).Sub(one, slippage))
	amountOutMin.Div(amountOutMin, one)

	input, err := u.abi.Pack("swapExactTokensForTokens",
		params.AmountIn,
		amountOutMin,
		path,
		params.Recipient,
		params.Deadline,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}

	return input, nil
}
