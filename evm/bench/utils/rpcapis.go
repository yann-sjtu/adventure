package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	//交易操作类
	EthSendTransaction = "eth_sendTransaction"
	EthSendRawTransaction = "eth_sendRawTransaction"

	//当前状态类
	EthSyncing = "eth_syncing"
	EthGasPrice = "eth_gasPrice"
	EthAccounts = "eth_accounts"
	EthBlockNumber = "eth_blockNumber"
	EthGetBalance = "eth_getBalance"
	EthGetStorageAt = "eth_getStorageAt"
	EthCall = "eth_call"
	EthEstimateGas = "eth_estimateGas"
	EthNewFilter = "eth_newFilter"
	EthNewBlockFilter = "eth_newBlockFilter"
	EthNewPendingTransactionFilter = "eth_newPendingTransactionFilter"
	EthUninstallFilter = "eth_uninstallFilter"
	EthSubscribe = "eth_subscribe"
	EthUnsubscribe = "eth_unsubscribe"

	//历史信息类
	EthProtocolVersion = "eth_protocolVersion"
	EthChainId = "eth_chainId"
	EthGetTransactionCount = "eth_getTransactionCount"
	EthGetBlockTransactionCountByNumber = "eth_getBlockTransactionCountByNumber"
	EthGetBlockTransactionCountByHash = "eth_getBlockTransactionCountByHash"
	EthGetCode = "eth_getCode"
	EthGetBlockByNumber = "eth_getBlockByNumber"
	EthGetBlockByHash = "eth_getBlockByHash"
	EthGetTransactionByHash = "eth_getTransactionByHash"
	EthGetTransactionByBlockHashAndIndex = "eth_getTransactionByBlockHashAndIndex"
	EthGetTransactionReceipt = "eth_getTransactionReceipt"
	EthGetLogs = "eth_getLogs"
	EthGetTransactionLogs = "eth_getTransactionLogs"
	EthGetTransactionbyBlockNumberAndIndex = "eth_getTransactionbyBlockNumberAndIndex"

	//待定类
	EthSign = "eth_sign"
	EthGetFilterChanges = "eth_getFilterChanges"
	EthGetFilterLogs = "eth_getFilterLogs"
)

const (
	jsonrpc						= "2.0"
	id							= 1
	contentType					= "application/json"
	post						= "POST"
	TestNetUrl					= "https://exchaintestrpc.okex.org"
	MainNetUrl					= "https://exchainrpc.okex.org"
)

type ReqBody struct {
	Jsonrpc  	string  		`json:"jsonrpc"`
	Method 		string			`json:"method"`
	Params		interface{}		`json:"params"`
	Id			int				`json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RPCResp struct {
	Error  		*RPCError       `json:"error"`
	Jsonrpc  	string  		`json:"jsonrpc"`
	Id     		int             `json:"id"`
	Result 		json.RawMessage `json:"result,omitempty"`
}


func NewReqBody(jsonrpc string, method string, params interface{}, id int) *ReqBody {
	return &ReqBody{
		Jsonrpc: jsonrpc,
		Method: method,
		Params: params,
		Id: 	id,
	}
}

/**
url:
测试网 https://exchaintestrpc.okex.org
主网 https://exchainrpc.okex.org
保留 http.Response 是为了做更多的字段验证
 */
func DoPost(url string, postBody []byte) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(post, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	startTime := time.Now()
	resp, respErr := client.Do(req)
	elapsed := time.Since(startTime)

	if respErr != nil {
		log.Println(postBody, strconv.FormatInt(elapsed.Milliseconds(),10) + "ms", "Fail", respErr)
		return nil, respErr
	}

	log.Println(postBody, strconv.FormatInt(elapsed.Milliseconds(),10) + "ms", "Success")
	defer resp.Body.Close()
	return resp, nil
}

/**
将response的body解码成RPCResp结构
只验证body中的返回信息
 */
func GetRespBody(resp *http.Response)(rpcResp *RPCResp, err error){
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)

	if err != nil {
		log.Println("fail", err)
		return nil, err
	}
	if rpcResp.Error != nil {
		log.Println("fail", rpcResp.Error)
		return nil, err
	}
	log.Println("success", string(rpcResp.Result))
	return rpcResp, nil
}

func PanicErr(err error){
	if err != nil {
		panic(err)
	}
}


/**
获取当前区块高度
没有参数
例子：
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' -H "Content-Type: application/json" https://exchaintestrpc.okex.org
{"jsonrpc":"2.0","id":1,"result":"0x92e29c"}

 */
func EthBlockNumberApi(url string) (*http.Response, error) {
	log.Println(fmt.Printf("*********start to run api function EthBlockNumberApi *******"))
	method := EthBlockNumber
	request := NewReqBody(jsonrpc,method,nil, id)
	//处理成为[]byte类型的req
	req, err := json.Marshal(*request)
	PanicErr(err)
	return DoPost(url, req)
}

/**
获取余额
报文中有两个参数：
1. Account Address
2. Block Number
例子：
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xAeFA44f2E8cb4871A0cA862a4E7C5f2761111886", "0x92e23d"],"id":1}' -H "Content-Type: application/json" https://exchaintestrpc.okex.org
{"jsonrpc":"2.0","id":1,"result":"0x1708e7a6bc8d6c00"}
 */
func EthGetBalanceApi(url string, params interface{}) (*http.Response, error) {
	log.Println(fmt.Printf("*********start to run api function EthGetBalanceApi *******"))
	method := EthGetBalance
	request := NewReqBody(jsonrpc, method, params, id)
	req, err := json.Marshal(*request)
	PanicErr(err)
	return DoPost(url, req)
}


/**
发送交易 eth_sendRawTransaction
params的构造，signTx->lrp->hex
例如：
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf86c8210ba843b9aca0082520894aefa44f2e8cb4871a0ca862a4e7c5f27611118868609184e72a0008081a6a070c0ea551370c8fb4bea6ac890f93f525517471045776a5c1366a862b9f84d2ea00c493ec6e3699ed56ce090ab0a68ffb943fe0cc735914572b65bb6eb699f3c5e"],"id":1}'
-H "Content-Type: application/json" https://exchaintestrpc.okex.org
{"jsonrpc":"2.0","id":1,"result":"0xac2fe3ac6bd424b8b6fb80d74b88ecfdc9347ceeee14d68f9eebf8ebe6f037a6"}
*/
func EthSendRawTransactionApi(url string, params interface{}) (*http.Response, error) {
	log.Println(fmt.Printf("*********start to run api function EthSendRawTransactionApi *******"))
	method := EthSendRawTransaction
	request := NewReqBody(jsonrpc, method, params, id)
	req, err := json.Marshal(*request)
	PanicErr(err)
	return DoPost(url, req)
}



