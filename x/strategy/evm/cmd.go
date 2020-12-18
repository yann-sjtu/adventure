package evm

import "github.com/spf13/cobra"

func EvmCmd() *cobra.Command {
	var evmCmd = &cobra.Command{
		Use:   "evm",
		Short: "evm cli of test strategy",
	}

	evmCmd.AddCommand(deployErc20Cmd())
	return evmCmd
}
