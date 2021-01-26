package farm_query

import (
	"fmt"
	"log"

	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	"github.com/spf13/cobra"
)

func FarmQueryCmd() *cobra.Command {
	farmUnlockCmd := &cobra.Command{
		Use:   "farm-query",
		Short: "farm-query",
		Args:  cobra.NoArgs,
		RunE:  runFarmQueryCmd,
	}

	//flags := farmUnlockCmd.Flags()
	return farmUnlockCmd
}

func runFarmQueryCmd(cmd *cobra.Command, args []string) error {
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	cli := clientManager.GetClient()
	for i := 0; i < len(addrs); i++ {
		acc, err := cli.Auth().QueryAccount(addrs[i])
		if err != nil {
			log.Printf("[%d]%s failed to query account info. err: %s", 901+i, addrs[i], err.Error())
			continue
		}
		fmt.Println(901+i, addrs[i], acc.GetCoins())
	}
	return nil
}

var addrs = monitorcommon.Addrs901To1000
