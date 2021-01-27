package top21

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/top21_shares_control/constant"
	"github.com/okex/adventure/x/monitor/top21_shares_control/keeper"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	flagTomlFilePath = "toml-path"
)

func Top21SharesControlCmd() *cobra.Command {
	top21SharesControlCmd := &cobra.Command{
		Use:   "top21-shares-control",
		Short: "shares controlling over top21 target validators",
		Args:  cobra.NoArgs,
		RunE:  runTop21SharesControlCmd,
	}

	flags := top21SharesControlCmd.Flags()
	flags.StringP(flagTomlFilePath, "p", "./config.toml", "the file path of config.toml")

	return top21SharesControlCmd
}

func runTop21SharesControlCmd(cmd *cobra.Command, args []string) error {
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
		round++
		fmt.Printf("============================== Round %d ==============================\n", round)
		// 1. sum shares
		enemyTotalShares, tarValsTotalShares, err := kp.SumShares()
		if err != nil {
			log.Println(err)
			continue
		}

		// 2. calculate how much shares to add
		kp.CalculateHowMuchToDeposit(enemyTotalShares, tarValsTotalShares)

		time.Sleep(constant.RoundInterval)
	}
	//	// analyse shares
	//	result, err := kp.AnalyseShares()
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//
	//	switch result.GetCode() {
	//	// TODO
	//	case 1:
	//		// TODO
	//	case 2, 3:
	//		if err = kp.RaisePercentageToPlunder(); err != nil {
	//			log.Println(err)
	//			continue
	//		}
	//
	//		fmt.Println("info to raise percentage to plunder successfully")
	//		fmt.Println("wait 1 minute ...")
	//		time.Sleep(constant.IntervalAfterTxBroadcast)
	//
	//	}
	//
	//	time.Sleep(constant.RoundInterval)
	//}
}
