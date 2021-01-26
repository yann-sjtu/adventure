package monitor

import (
	"github.com/okex/adventure/x/monitor/farm-control"
	farm_query "github.com/okex/adventure/x/monitor/farm-query"
	farmrmliquidity "github.com/okex/adventure/x/monitor/farm-rm-liquidity"
	farm_unlock "github.com/okex/adventure/x/monitor/farm-unlock"
	"github.com/okex/adventure/x/monitor/shares-control"
	"github.com/spf13/cobra"
)

func MonitorCmd() *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "gazing at the deep",
	}

	monitorCmd.AddCommand(
		shares_control.SharesControlCmd(),
		farm_control.FarmControlCmd(),
		farm_unlock.FarmUnlockCmd(),
		farmrmliquidity.FarmRemoveLiquidityCmd(),
		farm_query.FarmQueryCmd(),
	)

	return monitorCmd
}
