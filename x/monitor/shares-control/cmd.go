package shares_control

import (
	"github.com/spf13/cobra"
)

func SharesControlCmd() *cobra.Command {
	sharesControlCmd := &cobra.Command{
		Use:   "shares-control",
		Short: "shares controlling over all target validators",
		Args:  cobra.ExactArgs(2),
		RunE:  runSharesControlCmd,
	}

	return sharesControlCmd
}

func runSharesControlCmd(cmd *cobra.Command, args []string) error {
	return nil
}
