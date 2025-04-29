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
	dex    dex.Dex
	config configs.Config
}

func New(dexType configs.DexType, network configs.Network) (*SDK, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	d, err := factory.GetDex(dexType, network)
	if err != nil {
		return nil, fmt.Errorf("get dex: %w", err)
	}

	return &SDK{
		dex:    d,
		config: config,
	}, nil
}

func (s *SDK) BuildSwap(pairID string, amountIn *big.Int, recipient common.Address) ([]byte, error) {
	pair, ok := s.config.Pairs[pairID]
	if !ok {
		return nil, fmt.Errorf("pair not found: %s", pairID)
	}

	params := dex.SwapParams{
		TokenIn:   common.HexToAddress(pair.TokenIn),
		TokenOut:  common.HexToAddress(pair.TokenOut),
		AmountIn:  amountIn,
		Slippage:  pair.Slippage,
		Recipient: recipient,
		Deadline:  big.NewInt(time.Now().Add(10 * time.Minute).Unix()),
	}

	return s.dex.BuildSwapCallData(params)
}
