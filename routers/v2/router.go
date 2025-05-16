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
}

func New(routerAddress common.Address) (*V2Router, error) {
	return &V2Router{
		routerAddress: routerAddress,
	}, nil
}

func (v *V2Router) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	_ = params.Fee

	var parsedABI abi.ABI
	var err error
	var input []byte

	switch {
	case params.TokenIn == dex.NativeTokenAddress:
		parsedABI, err = abi.JSON(bytes.NewReader([]byte(swapExactETHForTokensABI)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse ETH swap ABI: %w", err)
		}
		path := []common.Address{params.WrappedNative, params.TokenOut}
		input, err = parsedABI.Pack(
			"swapExactETHForTokens",
			big.NewInt(0),
			path,
			params.Recipient,
			params.Deadline,
		)

	case params.TokenOut == dex.NativeTokenAddress:
		parsedABI, err = abi.JSON(bytes.NewReader([]byte(SwapExactTokensForETHABI)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse ETH swap ABI: %w", err)
		}
		path := []common.Address{params.TokenIn, params.WrappedNative}
		input, err = parsedABI.Pack(
			"swapExactTokensForETH",
			params.AmountIn,
			big.NewInt(0),
			path,
			params.Recipient,
			params.Deadline,
		)

	default:
		parsedABI, err = abi.JSON(bytes.NewReader([]byte(swapExactTokensForTokensABI)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse token swap ABI: %w", err)
		}
		path := []common.Address{params.TokenIn, params.TokenOut}
		input, err = parsedABI.Pack("swapExactTokensForTokens",
			params.AmountIn,
			big.NewInt(0),
			path,
			params.Recipient,
			params.Deadline)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}
	return input, nil
}
