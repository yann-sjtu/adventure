package utils

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/client"
)

type TxParam struct {
	to       ethcmm.Address
	amount   *big.Int
	gasLimit uint64
	gasPrice *big.Int
	data     []byte
}

var (
	lstTxHash = make([]string, 0)
	duration	int64
	ratio		float32
	tps			int64

)
/**
作用：用来计算并发携程一次发送完毕后的的成功率
 */
func GetTxTpsAndSuccessRatio(lstTxHash []string, cocurrent int64)(ratio float32, tps int64){
	num := len(lstTxHash)
	ratio = float32(num)/float32(cocurrent)
	tps = cocurrent/duration/1000
	return
}

func getTxHashList(gIndex int, cli client.Client, acc *EthAccount, e func(ethcmm.Address) []TxParam) ([]string){
	acc.Lock()
	defer acc.Unlock()

	caller := common.GetEthAddressFromPK(acc.GetPrivateKey())
	if err := acc.SetNonce(cli); err != nil {
		log.Println(fmt.Errorf("[g%d] failed to query %s nonce, error: %s", gIndex, caller, err))
		return lstTxHash
	}

	eParams := e(caller)
	for _, eParam := range eParams {
		txhash, err := cli.SendEthereumTx(acc.GetPrivateKey(), acc.GetNonce(), eParam.to, eParam.amount, eParam.gasLimit, eParam.gasPrice, eParam.data)
		if err != nil {
			log.Printf("[g%d] %s send tx err: %s\n", gIndex, caller, err)
			if strings.Contains(err.Error(), "already exists") {
				acc.AddNonce()
			} else if strings.Contains(err.Error(), "mempool is full") {
				time.Sleep(time.Second)
			} else if strings.Contains(err.Error(), "invalid nonce") {
				acc.AddNonce()
			}
		} else {
			log.Printf("[g%d] %s txhash: %s\n", gIndex, caller, txhash)
			lstTxHash = append(lstTxHash, txhash.String())
			acc.AddNonce()
		}
	}
	return lstTxHash
}


func NewTxParam(to ethcmm.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) TxParam {
	return TxParam{
		to,
		amount,
		gasLimit,
		gasPrice,
		data,
	}
}

func RunTxRpc(p BasepParam, e func(ethcmm.Address) []TxParam) {
	clients := client.GenerateClients(p.ips)    // generate CosmosClient or EthClient
	accounts := generateAccounts(p.privateKeys) // generate accounts

	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < p.concurrency; i++ {
		wg.Add(1)
		go func(gIndex int) {
			//j<1是为了获取一次交易
			for j := 0; j<1 ; j++ {
				aIndex := (gIndex + j*p.concurrency) % len(accounts) // make sure accounts will be picked in order by round-robin
				acc := accounts[aIndex]
				cli := clients[aIndex%len(clients)]

				getTxHashList(gIndex, cli, acc, e)
				//time.Sleep(time.Millisecond * time.Duration(p.sleep))
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	duration = time.Since(startTime).Milliseconds()
	elapsed := strconv.FormatInt(time.Since(startTime).Milliseconds(), 10) + "ms"
	ratio, tps = GetTxTpsAndSuccessRatio(lstTxHash,int64(p.concurrency))
	log.Printf(" %d tx send and receive txhash and cost time: %s\n", p.concurrency, elapsed)
	log.Printf(" tx send success ratio is : %d, and tx tps is : %d\n", ratio, tps)
}

func RunTxs(p BasepParam, e func(ethcmm.Address) []TxParam) {
	clients := client.GenerateClients(p.ips)    // generate CosmosClient or EthClient
	accounts := generateAccounts(p.privateKeys) // generate accounts

	for i := 0; i < p.concurrency; i++ {
		go func(gIndex int) {
			for j := 0; ; j++ {
				aIndex := (gIndex + j*p.concurrency) % len(accounts) // make sure accounts will be picked in order by round-robin
				acc := accounts[aIndex]
				cli := clients[aIndex%len(clients)]

				execute(gIndex, cli, acc, e)
				time.Sleep(time.Millisecond * time.Duration(p.sleep))
			}
		}(i)
	}

	select {}
}

func execute(gIndex int, cli client.Client, acc *EthAccount, e func(ethcmm.Address) []TxParam) {
	acc.Lock()
	defer acc.Unlock()

	caller := common.GetEthAddressFromPK(acc.GetPrivateKey())
	if err := acc.SetNonce(cli); err != nil {
		log.Println(fmt.Errorf("[g%d] failed to query %s nonce, error: %s", gIndex, caller, err))
		return
	}

	eParams := e(caller)
	for _, eParam := range eParams {
		txhash, err := cli.SendEthereumTx(acc.GetPrivateKey(), acc.GetNonce(), eParam.to, eParam.amount, eParam.gasLimit, eParam.gasPrice, eParam.data)
		if err != nil {
			log.Printf("[g%d] %s send tx err: %s\n", gIndex, caller, err)
			if strings.Contains(err.Error(), "already exists") {
				acc.AddNonce()
			} else if strings.Contains(err.Error(), "mempool is full") {
				time.Sleep(time.Second)
			} else if strings.Contains(err.Error(), "invalid nonce") {
				acc.AddNonce()
			}
		} else {
			log.Printf("[g%d] %s txhash: %s\n", gIndex, caller, txhash)
			acc.AddNonce()
		}
	}
}

