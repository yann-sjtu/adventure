package farm_statistics_reward

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/spf13/cobra"
)

func FarmRewardQueryCmd() *cobra.Command {
	farmUnlockCmd := &cobra.Command{
		Use:   "farm-reward-query",
		Short: "farm-reward-query",
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

	rewards := types.DecCoins{}
	for _, acc := range accounts {
		reward, err := cli.Farm().QueryEarnings(poolName, acc.Address)
		if err != nil {
			if strings.Contains(err.Error(), "hasn't locked") {
				//fmt.Printf("[%d] %s has no reward \n", acc.Index, acc.Address)
				continue
			}
			return err
		}
		rewards = rewards.Add(reward.AmountYielded...)
		fmt.Printf("[%d] %s get reward %s\n", acc.Index, acc.Address, reward.AmountYielded)
	}
	fmt.Println(rewards)
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


