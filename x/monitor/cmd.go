package monitor

import (
	farm_query "github.com/okex/adventure/x/monitor/farm-account-query"
	"github.com/okex/adventure/x/monitor/farm-control"
	"github.com/okex/adventure/x/monitor/tools"
	"github.com/spf13/cobra"
)

func MonitorCmd() *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "gazing at the deep",
	}

	monitorCmd.AddCommand(
		farm_control.FarmControlCmd(),
		farm_query.FarmQueryCmd(),
		tools.ToolsCmd(),
	)

	return monitorCmd
}
