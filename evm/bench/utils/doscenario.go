package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/client"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
	"log"

)



/**
遍历账户，获取账户余额， 转账，然后再次获取用户余额，验证下余额变化，余额减少
*/

func GetBalTxBal(p BasepParam, e func(ethcmm.Address) []TxParam) {
	clients := client.GenerateClients(p.ips)    // generate CosmosClient or EthClient
	accounts := generateAccounts(p.privateKeys) // generate accounts

	for i := 0; i < p.concurrency; i++ {
		go func(gIndex int) {
			for j := 0; ; j++ {
				aIndex := (gIndex + j*p.concurrency) % len(accounts) // make sure accounts will be picked in order by round-robin
				acc := accounts[aIndex]
				cli := clients[aIndex%len(clients)]

				//获取余额
				resp, _ := GetAccBalance(gIndex, acc, TestNetUrl)
				bal1 := string(resp.Result)
				//执行tx
				execute(gIndex, cli, acc, e)
				//再次获取余额
				resp2, _ := GetAccBalance(gIndex, acc, TestNetUrl)
				bal2 := string(resp2.Result)

				//验证bal2 小于 bal1
				AssertCompare(bal1, bal2, "bal1 should be greater than bal2")
				time.Sleep(time.Millisecond * time.Duration(p.sleep))
			}
		}(i)
	}

	select {}
}

/**
16进制数据比较，使用Big.Int
 */

func AssertCompare(val1 string, val2 string, errInfo string)  {

	a, err1 := hex.DecodeString(val1)
	b, err2 := hex.DecodeString(val2)
	if err1 != nil || err2 !=nil {
		log.Println("err1 is : %s; err2 is : %s", err1, err2)
	}
	intA := new(big.Int).SetBytes(a)
	intB := new(big.Int).SetBytes(b)
	n := intA.Cmp(intB)
	if n >0 {
		log.Println("success to assert")
		return
	}
	log.Println("fail to assert, error happen: %s, val1 is: %s; val2 is : %s", errInfo, intA, intB)
	return
}

func GetAccBalance(gIndex int, acc *EthAccount, url string)(rpcResp *RPCResp, err error){
	acc.Lock()
	defer acc.Unlock()
	params := make([]string, 0, 5)
	//构造request
	address := common.GetEthAddressFromPK(acc.GetPrivateKey())
	res, _:= GetBlockNumber(gIndex,acc,url)

	params = append(params, address.String())
	params = append(params, string(res.Result))

	resp, _ := EthGetBalanceApi(url, params)
	rpcResp, e := GetRespBody(resp)
	if  e != nil {
		log.Println(fmt.Errorf("[g%d] failed to get block number, error: %s", gIndex, err))
		return nil, e
	}
	return rpcResp,nil
}

func GetBlockNumber(gIndex int, acc *EthAccount, url string)(rpcResp *RPCResp, err error) {
	acc.Lock()
	defer acc.Unlock()

	resp, _ := EthBlockNumberApi(url)
	body, e := GetRespBody(resp)
	if  e != nil {
		log.Println(fmt.Errorf("[g%d] failed to get block number, error: %s", gIndex, err))
		return nil, e
	}
	return body, nil
}



func RespToMap(resp *http.Response) map[string]interface{}{
	tmp := make(map[string]interface{})
	r, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println(fmt.Errorf("Response to map is err: %s", e))
	}
	json.Unmarshal([]byte(string(r)),&tmp)
	return tmp
}