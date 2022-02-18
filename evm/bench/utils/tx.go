package utils

import (
	"fmt"
	"log"
	"math/big"
	"strings"
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

func NewTxParam(to ethcmm.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) TxParam {
	return TxParam{
		to,
		amount,
		gasLimit,
		gasPrice,
		data,
	}
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
