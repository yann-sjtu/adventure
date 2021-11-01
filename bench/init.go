package bench

import (
	"log"
	"time"

	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/spf13/cobra"
)

const (
	storageContract = "0xb81c8C0d691bA7A72704412cA0cF605427370Fd3"
)

func InitStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-storage",
		Short: "",
		Long:  "",
		Run:   initStorage,
	}
	flags := cmd.Flags()
	flags.IntVarP(&concurrency, "concurrency","g", 10, "set the number of tx number per second")
	flags.IntVarP(&sleepTimeTx, "sleepTime", "s",1, "")

	flags.StringVarP(&privkPath, "privkey-path", "p", "","")
	flags.StringSliceVarP(&rpc_hosts, "rpc-hosts","u", []string{}, "")
	flags.StringVarP(&chainID, "chain-id", "i","", "")

	return cmd
}

func initStorage(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	for i := 0; i < concurrency; i++ {
		go func(index int, privkey string) {
			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			deploy(privkey, rpcHost)

		}(i, privkeys[i])
	}

	select {}
}

func deploy(privkey string, host string) {
	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 30000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	addr := getCosmosAddress(privkey)
	accInfo, err := cli.Auth().QueryAccount(addr.String())
	if err != nil {
		panic(err)
	}
	nonce := accInfo.GetSequence()

	for {
		res, err := cli.Evm().SendTxEthereum(privkey, storageContract, "", "0xfe4b84df0000000000000000000000000000000000000000000000000000000000000002",20000000, nonce)
		if err != nil {
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}