# Dex-sdk

**dex-sdk** is a lightweight Go library for generating calldata for DEX swaps  
(e.g., Uniswap V2/V3, PancakeSwap V2/V3, QuickSwap V3) — useful for signing and sending transactions in your services or tools.

### Installation

```bash
go get github.com/A1exit/dex-sdk
```

### Supported DEXes
- Uniswap V2/V3
- PancakeSwap V2/V3
- QuickSwap V3

### Quick Example
```go
package main

import (
    "fmt"
    "log"
    "math/big"

    "github.com/A1exit/dex-sdk/sdk"
)

func main() {
    // Create SDK instance
    sdkInstance, err := sdk.New()
    if err != nil {
        log.Fatalf("failed to create sdk instance: %v", err)
    }

    // Prepare swap parameters
    pairID := "polygon-token2-eth-v3"  // pair ID from your config
    amount := big.NewInt(1e+13)        // amount to swap
    recipient := "0xCcD2466613fe5D185E5FA081A0B71040277606dd"  // recipient address

    // Generate swap calldata
    calldata, routerAddr, err := sdkInstance.BuildSwap(pairID, amount, recipient)
    if err != nil {
        log.Fatalf("failed to build swap calldata: %v", err)
    }

    fmt.Printf("Calldata: 0x%x\n", calldata)
    fmt.Printf("Router: %s\n", routerAddr.String())
}

### Configuration
Create a `config.yaml` file in the `configs` directory:

```yaml
networks:
  polygon:
    routers:
      quickswapv3: "0x..."
    wrapped_native: "0x..."
  bsc:
    routers:
      pancakev2: "0x..."
      pancakev3: "0x..."
    wrapped_native: "0x..."

pairs:
  polygon-token1-token2-v3:
    dex: quickswapv3
    network: polygon
    token_in: "0x..."
    token_out: "0x..."
    slippage: 0.5
    fee: 3000  # optional for V3
```

### Features
- Support for multiple DEXes (Uniswap V2/V3, PancakeSwap V2/V3, QuickSwap V3)
- Native token support (ETH, MATIC, BNB)
- Simple API for generating swap calldata
- Easy integration with any Ethereum-compatible network
