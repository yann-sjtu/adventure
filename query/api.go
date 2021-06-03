package query

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	ethBlockNumber      = "eth_blockNumber"
	ethGetBalance       = "eth_getBalance"
	ethGetBlockByNumber = "eth_getBlockByNumber"

	ethGasPrice              = "eth_gasPrice"
	ethGetCode               = "eth_getCode"
	ethGetTransactionCount   = "eth_getTransactionCount"
	ethGetTransactionReceipt = "eth_getTransactionReceipt"

	netVersion = "net_version"

	ethCall = "eth_call"
)

var (
	persistentBlockNumberRequest = EthBlockNumber()
	persistentGasPriceRequest    = EthGasPrice()
	persistentGetCodeReuqest     = EthGetCode()
)

func EthBlockNumber() Request {
	return CreateRequest(ethBlockNumber, nil)
}
func EthGetBalance() Request {
	return CreateRequest(ethGetBalance, []interface{}{queryAccount[rand.Intn(len(queryAccount))], latestBlockNumber})
}
func EthGetBlockByNumber() Request {
	rand.Seed(time.Now().UnixNano()+rand.Int63n(11111111111111))
	height := startHeight+rand.Int63n(402925)
	return CreateRequest(ethGetBlockByNumber, []interface{}{hexutil.Uint64(height), false})
}
func EthGasPrice() Request {
	return CreateRequest(ethGasPrice, nil)
}
func EthGetCode() Request {
	return CreateRequest(ethGetCode, []interface{}{usdtContract, latestBlockNumber})
}
func EthGetTransactionCount() Request {
	return CreateRequest(ethGetTransactionCount, []interface{}{queryAccount[rand.Intn(len(queryAccount))], latestBlockNumber})
}
func EthGetTransactionReceipt() Request {
	return CreateRequest(ethGetTransactionReceipt, []interface{}{txHashList[rand.Intn(len(txHashList))]})
}

func NetVersion() Request {
	return CreateRequest(netVersion, nil)
}

func EthCall() Request {
	callArgs := generateCallArgs()

	return CreateRequest(ethCall, []interface{}{callArgs, latestBlockNumber})
}

func generateCallArgs() map[string]string {
	callArgs := make(map[string]string)

	i := rand.Intn(100)%6
	switch i {
	case 0:
		callArgs["from"] = "0xd01bf1F0dB0E0F9998Ec01a45Cfc03116D0224bE"
		callArgs["to"] = "0x4a191907012673c9efde02a10a24c19db48bed0c"
		callArgs["value"] = "0x0"
		callArgs["data"] = "0x70a08231000000000000000000000000"+hexAddrs[rand.Intn(len(hexAddrs))][2:]
		callArgs["gas"] = "0x2dc6c0"
	case 1:
		callArgs["from"] = "0x5627f50785AE5EcE283FF79635628d505f1e7d3F"
		callArgs["to"] = "0xf3188861838F6B6b0a11f7619b2C0E00CAFf28b5"
		callArgs["data"] = "0x82821e15"
	case 2:
		callArgs["from"] = "0xf871EB07dA5DaE550731BC7Eb3aE7888253f28ac"
		callArgs["to"] = "0xf3188861838F6B6b0a11f7619b2C0E00CAFf28b5"
		callArgs["data"] = "0x8b19b75b"
	case 3:
		callArgs["to"] = "0x6f63aeca1250CfC9006bc0807204C42f2cb68a87"
		callArgs["data"] = "0x252dba42000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000003c00000000000000000000000000000000000000000000000000000000000000460000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000005a0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000006e00000000000000000000000000000000000000000000000000000000000000780000000000000000000000000000000000000000000000000000000000000082000000000000000000000000000000000000000000000000000000000000008c000000000000000000000000000000000000000000000000000000000000009600000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000aa00000000000000000000000000000000000000000000000000000000000000b400000000000000000000000000000000000000000000000000000000000000be00000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000000d200000000000000000000000000000000000000000000000000000000000000dc00000000000000000000000000000000000000000000000000000000000000e60000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341db00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341dc00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341dd00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341de00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341df00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e000000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e100000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e200000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e300000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e400000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e500000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e600000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e700000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e800000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341e900000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341ea00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341eb00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341ec00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341ed00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000245f8e26a700000000000000000000000000000000000000000000000000000000000341ee00000000000000000000000000000000000000000000000000000000"
	case 4:
		callArgs["to"] = "0x6f63aeca1250CfC9006bc0807204C42f2cb68a87"
		callArgs["data"] = "0x252dba4200000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000028000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000005a0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000006e00000000000000000000000000000000000000000000000000000000000000780000000000000000000000000000000000000000000000000000000000000082000000000000000000000000000000000000000000000000000000000000008c000000000000000000000000000000000000000000000000000000000000009600000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000aa00000000000000000000000000000000000000000000000000000000000000b400000000000000000000000000000000000000000000000000000000000000be00000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000000d200000000000000000000000000000000000000000000000000000000000000dc00000000000000000000000000000000000000000000000000000000000000e600000000000000000000000000000000000000000000000000000000000000f000000000000000000000000000000000000000000000000000000000000000fa0000000000000000000000000000000000000000000000000000000000000104000000000000000000000000000000000000000000000000000000000000010e00000000000000000000000000000000000000000000000000000000000001180000000000000000000000000000000000000000000000000000000000000122000000000000000000000000000000000000000000000000000000000000012c00000000000000000000000000000000000000000000000000000000000001360000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000014a0000000000000000000000000000000000000000000000000000000000000154000000000000000000000000000000000000000000000000000000000000015e00000000000000000000000000000000000000000000000000000000000001680000000000000000000000000000000000000000000000000000000000000172000000000000000000000000000000000000000000000000000000000000017c00000000000000000000000000000000000000000000000000000000000001860000000000000000000000000000000000000000000000000000000000000190000000000000000000000000000000000000000000000000000000000000019a00000000000000000000000000000000000000000000000000000000000001a400000000000000000000000000000000000000000000000000000000000001ae00000000000000000000000000000000000000000000000000000000000001b800000000000000000000000000000000000000000000000000000000000001c200000000000000000000000000000000000000000000000000000000000001cc00000000000000000000000000000000000000000000000000000000000001d60000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001303e00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001303f00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001304000000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306000000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306100000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306200000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306300000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306400000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306500000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000001306600000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022100000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022200000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022300000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022400000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022500000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022600000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022700000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022800000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022900000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022a00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022b00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022c00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022d00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022e00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002022f00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023000000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023100000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023200000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023300000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023400000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023500000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023600000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023700000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023800000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023900000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023a00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023b00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023c00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023d00000000000000000000000000000000000000000000000000000000000000000000000000000000a1341ad9250c44cd9bfbb5e050cee8f92e1941b8000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000241bf1e00f000000000000000000000000000000000000000000000000000000000002023e00000000000000000000000000000000000000000000000000000000"
	case 5:
		callArgs["to"] = "0xd9f780e7A9692b7De92aBdBb501714Dfe035A86D"
		callArgs["data"] = "0x7c7be7d1000000000000000000000000a9c809db065abf9973f0ab0dc7a74d68fea590f9000000000000000000000000000000000000000000000000000000000000020d"
	default:
	}
	return callArgs
}