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
		RunE:  runFarmRatioQueryCmd,
	}

	flags := farmUnlockCmd.Flags()
	flags.IntVarP(&startIndex, "start_index", "i", 0, "")
	flags.StringVarP(&poolName, "pool_name", "p", "", "")
	return farmUnlockCmd
}

var (
	startIndex = 0

	poolName = ""
)

func runFarmRatioQueryCmd(cmd *cobra.Command, args []string) error {
	accounts := monitorcommon.AddrsBook[startIndex/100]
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	cli := clientManager.GetClient()

	// query total locked lpt
	totalLpt, err := queryFarmPool(cli, poolName)
	if err != nil {
		return err
	}
	lptName := totalLpt.Denom

	ourTotalLpt := types.NewDecCoinFromDec(lptName, types.ZeroDec())
	for _, acc := range accounts {
		lockInfo, err := cli.Farm().QueryLockInfo(poolName, acc.Address)
		if err != nil {
			if strings.Contains(err.Error(), "hasn't locked") {
				fmt.Printf("[%d] %s lock %s\n", acc.Index, acc.Address, types.NewDecCoinFromDec(lptName, types.ZeroDec()))
				continue
			} else {
				return fmt.Errorf("failed to query %s lock-info: %s", acc.Address, err.Error())
			}
		} else {
			ourTotalLpt = ourTotalLpt.Add(lockInfo.Amount)
			fmt.Printf("[%d] %s lock %s\n", acc.Index, acc.Address, lockInfo.Amount)
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

