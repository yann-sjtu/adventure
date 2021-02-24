package rest_test

import (
	"github.com/spf13/cobra"
)

func RestTestCmd() *cobra.Command {
	var restTestCmd = &cobra.Command{
		Use:   "rest-test",
		Short: "rest test",
	}

	restTestCmd.AddCommand(
		deployErc20Cmd(),
	)
	return restTestCmd
}
