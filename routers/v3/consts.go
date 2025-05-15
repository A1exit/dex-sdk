package v3

const exactInputABI = `[{
	"inputs": [{
		"components": [
			{"internalType": "bytes", "name": "path", "type": "bytes"},
			{"internalType": "address", "name": "recipient", "type": "address"},
			{"internalType": "uint256", "name": "deadline", "type": "uint256"},
			{"internalType": "uint256", "name": "amountIn", "type": "uint256"},
			{"internalType": "uint256", "name": "amountOutMinimum", "type": "uint256"}
		],
		"internalType": "struct ISwapRouter.ExactInputParams",
		"name": "params",
		"type": "tuple"
	}],
	"name": "exactInput",
	"outputs": [
		{"internalType": "uint256", "name": "amountOut", "type": "uint256"}
	],
	"stateMutability": "payable",
	"type": "function"
}]`

const multicallABI = `[{
	"inputs": [
		{"internalType": "bytes[]", "name": "data", "type": "bytes[]"}
	],
	"name": "multicall", 
	"outputs": [
		{"internalType": "bytes[]", "name": "results", "type": "bytes[]"}
	],
	"stateMutability": "payable",
	"type": "function"
}]`

const unwrapWETH9ABI = `[{
	"inputs": [
		{"internalType": "uint256", "name": "amountMinimum", "type": "uint256"},
		{"internalType": "address", "name": "recipient", "type": "address"}
	],
	"name": "unwrapWETH9",
	"outputs": [],
	"stateMutability": "payable", 
	"type": "function"
}]`
