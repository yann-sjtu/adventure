package monitor

import (
	farm_control "github.com/okex/adventure/x/monitor/farm-control"
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
	)

	return monitorCmd
}
