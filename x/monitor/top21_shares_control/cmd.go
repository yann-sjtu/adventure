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
		// 0.round init
		err := kp.InitRound()
		if err != nil {
			log.Println(err)
			continue
		}

		// 1.found the intruder(stranger in top21, neither target vals and enemies)
		intruders := kp.CatchTheIntruders()

		// 2. get the highest shares of intruders
		limitShares := kp.GetTheHighestShares(intruders)

		// 3. get the targets vals that will be promote
		valAddrsStrToPromote := kp.GetTargetValAddrsStrToPromote(limitShares)
		if len(valAddrsStrToPromote) == 0 {
			// no target val to promote
			time.Sleep(constant.RoundInterval)
			continue
		}

		// 4. get the shares to add to the valAddrsStrToPromote
		requiredShares := kp.GetSharesToPromote(valAddrsStrToPromote, limitShares)

		_ = requiredShares
		time.Sleep(constant.RoundInterval)
	}
}
