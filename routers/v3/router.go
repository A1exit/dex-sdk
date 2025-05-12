package v3

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var _ dex.Router = (*V3Router)(nil)

type V3Router struct {
	routerAddress common.Address
	abi           abi.ABI
}

// exactInputABI defines the ABI for the exactInput method
const exactInputABI = `[{
	"inputs": [{
		"components": [
			{"internalType": "bytes", "name": "path", "type": "bytes"},
			{"internalType": "address", "name": "recipient", "type": "address"},
			{"internalType": "uint256", "name": "deadline", "type": "uint256"},
			{"internalType": "uint256", "name": "amountIn", "type": "uint256"},
			{"internalType": "uint256", "name": "amountOutMinimum", "type": "uint256"}
		],
		"internalType": "struct ISwapRouter.ExactInputParams",
		"name": "params",
		"type": "tuple"
	}],
	"name": "exactInput",
	"outputs": [
		{"internalType": "uint256", "name": "amountOut", "type": "uint256"}
	],
	"stateMutability": "payable",
	"type": "function"
}]`

func New(routerAddress common.Address) (*V3Router, error) {
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(exactInputABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	return &V3Router{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}, nil
}

func (v *V3Router) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	fee := uint32(3000)
	if params.Fee != nil {
		fee = *params.Fee
	}

	path := encodePath(params.TokenIn, params.TokenOut, fee)
	fmt.Println("path:", "0x"+common.Bytes2Hex(path))

	amountOutMin := big.NewInt(0)

	// Create the params struct
	swapParams := struct {
		Path             []byte
		Recipient        common.Address
		Deadline         *big.Int
		AmountIn         *big.Int
		AmountOutMinimum *big.Int
	}{
		Path:             path,
		Recipient:        params.Recipient,
		Deadline:         params.Deadline,
		AmountIn:         params.AmountIn,
		AmountOutMinimum: amountOutMin,
	}

	input, err := v.abi.Pack("exactInput", swapParams)
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
