package bench

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/spf13/cobra"
)

var (
	duplicate bool
	ethPort int

	contract string
	direct   bool

	id int64
	opts []int64
	times int64
)

func OperateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operate",
		Short: "",
		Long:  "",
		Run:   startOperate,
	}
	flags := cmd.Flags()
	flags.IntVarP(&concurrency, "concurrency","g", 10, "set the number of tx number per second")
	flags.IntVarP(&sleepTime, "sleepTime", "s",1, "")

	flags.StringVarP(&privkPath, "privkey-path", "p", "","")
	flags.StringSliceVarP(&rpc_hosts, "rpc-hosts","u", []string{}, "")
	flags.StringVarP(&chainID, "chain-id", "i","", "")

	flags.Int64Var(&id, "id", 0, "")
	flags.Int64SliceVar(&opts, "opts", []int64{}, "")
	flags.Int64Var(&times, "times", 0, "")

	flags.StringVar(&contract, "contract", "","")
	flags.BoolVar(&direct, "direct", false,"")
	flags.BoolVar(&duplicate, "duplicate", false,"")
	flags.IntVar(&ethPort, "eth-port", 0,"if not zero, query on eth port 26659")
	return cmd
}

func startOperate(cmd *cobra.Command, args []string) {
	var txdata []byte
	if direct {
		txdata = generateTxDataInDirect()
	} else {
		txdata = generateTxData()
	}

	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	for i := 0; i < concurrency; i++ {
		go func(index int, privkey string) {
			privateKey, err := crypto.HexToECDSA(privkey)
			if err != nil {
				panic(err)
			}

			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			operate(privateKey, rpcHost, txdata)

		}(i, privkeys[i])
		time.Sleep(time.Millisecond * 10)
	}

	select {}
}

func operate(privateKey *ecdsa.PrivateKey, host string, txdata []byte) {
	to := ethcmm.HexToAddress(contract)
	nonce := queryNonce(host, privateKey)
	fmt.Println(getCosmosAddress(privateKey).String(), nonce)

	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)
	for {
		res, err := cli.Evm().SendTxEthereum(privateKey, nonce, to, nil,2000000, evmtypes.DefaultGasPrice, txdata)
		if err != nil {
			log.Printf("[cosmos] send tx err: %s\n", err)
			if strings.Contains(err.Error(), "tx already exists in cache") {
				nonce++
			}
			time.Sleep(time.Second * time.Duration(sleepTime))
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
			if duplicate {
				sendDuplicateTx(cli, privateKey, nonce, to, txdata)
			}
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

//0x4b13e557000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002
func generateTxData() []byte {
	routerABI, err := abi.JSON(strings.NewReader(RouterABI))
	if err != nil {
		log.Fatal(err)
	}
	bigOpts := make([]*big.Int, len(opts))
	for i := 0; i < len(opts); i++ {
		bigOpts[i] = big.NewInt(opts[i])
	}
	txdata, err := routerABI.Pack("operate", big.NewInt(id), bigOpts, big.NewInt(times))
	if err != nil {
		log.Fatal(err)
	}
	return txdata
}

func generateTxDataInDirect() []byte {
	readABI, err := abi.JSON(strings.NewReader(ReadABI))
	if err != nil {
		log.Fatal(err)
	}
	bigOpts := make([]*big.Int, len(opts))
	for i := 0; i < len(opts); i++ {
		bigOpts[i] = big.NewInt(opts[i])
	}
	txdata, err := readABI.Pack("operate", bigOpts, big.NewInt(times))
	if err != nil {
		log.Fatal(err)
	}
	return txdata
}

func queryNonce(host string, privateKey *ecdsa.PrivateKey) (nonce uint64) {
 	if ethPort == 0 {
		cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
		cli := gosdk.NewClient(cfg)

		addr := getCosmosAddress(privateKey)
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

func sendDuplicateTx(cli gosdk.Client, privateKey *ecdsa.PrivateKey, nonce uint64, to ethcmm.Address, txdata []byte) {
	rand.Seed(time.Now().Unix())
	if rand.Intn(100) < 25 { // 25% chance to send duplicate txs
		for i := 1; i <= rand.Intn(3) + 1; i++ {
			gasPrice := big.NewInt(1).Mul(evmtypes.DefaultGasPrice, big.NewInt(int64(i)))

			res, err := cli.Evm().SendTxEthereum(privateKey, nonce, to, nil,2000000, gasPrice, txdata)
			if err != nil {
				log.Printf("[duplicate] error: %s\n", err)
			} else {
				log.Printf("[duplicate] txhash: %s\n", res.TxHash)
			}
		}
	}
}
