package shares_control

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"time"
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
	flags.StringP(flagValNumberInTop21, "n", "21", "the number of validators in top 21")
	flags.StringP(flagRewardsPercent, "p", "0.8", "the percentage of rewards to plunder")

	return sharesControlCmd
}

func runSharesControlCmd(cmd *cobra.Command, args []string) error {
	nValInTop21Str, err := cmd.Flags().GetString(flagValNumberInTop21)
	if err != nil {
		return err
	}

	nValInTop21, err := sdk.NewDecFromStr(nValInTop21Str)
	if err != nil {
		return err
	}

	percentToPlunderStr, err := cmd.Flags().GetString(flagRewardsPercent)
	if err != nil {
		return err
	}

	percentToPlunder, err := sdk.NewDecFromStr(percentToPlunderStr)
	if err != nil {
		return err
	}

	for {
		fmt.Println(nValInTop21)
		fmt.Println(percentToPlunder)

		time.Sleep(roundInterval)
	}
}
