package multi_transfer

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	round        int
	privkey      string
	contractAddr string
	sleepTime    int
	mode         string
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test-gas",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.IntVarP(&round, "round", "r", 1, "set Goroutine Num")
	flags.StringVarP(&privkey, "private-key", "p", "", "set the Priv Key path")
	flags.StringVarP(&contractAddr, "contract", "c", "", "set the Priv Key path")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 0, "set the sleep time")
	flags.StringVarP(&mode, "mode", "s", "sync", "set the mode of sync or block")
	return cmd
}

func testLoop(cmd *cobra.Command, args []string) {
	Init()
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.015okt", mode, 1500000)
	cli := clients.GetClient()
	info, err := utils.CreateAccountWithPrivateKey(privkey, "acc", common.PassWord)
	if err != nil {
		panic(err)
	}
	ethAddr, err := utils.ToHexAddress(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	payloadStr := hexutil.Encode(payload)

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	offset := uint64(0)

	for i := 0; i < round; i++ {
		res, err := cli.Evm().SendTxEthereum(privkey, contractAddr, "5.0", payloadStr, 2000000000, accInfo.GetSequence()+offset)
		if err != nil {
			log.Printf("[%s] %s failed to excute dyf in %s: %s\n", res.TxHash, ethAddr, contractAddr, err)
			continue
		} else {
			log.Printf("[%s] %s excute successfull in %s \n", res.TxHash, ethAddr, contractAddr)
			offset++
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}
