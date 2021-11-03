package bench

import (
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/spf13/cobra"
)


var (
	RouterTestContract string

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
	flags.IntVarP(&sleepTimeTx, "sleepTime", "s",1, "")

	flags.StringVarP(&privkPath, "privkey-path", "p", "","")
	flags.StringSliceVarP(&rpc_hosts, "rpc-hosts","u", []string{}, "")
	flags.StringVarP(&chainID, "chain-id", "i","", "")

	flags.Int64Var(&id, "id", 0, "")
	flags.Int64SliceVar(&opts, "opts", []int64{}, "")
	flags.Int64Var(&times, "times", 0, "")

	flags.StringVar(&RouterTestContract, "router-contract", "0xdA1BD71c96104F7794F263AcD04334870Cb428B7","")
	return cmd
}

func startOperate(cmd *cobra.Command, args []string) {
	txdata := generateTxData()

	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	for i := 0; i < concurrency; i++ {
		go func(index int, privkey string) {
			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			operate(privkey, rpcHost, txdata)

		}(i, privkeys[i])
	}

	select {}
}

func operate(privkey string, host string, txdata string) {
	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	addr := getCosmosAddress(privkey)
	accInfo, err := cli.Auth().QueryAccount(addr.String())
	if err != nil {
		panic(err)
	}
	nonce := accInfo.GetSequence()

	for {
		res, err := cli.Evm().SendTxEthereum(privkey, RouterTestContract, "", txdata,2000000, nonce)
		if err != nil {
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
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
