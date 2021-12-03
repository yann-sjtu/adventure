package evm_transfer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
)

var (
	ethPort int

	concurrencyTx   int
	sleepTimeTx     int
	privkPath       string
	rpc_hosts    []string

	chainID       string

	fixed  bool
)

func TransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evm-transfer",
		Short: "",
		Long:  "",
		Run:   transferTx,
	}
	flags := cmd.Flags()
	flags.IntVar(&concurrencyTx, "concurrency", 10, "set the number of tx number per second")
	flags.IntVar(&sleepTimeTx, "sleepTime", 1, "")

	flags.StringVar(&privkPath, "privkey-path", "", "")
	flags.StringSliceVar(&rpc_hosts, "rpc-hosts", []string{}, "")
	flags.StringVar(&chainID, "chain-id", "", "")

	flags.BoolVar(&fixed, "fixed", false, "")
	flags.IntVar(&ethPort, "eth-port", 0,"if not zero, query on eth port 26659")
	return cmd
}

func transferTx(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	to := ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64)))
	for i := 0; i < concurrencyTx; i++ {
		go func(index int, privkey string) {
			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			transfer(privkey, rpcHost, to.String())

		}(i, privkeys[i])
	}

	select {}
}

func transfer(privkey string, host string, to string) {
	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 30000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	nonce := queryNonce(host, privkey)
	fmt.Println(getCosmosAddress(privkey).String(), nonce)

	for {
		if !fixed {
			to = ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64))).String()
		}

		res, err := cli.Evm().SendTxEthereum2(privkey, to, "0.000000001", "",21000, nonce)
		if err != nil {
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}

func getEthAddress(privateKey *ecdsa.PrivateKey) ethcmm.Address {
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

func queryNonce(host string, privkey string) (nonce uint64) {
	if ethPort == 0 {
		cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
		cli := gosdk.NewClient(cfg)

		addr := getCosmosAddress(privkey)
		for i := 0; ; i++ {
			accInfo, err := cli.Auth().QueryAccount(addr.String())
			if err != nil {
				log.Printf("[cosmos] query %s error: %s\n", addr.String(), err)
				time.Sleep(time.Second)
				continue
			}
			nonce = accInfo.GetSequence()
			return
		}
	} else {
		str := strings.Split(host, ":")
		ethhost := str[0] + ":" + str[1] + ":" + strconv.Itoa(ethPort)
		client, err := ethclient.Dial(ethhost)
		if err != nil {
			panic(err)
		}

		privateKey, err := crypto.HexToECDSA(privkey)
		if err != nil {
			panic(err)
		}

		addr := getEthAddress(privateKey)
		for i := 0; ; i++ {
			nonce, err = client.PendingNonceAt(context.Background(), addr)
			if err != nil {
				log.Printf("[eth] query %s error: %s\n", addr.String(), err)
				time.Sleep(time.Second)
				continue
			}
			return
		}
	}
}