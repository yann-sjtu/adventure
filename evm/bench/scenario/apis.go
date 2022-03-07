package scenario

import (
	"encoding/json"
	"net/http"
)

const (
	ethBlockNumber      		= "eth_blockNumber"
	ethGetBalance       		= "eth_getBalance"
	ethGetBlockByNumber 		= "eth_getBlockByNumber"
	ethGasPrice              	= "eth_gasPrice"
	ethGetCode               	= "eth_getCode"
	ethGetTransactionCount   	= "eth_getTransactionCount"
	ethGetTransactionReceipt 	= "eth_getTransactionReceipt"
	netVersion 					= "net_version"
	ethCall 					= "eth_call"
)

const (
	jsonrpc						= "2.0"
	id							= 1
)

type ReqBody struct {
	jsonrpc  	string
	method 		string
	params		interface{}
	id			int
}

func NewReqBody(jsonrpc string, method string, params interface{}, id int) *ReqBody {
	return &ReqBody{
		jsonrpc: jsonrpc,
		method: method,
		params: params,
		id: 	id,
	}
}

//获取当前区块高度
func EthBlockNumberApi(url string) (*http.Response, error) {
	method := ethBlockNumber
	request := NewReqBody(jsonrpc,method,nil, id)
	//处理成为[]byte类型的req
	req, err := json.Marshal(*request)
	if err != nil {
		panic(err)
	}
	return DoPost(url, req)
}

//获取余额
func EthGetBalanceApi(url string, params interface{}) (*http.Response, error) {
	method := ethGetBalance
	request := NewReqBody(jsonrpc, method, params, id)
	req, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	return DoPost(url, req)
}

