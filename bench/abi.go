package bench

const RouterABI = `[
	{
		"inputs": [
			{
				"internalType": "address[]",
				"name": "_tests",
				"type": "address[]"
			}
		],
		"name": "append",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "id",
				"type": "uint256"
			},
			{
				"internalType": "uint256[]",
				"name": "opts",
				"type": "uint256[]"
			},
			{
				"internalType": "uint256",
				"name": "times",
				"type": "uint256"
			}
		],
		"name": "operate",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "tests",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

const OperateABI = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_container",
				"type": "address"
			}
		],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [],
		"name": "container",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "counter",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256[]",
				"name": "opts",
				"type": "uint256[]"
			},
			{
				"internalType": "uint256",
				"name": "times",
				"type": "uint256"
			}
		],
		"name": "operate",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`
