package final_top_21

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/final_top_21_control/constant"
	"github.com/okex/adventure/x/monitor/final_top_21_control/keeper"
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
		//
		//	// 2. get the highest shares of intruders
		//	limitShares := kp.GetTheHighestShares(intruders)
		//
		//	// 3. get the targets vals that will be promote
		//	valAddrsStrToPromote := kp.GetTargetValAddrsStrToPromote(limitShares)
		//	if len(valAddrsStrToPromote) == 0 {
		//		// no target val to promote
		//		continue
		//	}
		//
		//	// 4. get the shares to add to the valAddrsStrToPromote
		//	requiredShares := kp.GetSharesToPromote(valAddrsStrToPromote, limitShares)
		//
		//	// 5. pick a worker to promote vals
		//	workers := kp.PickWorker(valAddrsStrToPromote)
		//
		//	// 6. calculate tokens to deposit with the requiredShares
		//	tokensToDeposit := kp.CalculateTokenToDeposit(requiredShares)
		//
		//	// 7. pre check for deposit
		//	electedWorkers, err := kp.PrecheckWorker(workers, tokensToDeposit)
		//	if err != nil {
		//		log.Println(err)
		//		continue
		//	}
		//
		//	_ = electedWorkers
		//	// 8. info to broadcast
		//	//err := kp.InfoToDeposit()
	}
}
