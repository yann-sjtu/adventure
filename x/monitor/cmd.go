package monitor

import (
	"github.com/okex/adventure/x/monitor/cval_control"
	farm_query "github.com/okex/adventure/x/monitor/farm-account-query"
	"github.com/okex/adventure/x/monitor/farm-control"
	farm_ratio_query "github.com/okex/adventure/x/monitor/farm-ratio-query"
	farmrmliquidity "github.com/okex/adventure/x/monitor/farm-rm-liquidity"
	farm_statistics_reward "github.com/okex/adventure/x/monitor/farm-statistics-reward"
	farm_unlock "github.com/okex/adventure/x/monitor/farm-unlock"
	final_top_21 "github.com/okex/adventure/x/monitor/final_top_21_control"
	"github.com/okex/adventure/x/monitor/reward_plunderer"
	"github.com/okex/adventure/x/monitor/shares-control"
	"github.com/okex/adventure/x/monitor/tools"
	top21 "github.com/okex/adventure/x/monitor/top21_shares_control"
	"github.com/spf13/cobra"
)

func MonitorCmd() *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "gazing at the deep",
	}

	monitorCmd.AddCommand(
		reward_plunderer.RewardPlundererCmd(),
		final_top_21.FinalTop21SharesControlCmd(),
		top21.Top21SharesControlCmd(),
		shares_control.SharesControlCmd(),
		cval_control.CValControlCmd(),
		farm_control.FarmControlCmd(),
		farm_unlock.FarmUnlockCmd(),
		farmrmliquidity.FarmRemoveLiquidityCmd(),
		farm_query.FarmQueryCmd(),
		farm_ratio_query.FarmRatioQueryCmd(),
		tools.ToolsCmd(),
		farm_statistics_reward.FarmRewardQueryCmd(),
	)

	return monitorCmd
}
