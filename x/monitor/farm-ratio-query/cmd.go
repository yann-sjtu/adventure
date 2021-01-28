package farm_ratio_query

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/spf13/cobra"
)

func FarmRatioQueryCmd() *cobra.Command {
	farmUnlockCmd := &cobra.Command{
		Use:   "farm-ratio-query",
		Short: "farm-ratio-query",
		Args:  cobra.NoArgs,
		RunE:  runFarmQueryCmd,
	}

	flags := farmUnlockCmd.Flags()
	flags.IntVarP(&startIndex, "start_index", "i", 901, "")
	flags.StringVarP(&poolName, "pool_name", "p", "", "")
	return farmUnlockCmd
}

var (
	startIndex = 0

	poolName = ""
)

func runFarmQueryCmd(cmd *cobra.Command, args []string) error {
	addrs := monitorcommon.AddrsBook[startIndex/100]
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	cli := clientManager.GetClient()

	// query total locked lpt
	totalLpt, err := queryFarmPool(cli, poolName)
	if err != nil {
		return err
	}
	lptName := totalLpt.Denom

	ourTotalLpt := types.NewDecCoinFromDec(lptName, types.ZeroDec())
	for i := 0; i < len(addrs); i++ {
		lockInfo, err := cli.Farm().QueryLockInfo(poolName, addrs[i])
		if err != nil {
			if strings.Contains(err.Error(), "hasn't locked") {
				fmt.Printf("[%d] %s lock %s\n", startIndex+i, addrs[i], types.NewDecCoinFromDec(lptName, types.ZeroDec()))
				continue
			} else {
				return fmt.Errorf("failed to query %s lock-info: %s", addrs[i], err.Error())
			}
		} else {
			ourTotalLpt = totalLpt.Add(lockInfo.Amount)
			fmt.Printf("[%d] %s lock %s\n", startIndex+i, addrs[i], lockInfo.Amount)
		}
	}

	fmt.Printf("ourTotalLpt: %s \n", ourTotalLpt)
	fmt.Printf("   totalLpt: %s \n", totalLpt)
	fmt.Printf("      ratio: %s \n", ourTotalLpt.Amount.Quo(totalLpt.Amount))
	return nil
}

func queryFarmPool(cli *gosdk.Client, poolName string) (types.DecCoin, error) {
	pool, err := cli.Farm().QueryPool(poolName)
	if err != nil {
		return types.DecCoin{}, err
	}
	fmt.Println("whole total locked:", pool.TotalValueLocked)
	return pool.TotalValueLocked, nil
}

