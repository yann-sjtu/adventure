package farm_account_query

import (
	"fmt"
	"log"

	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	"github.com/spf13/cobra"
)

func FarmQueryCmd() *cobra.Command {
	farmUnlockCmd := &cobra.Command{
		Use:   "farm-account-query",
		Short: "farm-account-query",
		Args:  cobra.NoArgs,
		RunE:  runFarmQueryCmd,
	}

	flags := farmUnlockCmd.Flags()
	flags.IntVarP(&startIndex, "start_index", "i", 901, "")
	return farmUnlockCmd
}

var (
	startIndex = 0
)

func runFarmQueryCmd(cmd *cobra.Command, args []string) error {
	addrs := monitorcommon.AddrsBook[startIndex/100]

	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	cli := clientManager.GetClient()
	for i := 0; i < len(addrs); i++ {
		acc, err := cli.Auth().QueryAccount(addrs[i])
		if err != nil {
			log.Printf("[%d]%s failed to query account info. err: %s", startIndex+i, addrs[i], err.Error())
			continue
		}
		fmt.Println(startIndex+i, addrs[i], acc.GetCoins())
	}
	return nil
}
