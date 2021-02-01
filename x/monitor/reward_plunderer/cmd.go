package reward_plunderer

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/final_top_21_control/utils"
	"github.com/okex/adventure/x/monitor/reward_plunderer/constant"
	"github.com/okex/adventure/x/monitor/reward_plunderer/keeper"
	"github.com/spf13/cobra"
	"log"
	"time"
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

	var round int
	for {
		time.Sleep(constant.RoundInterval)
		round++
		fmt.Printf("============================== Round %d ==============================\n", round)
		// 0.round init
		err := kp.InitRound()
		if err != nil {
			log.Println(err)
			continue
		}

		// 1. check warning of current plundered pct
		if !kp.CheckPlunderedPctWarning() {
			log.Println("all rewards are under control")
			continue
		}

		// 2. generate tokens to deposit
		tokenToDeposit := utils.GenerateRandomTokensToDeposit(1, 1000)

		// 3. pick a worker that has enough balance for tokenToDeposit
		worker, err := kp.PickEfficientWorker(tokenToDeposit)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 4. send msg
		if err := kp.InfoToDeposit(worker, tokenToDeposit); err != nil {
			log.Println(err.Error())
		}
	}
}
