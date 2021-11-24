package bench

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/spf13/cobra"
)


var (
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
	flags.IntVar(&ethPort, "eth-port", 0,"if not zero, query on eth port 26659")
	return cmd
}

func startOperate(cmd *cobra.Command, args []string) {
	var txdata string
	if direct {
		txdata = generateTxDataInDirect()
	} else {
		txdata = generateTxData()
	}

	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	for i := 0; i < concurrency; i++ {
		go func(index int, privkey string) {
			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			operate(privkey, rpcHost, txdata)

		}(i, privkeys[i])
		time.Sleep(time.Millisecond * 10)
	}

	select {}
}

func operate(privkey string, host string, txdata string) {
	nonce := queryNonce(host, privkey)
	fmt.Println(getCosmosAddress(privkey).String(), nonce)

	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)
	for {
		res, err := cli.Evm().SendTxEthereum(privkey, contract, "", txdata,2000000, nonce)
		if err != nil {
			log.Printf("[cosmos] send tx err: %s\n", err)
			if strings.Contains(err.Error(), "tx already exists in cache") {
				nonce++
			}
			time.Sleep(time.Second * time.Duration(sleepTime))
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

//0x4b13e557000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002
func generateTxData() string {
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
	return ethcmm.Bytes2Hex(txdata)
}

func generateTxDataInDirect() string {
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
	return ethcmm.Bytes2Hex(txdata)
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
