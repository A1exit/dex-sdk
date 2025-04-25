package factory

import (
	"fmt"
	"github.com/A1exit/dex-sdk/dex"
	"github.com/A1exit/dex-sdk/uniswapv2"

	"github.com/A1exit/dex-sdk/config"
)

func GetDex(cfg config.SDKConfig, dexType config.DexType, net config.Network) (dex.Dex, error) {
	switch dexType {
	case config.UniswapV2:
		routerAddr, err := cfg.GetRouterAddress(dexType, net)
		if err != nil {
			return nil, fmt.Errorf("get router address: %w", err)
		}
		abiPath, err := cfg.GetABIPath(dexType)
		if err != nil {
			return nil, fmt.Errorf("get abi path: %w", err)
		}
		return uniswapv2.New(routerAddr, abiPath), nil
	default:
		return nil, fmt.Errorf("unsupported domain: %s", dexType)
	}
}
