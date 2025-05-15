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
}

func New(routerAddress common.Address) (*V3Router, error) {
	return &V3Router{
		routerAddress: routerAddress,
	}, nil
}

func (v *V3Router) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	var input []byte
	var err error
	var parsedABI abi.ABI

	var tokenIn common.Address
	var tokenOut common.Address
	var recipient common.Address

	fee := uint32(3000)
	if params.Fee != nil {
		fee = *params.Fee
	}

	switch {
	case params.TokenIn == dex.NativeTokenAddress:
		tokenIn = params.WrappedNative
		tokenOut = params.TokenOut
		recipient = params.Recipient
	case params.TokenOut == dex.NativeTokenAddress:
		tokenIn = params.TokenIn
		tokenOut = params.WrappedNative
		recipient = v.routerAddress
	default:
		tokenIn = params.TokenIn
		tokenOut = params.TokenOut
		recipient = params.Recipient
	}

	path := encodePath(tokenIn, tokenOut, fee)
	fmt.Println("path:", "0x"+common.Bytes2Hex(path))
	amountOutMin := big.NewInt(0)

	parsedABI, err = abi.JSON(bytes.NewReader([]byte(exactInputABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ETH swap ABI: %w", err)
	}

	swapParams := struct {
		Path             []byte
		Recipient        common.Address
		Deadline         *big.Int
		AmountIn         *big.Int
		AmountOutMinimum *big.Int
	}{
		Path:             path,
		Recipient:        recipient,
		Deadline:         params.Deadline,
		AmountIn:         params.AmountIn,
		AmountOutMinimum: amountOutMin,
	}

	input, err = parsedABI.Pack("exactInput", swapParams)
	fmt.Println("input:", "0x"+common.Bytes2Hex(input))

	if err != nil {
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}

	if tokenOut != params.WrappedNative {
		return input, nil
	} else {

		unwrapABI, err := abi.JSON(bytes.NewReader([]byte(unwrapWETH9ABI)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse unwrapWETH9 ABI: %w", err)
		}
		unwrapInput, err := unwrapABI.Pack("unwrapWETH9", amountOutMin, params.Recipient)
		fmt.Println("unwrapInput:", "0x"+common.Bytes2Hex(unwrapInput))
		if err != nil {
			return nil, fmt.Errorf("failed to pack unwrapWETH9 calldata: %w", err)
		}
		multicallABI, err := abi.JSON(bytes.NewReader([]byte(multicallABI)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse multicall ABI: %w", err)
		}
		multicallInput, err := multicallABI.Pack("multicall", [][]byte{input, unwrapInput})
		if err != nil {
			return nil, fmt.Errorf("failed to pack multicall calldata: %w", err)
		}

		return multicallInput, nil
	}
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
