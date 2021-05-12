package evmtx

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm"
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

			for {
				// mint、approve、transfer
				// 1. estimate gas
				param := generateTxParams(ethAddr)
				rpcRes, err := CallWithError("eth_estimateGas", param)
				if err != nil {
					log.Println(err)
					continue
				}
				var gas hexutil.Uint64
				if err = json.Unmarshal(rpcRes.Result, &gas); err != nil {
					panic(err)
				}
				fmt.Println(uint64(gas))

				// 2. fetch gas price
				rpcRes, err = CallWithError("eth_gasPrice", nil)
				if err != nil {
					log.Println(err)
					continue
				}
				var gasPrice hexutil.Big
				if err = json.Unmarshal(rpcRes.Result, &gasPrice); err != nil {
					panic(err)
				}
				fmt.Println(gasPrice.String())

				// 3. eth_getTransactionCount
				rpcRes, err = CallWithError("eth_getTransactionCount", []interface{}{ethAddr, "pending"})
				var nonce hexutil.Uint64
				if err = json.Unmarshal(rpcRes.Result, &nonce); err != nil {
					panic(err)
				}
				fmt.Println(uint64(nonce))

				// sendRawTransaction
				data := signTx(privateKey, uint64(nonce), uint64(gas), (*big.Int)(&gasPrice))
				rpcRes, err = CallWithError("eth_sendRawTransaction", []interface{}{data})
				var txhash ethcmn.Hash
				if err = json.Unmarshal(rpcRes.Result, &txhash); err != nil {
					panic(err)
				}
				fmt.Println(txhash.String())

				// getTransactionReceipt
			}

		}(privkeys[i])
	}

	select {}
}

// eth_getcode?