package farm_rm_liquidity

import "github.com/spf13/cobra"

func FarmRemoveLiquidityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-liquidity",
		Short: "remove liquidity",
		Args:  cobra.NoArgs,
		RunE:  runFarmRemoveLiquidityCmd,
	}
}

func runFarmRemoveLiquidityCmd(cmd *cobra.Command, args []string) error {

	return nil
}
