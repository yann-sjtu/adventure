package farm_unlock

import (
	"log"
	"time"

	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	farm_control "github.com/okex/adventure/x/monitor/farm-control"
	"github.com/spf13/cobra"
)

var (
	startIndex = 0

	poolName = ""
)

func FarmUnlockCmd() *cobra.Command {
	farmUnlockCmd := &cobra.Command{
		Use:   "farm-unlock",
		Short: "farm-unlock",
		Args:  cobra.NoArgs,
		RunE:  runFarmUnlocklCmd,
	}

	flags := farmUnlockCmd.Flags()
	flags.IntVarP(&startIndex, "start_index", "i", 901, "")
	flags.StringVar(&poolName, "pool_name", "1st_pool_okt_usdt", "")

	return farmUnlockCmd
}

func runFarmUnlocklCmd(cmd *cobra.Command, args []string) error {
	addrs := monitorcommon.AddrsBook[startIndex/100]
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	cli := clientManager.GetClient()
	for i := 0; i < len(addrs); i++ {
		index := i + startIndex
		lockinfo, err := cli.Farm().QueryLockInfo(poolName, addrs[i])
		if err != nil {
			log.Printf("[%d]%s failed to query lock info. err: %s", index, addrs[i], err.Error())
			continue
		}

		if !lockinfo.Amount.IsZero() {
			acc, err := cli.Auth().QueryAccount(addrs[i])
			if err != nil {
				log.Printf("[%d]%s failed to query account info. err: %s", index, addrs[i], err.Error())
				continue
			}

			unlockMsg := farm_control.NewMsgUnLock(acc.GetAccountNumber(), acc.GetSequence(), lockinfo.Amount, addrs[i])
			err = monitorcommon.SendMsg(monitorcommon.Unfarmlp, unlockMsg, index)
			if err != nil {
				log.Printf("[%d] %s failed to unlock: %s\n", index, addrs[i], err)
				continue
			}
			log.Printf("[%d] %s send unlock msg: %+v\n", index, addrs[i], unlockMsg.Msgs[0])
		}

		if i%20 == 0 && i != 0 {
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
	return nil
}
