package evmtx

import (
	"encoding/json"
	"log"
	"math/big"
	"sync"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm"
	"github.com/okex/adventure/x/strategy/evm/template/DYF"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	sleepTime   int
	host        string
	concurrency int
	privkPath   string
)

func BenchTxCmd() *cobra.Command {
	// add flags
	cmd := &cobra.Command{
		Use:   "bench-tx",
		Short: "",
		Long:  "",
		Run:   benchTx,
		PreRun: func(cmd *cobra.Command, args []string) {
			evm.InitTemplate()

			depositPayloadStr = hexutil.Encode(UniswapV2.BuildWethDepositPayload())
			approvePayloadStr = hexutil.Encode(UniswapV2.BuildWethApprovePayload(
				routerAddr, 5,
			))
			swapPayloadStr = hexutil.Encode(UniswapV2.BuildSwapExactTokensForTokensPayload(
				big.NewInt(1000), big.NewInt(0),
				[]string{wethAddr, usdtAddr}, "0x2B5Cf24AeBcE90f0B8f80Bc42603157b27cFbf47",
				time.Now().Add(time.Hour*8640).Unix(),
			))
			dyfPayloadStr = hexutil.Encode(DYF.BuildExcutePayload())
		},
	}
	flags := cmd.Flags()
	flags.IntVarP(&concurrency, "concurrency", "g", 1, "set the number of query concurrent number per second")
	flags.IntVarP(&sleepTime, "sleeptime", "t", 1, "")
	flags.StringVarP(&host, "host", "o", "https://exchaintestrpc.okex.org", "")
	flags.StringVarP(&privkPath, "privkey", "p", "", "")
	return cmd
}

func benchTx(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath, concurrency)
	for i := 0; i < concurrency; i++ {
		go func(privkey string) {
			privateKey, err := crypto.HexToECDSA(privkey)
			if err != nil {
				log.Fatalf("failed to switch unencrypted private key -> secp256k1 private key: %+v", err)
			}
			info, err := utils.CreateAccountWithPrivateKey(privkey, "acc", common.PassWord)
			if err != nil {
				panic(err)
			}
			ethAddrHex, err := utils.ToHexAddress(info.GetAddress().String())
			if err != nil {
				panic(err)
			}
			ethAddr := ethAddrHex.String()

			for r := 0; ;r++{
				var wg sync.WaitGroup
				wg.Add(3)

				// mint、approve、transfer
				// 1. estimate gas
				var gas hexutil.Uint64
				param := generateTxParams(ethAddr, r%4)
				go func(p []map[string]string) {
					defer wg.Done()
					rpcRes, err := CallWithError("eth_estimateGas", p)
					if err != nil {
						log.Println("eth_estimateGas", err)
						return
					}
					if err = json.Unmarshal(rpcRes.Result, &gas); err != nil {
						panic(err)
					}
					//fmt.Println(uint64(gas))
				}(param)

				// 2. fetch gas price
				var gasPrice hexutil.Big
				go func() {
					defer wg.Done()
					rpcRes, err := CallWithError("eth_gasPrice", nil)
					if err != nil {
						log.Println("eth_gasPrice", err)
						return
					}
					if err = json.Unmarshal(rpcRes.Result, &gasPrice); err != nil {
						panic(err)
					}
					//fmt.Println(gasPrice.String())
				}()

				// 3. eth_getTransactionCount
				var nonce hexutil.Uint64
				go func() {
					defer wg.Done()
					rpcRes, err := CallWithError("eth_getTransactionCount", []interface{}{ethAddr, "pending"})
					if err != nil {
						log.Println("eth_getTransactionCount", err)
						return
					}
					if err = json.Unmarshal(rpcRes.Result, &nonce); err != nil {
						panic(err)
					}
					//fmt.Println(uint64(nonce))
				}()
				wg.Wait()

				// 4. eth_sendRawTransaction
				data := signTx(privateKey, nonce, param[0]["to"], param[0]["value"], gas, gasPrice, param[0]["data"])
				rpcRes, err := CallWithError("eth_sendRawTransaction", []interface{}{data})
				if err != nil {
					log.Println("eth_sendRawTransaction", err)
					continue
				}
				var txhash ethcmn.Hash
				if err = json.Unmarshal(rpcRes.Result, &txhash); err != nil {
					panic(err)
				}
				//fmt.Println(txhash.String())

				// 5. getTransactionReceipt
				go func(hash string) {
					for {
						rpcRes, err := CallWithError("eth_getTransactionReceipt", []interface{}{hash})
						if err != nil {
							log.Println("eth_getTransactionReceipt", err)
							return
						}
						if string(rpcRes.Result) == "null" {
							time.Sleep(time.Second*2)
							continue
						}

						var receipt map[string]interface{}
						if err = json.Unmarshal(rpcRes.Result, &receipt); err != nil {
							panic(err)
						}
						if receipt["status"].(string) == hexutil.Uint(1).String() {
							//log.Println("done")
						}
						break
					}
				}(txhash.String())

				time.Sleep(time.Second*time.Duration(sleepTime))
			}

		}(privkeys[i])
	}

	select {}
}

// eth_getcode?