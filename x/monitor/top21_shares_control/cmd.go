package top21

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/top21_shares_control/keeper"
	"github.com/okex/adventure/x/monitor/top21_shares_control/utils"
	"github.com/spf13/cobra"
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

	targetValAddrsStr, err := kp.GetTargetValsAddr(kp.GetEnemyValAddrs(), 2)
	if err != nil {
		return err
	}

	accAddrsStr, err := utils.ConvertValAddrsStr2AccAddrsStr(targetValAddrsStr)
	if err != nil {
		return err
	}

	fmt.Println(len(accAddrsStr))
	for _, a := range accAddrsStr {
		fmt.Println(a)
	}

	strs := []string{
		"okexchainvaloper1tkwxgcpvptua0q0h5tn0at58ufnjdue7kf5fvp",
		"okexchainvaloper18v23ln9ycrtg0mrwsm004sh4tdknudtddffjr5",
	}

	accAddrsStr, err = utils.ConvertValAddrsStr2AccAddrsStr(strs)
	if err != nil {
		return err
	}

	fmt.Println(len(accAddrsStr))
	for _, a := range accAddrsStr {
		fmt.Println(a)
	}

	return nil
	//var round int
	//for {
	//	round++
	//	fmt.Printf("============================== Round %d ==============================\n", round)
	//	// 1. sum shares
	//	enemyTotalShares, tarValsTotalShares, err := kp.SumShares()
	//	if err != nil {
	//		log.Println(err)
	//		continue
	//	}
	//
	//	// 2. calculate how much shares to add
	//	_, _, err = kp.CalculateHowMuchToDeposit(enemyTotalShares, tarValsTotalShares)
	//
	//	// 3. info to deposit
	//	time.Sleep(constant.RoundInterval)
	//}
}
