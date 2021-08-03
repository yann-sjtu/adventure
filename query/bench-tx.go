package query

import (
	bytes2 "bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	adcomm "github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

var (
	concurrencyTx   int
	sleepTimeTx     int
	hostTx          string
	chainIdTx       int
	privkPath       string
	contractAddress string
	abiPath         string
)

var (
	sampleContractABI abi.ABI
)

func BenchTxCmd() *cobra.Command {
	// add flags
	//adventure bench-make-tx --concurrency 1 --sleepTime 4 --chainId 65 --host https://exchaintestrpc.okex.org --privkeyPath /Users/shaoyunzhan/Documents/wp/go_wp/adventure/template/privkey/priv_test_block --abiPath /Users/shaoyunzhan/Documents/wp/go_wp/adventure/template/contract/TestBlock.abi  --contractAddress 0x760f3336C4fF1b25C9dFAB105B693E31B4475a15 -c /Users/shaoyunzhan/Documents/wp/go_wp/adventure/config.toml
	cmd := &cobra.Command{
		Use:   "bench-make-tx",
		Short: "",
		Long:  "",
		Run:   benchTx,
	}
	flags := cmd.Flags()
	flags.IntVar(&concurrencyTx, "concurrency", 10, "set the number of tx number per second")
	flags.IntVar(&sleepTimeTx, "sleepTime", 1, "")
	flags.StringVar(&hostTx, "host", "https://exchaintestrpc.okex.org", "")
	flags.IntVar(&chainIdTx, "chainId", 65, "")
	flags.StringVar(&privkPath, "privkeyPath", "", "")
	flags.StringVar(&contractAddress, "contractAddress", "", "")
	flags.StringVar(&abiPath, "abiPath", "", "")
	return cmd
}

func benchTx(cmd *cobra.Command, args []string) {
	abiByte, err := ioutil.ReadFile(abiPath)
	if err != nil {
		log.Fatal(err)
	}

	sampleContractABI, err = abi.JSON(bytes2.NewReader(abiByte))
	if err != nil {
		log.Fatal(err)
	}

	privkeys := adcomm.GetPrivKeyFromPrivKeyFile(privkPath)
	for {
		for i := 0; i < concurrencyTx; i++ {
			go func(privkey string) {
				privateKey, _ := crypto.HexToECDSA(privkey)
				client, err := ethclient.Dial(host)
				if err != nil {
					log.Fatalf("failed to initialize client: %+v", err)
				}
				nonce, _ := client.PendingNonceAt(context.Background(), getAddress(privateKey))

				for {
					unsignedTx := buildUnsignedTx(nonce, common.HexToAddress(contractAddress))
					// 2. sign unsignedTx -> rawTx
					signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(int64(chainIdTx))), privateKey)
					if err != nil {
						log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
					}

					// 3. send rawTx
					err = client.SendTransaction(context.Background(), signedTx)
					if err != nil {
						log.Fatal(err)
					}

					nonce++
					time.Sleep(time.Second * time.Duration(sleepTimeTx))
				}
			}(privkeys[i])
		}
	}

}

func getAddress(privateKey *ecdsa.PrivateKey) common.Address  {
	pubkeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if ok != true {
		panic(fmt.Errorf("convert into pubkey failed"))
	}
	fromAddress := crypto.PubkeyToAddress(*pubkeyECDSA)
	return fromAddress
}

func buildUnsignedTx(nonce uint64, contractAddr common.Address) *types.Transaction {
	gasPrice := big.NewInt(1000000000)
	value := big.NewInt(0)
	gasLimit := uint64(3000000)

	data, err := sampleContractABI.Pack("operate")
	if err != nil {
		log.Fatal(err)
	}
	trans := types.NewTransaction(nonce, contractAddr, value, gasLimit, gasPrice, data)
	log.Println("tx ,", trans.Hash())
	return trans
}
