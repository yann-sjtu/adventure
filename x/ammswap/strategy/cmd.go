package strategy

import "github.com/spf13/cobra"

func SwapCmd() *cobra.Command {
	var swapCmd = &cobra.Command{
		Use:   "swap",
		Short: "swap cli about test strategy",
	}

	swapCmd.AddCommand(arbitrageCmd(), createCmd(), addSwapRemove())
	return swapCmd
}
