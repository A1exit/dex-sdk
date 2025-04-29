package configs

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v3"
)

type DexType string
type Network string

const (
	UniswapV2 DexType = "uniswapv2"
	UniswapV3 DexType = "uniswapv3"

	Mainnet  Network = "mainnet"
	Sepolia  Network = "sepolia"
	BSC      Network = "bsc"
	Arbitrum Network = "arbitrum"
)

type NetworkRouters struct {
	Routers map[DexType]string `yaml:"routers"`
}

type PairConfig struct {
	Dex      DexType `yaml:"dex"`
	Network  Network `yaml:"network"`
	TokenIn  string  `yaml:"token_in"`
	TokenOut string  `yaml:"token_out"`
	Slippage float64 `yaml:"slippage"`
	Fee      *uint32 `yaml:"fee,omitempty"`
}

type Config struct {
	Networks map[Network]NetworkRouters `yaml:"networks"`
	Pairs    map[string]PairConfig      `yaml:"pairs"`
}

var DefaultConfigPath = "configs/config.yaml"
var cachedConfig *Config

func LoadConfig() (Config, error) {
	if cachedConfig != nil {
		return *cachedConfig, nil
	}
	data, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("read config yaml: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config yaml: %w", err)
	}
	return cfg, nil
}

func (cfg Config) GetRouterAddress(net Network, dex DexType) (common.Address, error) {
	netCfg, ok := cfg.Networks[net]
	if !ok {
		return common.Address{}, fmt.Errorf("network not found: %s", net)
	}
	addrStr, ok := netCfg.Routers[dex]
	if !ok {
		return common.Address{}, fmt.Errorf("router for dex %s not found in network %s", dex, net)
	}
	return common.HexToAddress(addrStr), nil
}

func (cfg Config) GetPair(pairID string) (PairConfig, error) {
	pair, ok := cfg.Pairs[pairID]
	if !ok {
		return PairConfig{}, fmt.Errorf("pair not found: %s", pairID)
	}
	return pair, nil
}
