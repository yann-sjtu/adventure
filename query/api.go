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
	rand.Seed(time.Now().UnixNano() + rand.Int63n(11111111111111))
	height := startHeight + rand.Int63n(2000)
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
	callArgs := make(map[string]string)
	callArgs["from"] = "0xd01bf1F0dB0E0F9998Ec01a45Cfc03116D0224bE"
	callArgs["to"] = "0x4a191907012673c9efde02a10a24c19db48bed0c"
	callArgs["value"] = "0x0"
	callArgs["data"] = "0x70a08231000000000000000000000000" + hexAddrs[rand.Intn(len(hexAddrs))][2:]
	callArgs["gas"] = "0x2dc6c0"

	return CreateRequest(ethCall, []interface{}{callArgs, latestBlockNumber})
}
