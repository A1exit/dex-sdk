package sdk

import (
	"github.com/A1exit/dex-sdk/config"
	"github.com/A1exit/dex-sdk/dex"
	"github.com/A1exit/dex-sdk/factory"
)

type SDK struct {
	dex dex.Dex
}

func New(dexType config.DexType, net config.Network) (*SDK, error) {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		return nil, err
	}

	d, err := factory.GetDex(cfg, dexType, net)
	if err != nil {
		return nil, err
	}

	return &SDK{dex: d}, nil
}

func (s *SDK) BuildSwap(params dex.SwapParams) ([]byte, error) {
	return s.dex.BuildSwapCallData(params)
}
