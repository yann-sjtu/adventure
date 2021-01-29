package final_top_21

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/final_top_21_control/constant"
	"github.com/okex/adventure/x/monitor/final_top_21_control/keeper"
	"github.com/okex/adventure/x/monitor/final_top_21_control/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	flagTomlFilePath = "toml-path"
)

func FinalTop21SharesControlCmd() *cobra.Command {
	finalTop21SharesControlCmd := &cobra.Command{
		Use:   "final-top21-shares-control",
		Short: "shares controlling over top21 target validators",
		Args:  cobra.NoArgs,
		RunE:  runFinalTop21SharesControlCmd,
	}

	flags := finalTop21SharesControlCmd.Flags()
	flags.StringP(flagTomlFilePath, "p", "./config.toml", "the file path of config.toml")

	return finalTop21SharesControlCmd
}

func runFinalTop21SharesControlCmd(cmd *cobra.Command, args []string) error {
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

		// 1.found the intruder(stranger in top21, neither target vals and enemies)
		intruders := kp.CatchTheIntruders()
		if len(intruders) == 0 {
			log.Println("no intruders and everything goes well")
			continue
		}

		// 2. generate tokens to deposit
		tokenToDeposit := utils.GenerateRandomTokensToDeposit(500, 1000)

		// 3. pick a worker that has enough balance for tokenToDeposit
		worker, err := kp.PickEfficientWorker(tokenToDeposit)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 4. send msg
		if 	err := kp.SendMsgs(worker, tokenToDeposit); err != nil {
			log.Println(err.Error())
			continue
		}
	}
}
