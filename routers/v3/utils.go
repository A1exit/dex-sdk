package v3

import (
	"github.com/A1exit/dex-sdk/dex"
	"github.com/ethereum/go-ethereum/common"
)

func EncodePath(tokenIn, tokenOut common.Address, fee *uint32) []byte {
	var path []byte
	path = append(path, tokenIn.Bytes()...)
	if fee != nil {
		path = append(path, Uint24ToBytes(*fee)...)
	}
	path = append(path, tokenOut.Bytes()...)
	return path
}

func Uint24ToBytes(v uint32) []byte {
	return []byte{byte(v >> 16), byte(v >> 8), byte(v)}
}

func GetSwapTokens(params dex.SwapParams, routerAddress common.Address) (tokenIn, tokenOut, recipient common.Address) {
	switch {
	case params.TokenIn == dex.NativeTokenAddress:
		tokenIn = params.WrappedNative
		tokenOut = params.TokenOut
		recipient = params.Recipient
	case params.TokenOut == dex.NativeTokenAddress:
		tokenIn = params.TokenIn
		tokenOut = params.WrappedNative
		recipient = routerAddress
	default:
		tokenIn = params.TokenIn
		tokenOut = params.TokenOut
		recipient = params.Recipient
	}
	return
}
