package sdk

import (
	"fmt"
	"math/big"
	"time"

	"github.com/A1exit/dex-sdk/configs"
	"github.com/A1exit/dex-sdk/dex"
	"github.com/A1exit/dex-sdk/factory"
	"github.com/ethereum/go-ethereum/common"
)

type SDK struct {
	config configs.Config
}

func New() (*SDK, error) {
	cfg, err := configs.LoadDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &SDK{config: cfg}, nil
}

func (s *SDK) BuildSwap(pairID string, amountIn *big.Int, recipient string) ([]byte, common.Address, error) {
	if amountIn == nil || amountIn.Sign() <= 0 {
		return nil, common.Address{}, fmt.Errorf("amount must be greater than zero")
	}
	if !common.IsHexAddress(recipient) {
		return nil, common.Address{}, fmt.Errorf("invalid recipient address: %s", recipient)
	}
	pair, err := s.config.GetPair(pairID)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("get pair: %w", err)
	}
	router, err := factory.GetDex(s.config, pair.Dex, pair.Network)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("get router: %w", err)
	}
	params := dex.SwapParams{
		TokenIn:   common.HexToAddress(pair.TokenIn),
		TokenOut:  common.HexToAddress(pair.TokenOut),
		AmountIn:  amountIn,
		Slippage:  pair.Slippage,
		Fee:       pair.Fee,
		Recipient: common.HexToAddress(recipient),
		Deadline:  big.NewInt(time.Now().Add(10 * time.Minute).Unix()),
	}

	calldata, err := router.BuildSwapCallData(params)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("build calldata: %w", err)
	}

	routerAddr, err := s.config.GetRouterAddress(pair.Network, pair.Dex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("get router address: %w", err)
	}

	return calldata, routerAddr, nil
}
