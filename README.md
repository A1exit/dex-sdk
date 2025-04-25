# Dex-sdk

**dex-sdk** is a lightweight Go library for generating calldata for DEX swaps  
(e.g., Uniswap V2, PancakeSwap) — useful for signing and sending transactions in your services or tools.

### Installation

```bash
go get github.com/A1exit/dex-sdk
```

### Quick Example
```
sdkInstance, _ := sdk.New(config.UniswapV2, config.BscTestnet)

params := dex.SwapParams{
    TokenIn:   common.HexToAddress("0x..."),
    TokenOut:  common.HexToAddress("0x..."),
    AmountIn:  big.NewInt(1e18), // 1 token in wei
    Slippage:  0.5,
    Recipient: common.HexToAddress("0x..."),
    Deadline:  big.NewInt(time.Now().Add(10 * time.Minute).Unix()),
    ExactOut:  false,
}

calldata, _ := sdkInstance.BuildSwap(params)
fmt.Printf("Calldata: 0x%x\n", calldata)
```
