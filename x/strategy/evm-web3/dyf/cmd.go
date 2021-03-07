package dyf

import (
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/DYF"
	"github.com/okex/adventure/x/strategy/evm/tools"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	goroutineNum int
	privkeyPath  string
	sleepTime    int
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dyf",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.IntVarP(&goroutineNum, "goroutine-num", "g", 1, "set Goroutine Num")
	flags.StringVarP(&privkeyPath, "private-path", "p", "", "set the Priv Key path")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 0, "set the sleep time")

	return cmd
}

const dyfAddr = "0xd78e1680e93bF57580F299d75B364e591873a8d3"

func testLoop(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkeyPath)
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.015okt", types.BroadcastBlock, 1500000)
	succ, fail := tools.NewCounter(-1), tools.NewCounter(-1)

	var wg sync.WaitGroup
	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			privkey := privkeys[index]
			cli := clients.GetClient()
			info, err := utils.CreateAccountWithPrivateKey(privkey, "acc", common.PassWord)
			if err != nil {
				panic(err)
			}

			ethAddr, err  := utils.ToHexAddress(info.GetAddress().String())
			if err != nil {
				panic(err)
			}

			goPayloadStr := hexutil.Encode(DYF.BuildExcutePayload())
			for {
				accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
				if err != nil {
					time.Sleep(time.Second*15)
					continue
				}
				// Let Us GO GO GO !!!!!!
				// 1. add liquididy
				res, err := cli.Evm().SendTxEthereum(privkey, dyfAddr, "", goPayloadStr, 1500000, accInfo.GetSequence())
				if err != nil {
					log.Printf("(%d)[%s] %s failed to excute dyf in %s: %s\n", fail.Add(), res.TxHash, ethAddr, dyfAddr, err)
					continue
				} else {
					log.Printf("(%d)[%s] %s excute dyf successfull in %s \n", succ.Add(), res.TxHash, ethAddr, dyfAddr)
				}
				time.Sleep(time.Second*time.Duration(sleepTime))
			}
		}(i)
	}
	wg.Wait()
}
