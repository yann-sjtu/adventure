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
)

func EthBlockNumber() Request {
	return CreateRequest(ethBlockNumber, nil)
}
func EthGetBalance() Request {
	return CreateRequest(ethGetBalance, []interface{}{queryAccount[rand.Intn(len(queryAccount))], latestBlockNumber})
}
func EthGetBlockByNumber() Request {
	rand.Seed(time.Now().UnixNano()+rand.Int63n(11111111111111))
	height := startHeight+rand.Int63n(120000)
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
