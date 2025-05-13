package factory

import (
	"fmt"

	"github.com/A1exit/dex-sdk/configs"
	"github.com/A1exit/dex-sdk/dex"
	v2 "github.com/A1exit/dex-sdk/routers/v2"
	v3 "github.com/A1exit/dex-sdk/routers/v3"
)

func GetDex(config configs.Config, dexType configs.DexType, net configs.Network) (dex.Router, error) {
	routerAddr, err := config.GetRouterAddress(net, dexType)
	if err != nil {
		return nil, fmt.Errorf("get router address: %w", err)
	}
	var router dex.Router
	switch dexType {
	case configs.UniswapV2, configs.PancakeV2:
		router, err = v2.New(routerAddr)
	case configs.UniswapV3, configs.PancakeV3:
		router, err = v3.New(routerAddr)
	default:
		return nil, fmt.Errorf("unsupported dex type: %s", dexType)
	}

	if err != nil {
		return nil, fmt.Errorf("create %s router: %w", dexType, err)
	}

	return router, nil
}
