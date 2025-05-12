package pancakev2

import (
	"bytes"
	"fmt"
	"math/big"
	"os"

	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const abiPath = "routers/pancakev2/abi/PancakeV2Router.abi.json"

var _ dex.Router = (*PancakeV2)(nil)

type PancakeV2 struct {
	routerAddress common.Address
	abi           abi.ABI
}

func New(routerAddress common.Address) (*PancakeV2, error) {
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %w", err)
	}
	parsedABI, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	return &PancakeV2{
		routerAddress: routerAddress,
		abi:           parsedABI,
	}, nil
}

func (p *PancakeV2) Name() string {
	return "pancakev2"
}

func (p *PancakeV2) BuildSwapCallData(params dex.SwapParams) ([]byte, error) {
	_ = params.Fee

	path := []common.Address{params.TokenIn, params.TokenOut}

	slippageMultiplier := big.NewInt(10000 - int64(params.Slippage*10000))
	amountOutMin := new(big.Int).Mul(params.AmountIn, slippageMultiplier)
	amountOutMin.Div(amountOutMin, big.NewInt(10000))

	input, err := p.abi.Pack("swapExactTokensForTokens",
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
