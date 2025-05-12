package factory

import (
	"fmt"

	"github.com/A1exit/dex-sdk/routers/pancakev3"

	"github.com/A1exit/dex-sdk/configs"
	"github.com/A1exit/dex-sdk/dex"
	"github.com/A1exit/dex-sdk/routers/pancakev2"
	"github.com/A1exit/dex-sdk/routers/uniswapv2"
	"github.com/A1exit/dex-sdk/routers/uniswapv3"
)

func GetDex(config configs.Config, dexType configs.DexType, net configs.Network) (dex.Router, error) {
	routerAddr, err := config.GetRouterAddress(net, dexType)
	if err != nil {
		return nil, fmt.Errorf("get router address: %w", err)
	}

	var router dex.Router
	switch dexType {
	case configs.UniswapV2:
		router, err = uniswapv2.New(routerAddr)
	case configs.UniswapV3:
		router, err = uniswapv3.New(routerAddr)
	case configs.PancakeV3:
		router, err = pancakev3.New(routerAddr)
	case configs.PancakeV2:
		router, err = pancakev2.New(routerAddr)
	default:
		return nil, fmt.Errorf("unsupported dex type: %s", dexType)
	}

	if err != nil {
		return nil, fmt.Errorf("create %s router: %w", dexType, err)
	}

	return router, nil
}
