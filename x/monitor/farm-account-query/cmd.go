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
	accounts := monitorcommon.AddrsBook[startIndex/100]

	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	cli := clientManager.GetClient()
	for _, acc := range accounts {
		accInfo, err := cli.Auth().QueryAccount(acc.Address)
		if err != nil {
			log.Printf("[%d]%s failed to query account info. err: %s", acc.Index, acc.Address, err.Error())
			continue
		}
		fmt.Println(acc.Index, acc.Address, accInfo.GetCoins())
	}
	return nil
}
