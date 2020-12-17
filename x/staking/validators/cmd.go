package validators

import (
	"github.com/spf13/cobra"
)

func StakingCmd() *cobra.Command {
	var stakingCmd = &cobra.Command{
		Use:   "staking",
		Short: "staking cli about validators",
	}

	stakingCmd.AddCommand(valsLoopTestCmd())
	stakingCmd.AddCommand(createValidatorsCmd())
	stakingCmd.AddCommand(getUnjailCmd())
	stakingCmd.AddCommand(queryValidatorsCmd())
	return stakingCmd
}
