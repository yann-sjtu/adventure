package tools

import (
	"github.com/spf13/cobra"
)

func ToolsCmd() *cobra.Command {
	toolsCmd := &cobra.Command{
		Use:   "tools",
		Short: "useful tools",
	}

	toolsCmd.AddCommand(
		valAddrToAccAddrCmd())

	return toolsCmd
}
