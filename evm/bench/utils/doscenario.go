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
	aIndex := 0
	acc := accounts[aIndex]
	cli := clients[aIndex%len(clients)]
	log.Println(fmt.Errorf("[concurrency%d][aIndex%d] start to execBalTxBal test",0,0))
	execBalTxBal(aIndex, cli, acc, e)
}

func GetBalTxBal_bak(p BasepParam, e func(ethcmm.Address) []TxParam) {
	clients := client.GenerateClients(p.ips)    // generate CosmosClient or EthClient
	accounts := generateAccounts(p.privateKeys) // generate accounts

	for i := 0; i < p.concurrency; i++ {
		go func(gIndex int) {
			for j := 0; ; j++ {
				aIndex := (gIndex + j*p.concurrency) % len(accounts) // make sure accounts will be picked in order by round-robin
				acc := accounts[aIndex]
				cli := clients[aIndex%len(clients)]
				log.Println(fmt.Errorf("[concurrency%d][gIndex%d] start to execBalTxBal test",i,gIndex))
				execBalTxBal(gIndex, cli, acc, e)
				time.Sleep(time.Millisecond * time.Duration(p.sleep))
			}
		}(i)
	}

	select {}
}

func execBalTxBal(gIndex int, cli client.Client, acc *EthAccount, e func(ethcmm.Address) []TxParam){
	//获取余额
	resp, _ := GetAccBalance(gIndex, acc, TestNetUrl)
	bal1 := string(resp.Result)
	//执行tx
	execute(gIndex, cli, acc, e)
	//再次获取余额
	resp2, _ := GetAccBalance(gIndex, acc, TestNetUrl)
	bal2 := string(resp2.Result)

	//验证bal2 小于 bal1
	bRet := AssertCompare(bal1, bal2, "bal1 should be greater than bal2")
	sRet := "FAIL"
	if bRet == true{
		sRet = "SUCCESS"
	}
	log.Println(fmt.Errorf(" ****[%d] finish execBalTxBal, and %s\n", gIndex,sRet))
}

/**
16进制数据比较，使用Big.Int
 */

func AssertCompare(val1 string, val2 string, errInfo string) (bRet bool) {
	log.Println(fmt.Errorf("*********start to do comparison AssertCompare *******"))
	bRet = false
	a, err1 := hex.DecodeString(val1)
	b, err2 := hex.DecodeString(val2)
	if err1 != nil || err2 !=nil {
		log.Println(fmt.Errorf("err1 is : %s; err2 is : %s", err1, err2))
	}
	intA := new(big.Int).SetBytes(a)
	intB := new(big.Int).SetBytes(b)
	n := intA.Cmp(intB)
	if n >0 {
		log.Println(fmt.Printf("success to assert"))
		bRet = true
		return bRet
	}
	log.Println(fmt.Errorf("fail to assert, error happen: %s, val1 is: %s; val2 is : %s", errInfo, intA, intB))
	return bRet
}

func GetAccBalance(gIndex int, acc *EthAccount, url string)(rpcResp *RPCResp, err error){

	//prefix := "0x"
	params := make([]string, 0, 5)
	//构造request
	address := common.GetEthAddressFromPK(acc.GetPrivateKey())
	res, _:= GetBlockNumber(gIndex,acc,url)
	params = append(params, address.String())
	params = append(params, string(res.Result))
	log.Println(fmt.Errorf("params[0] is : %s;params[1] is : %s; ", params[0], params[1]))
	//调用函数，获得返回
	rpcResp, err = EthGetBalanceApi(url, params)

	if  err != nil {
		log.Println(fmt.Errorf("[g%d] failed to get account balance, error: %s", gIndex, err))
		return nil, err
	}
	return rpcResp,nil
}

func GetBlockNumber(gIndex int, acc *EthAccount, url string)(rpcResp *RPCResp, err error) {

	rpcResp, err = EthBlockNumberApi(url)
	if  err != nil {
		log.Println(fmt.Errorf("[g%d] failed to get block number, error: %s", gIndex, err))
		return nil, err
	}
	return rpcResp, nil
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