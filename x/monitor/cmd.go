package monitor

import (
	"github.com/okex/adventure/x/monitor/farm-control"
	farm_rm_liquidity "github.com/okex/adventure/x/monitor/farm-rm-liquidity"
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
		farm_rm_liquidity.FarmRemoveLiquidityCmd(),
	)

	return monitorCmd
}
