package top21

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/cval_control/constant"
	"github.com/okex/adventure/x/monitor/top21_shares_control/keeper"
	"github.com/spf13/cobra"
)

func Top21SharesControlCmd() *cobra.Command {
	top21SharesControlCmd := &cobra.Command{
		Use:   "top21-shares-control",
		Short: "shares controlling over top21 target validators",
		Args:  cobra.NoArgs,
		RunE:  runTop21SharesControlCmd,
	}

	flags := top21SharesControlCmd.Flags()
	flags.StringP(constant.FlagTomlFilePath, "p", "./config.toml", "the file path of config.toml")

	return top21SharesControlCmd
}

func runTop21SharesControlCmd(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(constant.FlagTomlFilePath)
	if err != nil {
		return err
	}

	kp := keeper.NewKeeper()
	err = kp.Init(path)
	if err != nil {
		return err
	}

	fmt.Println(kp)

	//var round int
	//for {
	//	round++
	//	fmt.Printf("============================== Round %d ==============================\n", round)
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

	return nil
}
