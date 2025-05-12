package v2

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var _ dex.Router = (*V2Router)(nil)

type V2Router struct {
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

func New(routerAddress common.Address) (*V2Router, error) {
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(swapExactTokensForTokensABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	return &V2Router{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}, nil
}

func (v *V2Router) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	_ = params.Fee

	path := []common.Address{params.TokenIn, params.TokenOut}

	input, err := v.abi.Pack("swapExactTokensForTokens",
		params.AmountIn,
		big.NewInt(0),
		path,
		params.Recipient,
		params.Deadline,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}

	return input, nil
}
