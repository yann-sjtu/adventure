package shares_control

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	flagValNumberInTop21 = "val_in_top_21"
	flagRewardsPercent   = "rewards_percentage"
)

func SharesControlCmd() *cobra.Command {
	sharesControlCmd := &cobra.Command{
		Use:   "shares-control",
		Short: "shares controlling over all target validators",
		Args:  cobra.NoArgs,
		RunE:  runSharesControlCmd,
	}

	flags := sharesControlCmd.Flags()
	flags.Int64P(flagValNumberInTop21, "n", 21, "the number of validators in top 21")
	flags.Int64P(flagRewardsPercent, "p", 100, "the percentage of rewards to plunder")

	return sharesControlCmd
}

func runSharesControlCmd(cmd *cobra.Command, args []string) error {
	nValInTop21, err := cmd.Flags().GetInt64(flagValNumberInTop21)
	if err != nil {
		return err
	}

	percentToPlunder, err := cmd.Flags().GetInt64(flagRewardsPercent)
	if err != nil {
		return err
	}

	fmt.Println(nValInTop21, percentToPlunder)
	return nil
}
