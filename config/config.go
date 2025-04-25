package config

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v3"
)

type Network string
type DexType string

const (
	Mainnet   Network = "mainnet"
	Sepolia   Network = "sepolia"
	UniswapV2 DexType = "uniswapv2"
)

type NetworkConfig struct {
	RouterAddress string `yaml:"router_address"`
}

type DexConfig struct {
	ABIPath  string                    `yaml:"abi_path"`
	Networks map[Network]NetworkConfig `yaml:"networks"`
}

type SDKConfig map[DexType]DexConfig

var DefaultConfigPath = "config/config.yaml"

func LoadConfigFromYAML(path string) (SDKConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read yaml: %w", err)
	}
	var cfg SDKConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return cfg, nil
}

func LoadDefaultConfig() (SDKConfig, error) {
	return LoadConfigFromYAML(DefaultConfigPath)
}

func (cfg SDKConfig) GetRouterAddress(dex DexType, net Network) (common.Address, error) {
	dexCfg, ok := cfg[dex]
	if !ok {
		return common.Address{}, fmt.Errorf("domain not found: %s", dex)
	}
	netCfg, ok := dexCfg.Networks[net]
	if !ok {
		return common.Address{}, fmt.Errorf("network not found: %s", net)
	}
	return common.HexToAddress(netCfg.RouterAddress), nil
}

func (cfg SDKConfig) GetABIPath(dex DexType) (string, error) {
	dexCfg, ok := cfg[dex]
	if !ok {
		return "", fmt.Errorf("domain not found: %s", dex)
	}
	return dexCfg.ABIPath, nil
}
