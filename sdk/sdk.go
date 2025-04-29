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
	router dex.Router
	config configs.Config
}

func New(dexType configs.DexType, network configs.Network) (*SDK, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	r, err := factory.GetRouter(dexType, network)
	if err != nil {
		return nil, fmt.Errorf("get router: %w", err)
	}

	return &SDK{
		router: r,
		config: config,
	}, nil
}

func (s *SDK) BuildSwap(pairID string, amountIn *big.Int, recipient string) ([]byte, error) {
	if amountIn == nil || amountIn.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid amountIn: must be greater than 0")
	}
	pair, ok := s.config.Pairs[pairID]
	if !ok {
		return nil, fmt.Errorf("pair not found: %s", pairID)
	}
	if !common.IsHexAddress(recipient) {
		return nil, fmt.Errorf("invalid recipient: %s", recipient)
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

	return s.router.BuildSwapCallData(params)
}
