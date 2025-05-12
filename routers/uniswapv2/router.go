package uniswapv2

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var _ dex.Router = (*UniV2)(nil)

type UniV2 struct {
	routerAddress common.Address
	abi           abi.ABI
}

// swapExactTokensForTokensABI defines the ABI for the swapExactTokensForTokens method
const swapExactTokensForTokensABI = `[{
	"inputs": [
		{"internalType": "uint256", "name": "amountIn", "type": "uint256"},
		{"internalType": "uint256", "name": "amountOutMin", "type": "uint256"},
		{"internalType": "address[]", "name": "path", "type": "address[]"},
		{"internalType": "address", "name": "to", "type": "address"},
		{"internalType": "uint256", "name": "deadline", "type": "uint256"}
	],
	"name": "swapExactTokensForTokens",
	"outputs": [
		{"internalType": "uint256[]", "name": "amounts", "type": "uint256[]"}
	],
	"stateMutability": "nonpayable",
	"type": "function"
}]`

func New(routerAddress common.Address) (*UniV2, error) {
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(swapExactTokensForTokensABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	return &UniV2{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}, nil
}

func (u *UniV2) Name() string {
	return "uniswapv2"
}

func (u *UniV2) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	_ = params.Fee

	path := []common.Address{params.TokenIn, params.TokenOut}

	slippageMultiplier := big.NewInt(10000 - int64(params.Slippage*10000))
	amountOutMin := new(big.Int).Mul(params.AmountIn, slippageMultiplier)
	amountOutMin.Div(amountOutMin, big.NewInt(10000))

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
