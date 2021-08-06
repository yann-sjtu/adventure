package query

import (
	bytes2 "bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	adcomm "github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	adtypes "github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	concurrencyTx   int
	sleepTimeTx     int
	privkPath       string
	contractAddress string
	abiPath         string

	rest_host       string
	rest_chainId    int
	rpc_hosts       []string
	rpc_chainId     string
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

	flags.StringVar(&privkPath, "privkeyPath", "", "")
	flags.StringVar(&contractAddress, "contractAddress", "", "")
	flags.StringVar(&abiPath, "abiPath", "", "")

	flags.StringVar(&rest_host, "rest-host", "", "")
	flags.IntVar(&rest_chainId, "rest-chainid", 65, "")
	flags.StringSliceVar(&rpc_hosts, "rpc-hosts", []string{}, "")
	flags.StringVar(&rpc_chainId, "rpc-chainid", "", "")
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
	var wg sync.WaitGroup
	wg.Add(concurrencyTx)
	for i := 0; i < concurrencyTx; i++ {
		go func(index int, privkey string) {
			defer wg.Done()

			if rest_host != "" {
				sendTxToRestNodes(privkey, rest_host)
			} else if rpc_hosts != nil {
				rpcHost := rpc_hosts[index%len(rpc_hosts)]
				sendTxToRpcNodes(privkey, rpcHost)
			} else {
				panic(fmt.Errorf("no host"))
			}

		}(i, privkeys[i])
	}
	wg.Wait()
}

func sendTxToRestNodes(privkey string, host string) {
	privateKey, _ := crypto.HexToECDSA(privkey)
	client, err := ethclient.Dial(host)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), getEthAddress(privateKey))
	if err != nil {
		log.Fatalf("failed to fetch noce: %+v", err)
	}
	//fmt.Println(getEthAddress(privateKey), nonce)

	for {
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Println("failed to fetch gas price:", err)
			gasPrice = big.NewInt(1000000000)
		}

		// 2. sign unsignedTx -> rawTx
		signedTx, err := types.SignTx(
			buildUnsignedTx(nonce, gasPrice, common.HexToAddress(contractAddress)),
			types.NewEIP155Signer(big.NewInt(int64(rest_chainId))),
			privateKey,
		)
		if err != nil {
			log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
		}

		// 3. send rawTx
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Printf("err: %s\n", err)
			continue
		}
		log.Println("txhash:", signedTx.Hash().String(), "gasPrice", gasPrice.String())
		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}

func sendTxToRpcNodes(privkey string, host string) {
	cfg, _ := adtypes.NewClientConfig(host, rpc_chainId, adtypes.BroadcastSync, "", 1000000, 1.5, "0.0000000001"+adcomm.NativeToken)
	cli := gosdk.NewClient(cfg)

	addr := getCosmosAddress(privkey)
	accInfo, err := cli.Auth().QueryAccount(addr.String())
	if err != nil {
		panic(err)
	}
	//fmt.Println(addr.String(), accInfo.GetSequence())

	gasPrice := big.NewInt(int64((rand.Intn(20)+1) * 100000000))
	payload := buildCosmosTxData()
	index := 0
	for {
		res, err := cli.Evm().SendTxEthereum(privkey, contractAddress, "", common.Bytes2Hex(payload), 1000000, accInfo.GetSequence()+uint64(index), gasPrice)
		if err != nil {
			log.Printf("err: %s\n", err)
			if strings.Contains(err.Error(), "mempool") || strings.Contains(err.Error(), "EOF") {
				time.Sleep(time.Second * 10)
				continue
			}
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		index++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}

func getEthAddress(privateKey *ecdsa.PrivateKey) common.Address {
	pubkeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if ok != true {
		panic(fmt.Errorf("convert into pubkey failed"))
	}
	fromAddress := crypto.PubkeyToAddress(*pubkeyECDSA)
	return fromAddress
}

func getCosmosAddress(privkey string) sdk.Address {
	privateKey, _ := crypto.HexToECDSA(privkey)
	cosmosAddr, err := utils.ToCosmosAddress(getEthAddress(privateKey).String())
	if err != nil {
		panic(err)
	}
	return cosmosAddr
}

func buildCosmosTxData() []byte {
	data, err := sampleContractABI.Pack("operate")
	if err != nil {
		panic(err)
	}
	return data
}

func buildUnsignedTx(nonce uint64, gasPrice *big.Int, contractAddr common.Address) *types.Transaction {
	value := big.NewInt(0)
	gasLimit := uint64(3000000)

	data, err := sampleContractABI.Pack("operate")
	if err != nil {
		log.Fatal(err)
	}
	trans := types.NewTransaction(nonce, contractAddr, value, gasLimit, gasPrice, data)
	return trans
}
