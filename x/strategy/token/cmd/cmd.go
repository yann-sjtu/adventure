package cmd

import "github.com/spf13/cobra"

func TokenCmd() *cobra.Command {
	var swapCmd = &cobra.Command{
		Use:   "token",
		Short: "token cli",
	}

	swapCmd.AddCommand(issueCmd())
	return swapCmd
}
