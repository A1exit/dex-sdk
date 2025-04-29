package factory

import (
	"fmt"

	"github.com/A1exit/dex-sdk/configs"
	"github.com/A1exit/dex-sdk/dex"
	"github.com/A1exit/dex-sdk/uniswapv2"
	"github.com/A1exit/dex-sdk/uniswapv3"
)

func GetDex(dexType configs.DexType, net configs.Network) (dex.Dex, error) {
	networksConfig, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load networks config: %w", err)
	}

	routerAddr, err := networksConfig.GetRouterAddress(net, dexType)
	if err != nil {
		return nil, fmt.Errorf("get router address: %w", err)
	}

	switch dexType {
	case configs.UniswapV2:
		return uniswapv2.New(routerAddr), nil
	case configs.UniswapV3:
		return uniswapv3.New(routerAddr), nil
	default:
		return nil, fmt.Errorf("unsupported dex type: %s", dexType)
	}
}
