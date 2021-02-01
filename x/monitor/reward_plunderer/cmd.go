package reward_plunderer

import (
	"github.com/okex/adventure/x/monitor/reward_plunderer/keeper"
	"github.com/spf13/cobra"
)

const (
	flagTomlFilePath = "toml-path"
)

func RewardPlundererCmd() *cobra.Command {
	rewardPlundererCmd := &cobra.Command{
		Use:   "reward-plunderer",
		Short: "plunder total okt reward of staking",
		Args:  cobra.NoArgs,
		RunE:  runRewardPlundererCmd,
	}

	flags := rewardPlundererCmd.Flags()
	flags.StringP(flagTomlFilePath, "p", "./config.toml", "the file path of config.toml")

	return rewardPlundererCmd
}

func runRewardPlundererCmd(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(flagTomlFilePath)
	if err != nil {
		return err
	}

	kp := keeper.NewKeeper()
	err = kp.Init(path)
	if err != nil {
		return err
	}

	return nil
}
