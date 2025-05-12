package pancakev3

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var _ dex.Router = (*PancakeV3)(nil)

type PancakeV3 struct {
	routerAddress common.Address
	abi           abi.ABI
}

type ExactInputParams struct {
	Path             []byte
	Recipient        common.Address
	Deadline         *big.Int
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
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

func New(routerAddress common.Address) (*PancakeV3, error) {
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(exactInputABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	return &PancakeV3{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}, nil
}

func (p *PancakeV3) Name() string {
	return "pancakev3"
}

func (p *PancakeV3) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	fee := uint32(3000)
	if params.Fee != nil {
		fee = *params.Fee
	}

	path := encodePath(params.TokenIn, params.TokenOut, fee)
	fmt.Println("path:", "0x"+common.Bytes2Hex(path))

	payload := ExactInputParams{
		Path:             path,
		Recipient:        params.Recipient,
		Deadline:         params.Deadline,
		AmountIn:         params.AmountIn,
		AmountOutMinimum: calcAmountOutMin(params.AmountIn, params.Slippage),
	}

	input, err := p.abi.Pack("exactInput", payload)
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
