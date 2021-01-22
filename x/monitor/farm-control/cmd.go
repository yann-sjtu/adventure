package farm_control

import (
	"time"

	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

func FarmControlCmd() *cobra.Command {
	farmControlCmd := &cobra.Command{
		Use:   "farm-control",
		Short: "farm controlling in the first miner on a farm pool",
		Args:  cobra.NoArgs,
		RunE:  runFarmControlCmd,
	}

	//flags := farmControlCmd.Flags()

	return farmControlCmd
}

func runFarmControlCmd(cmd *cobra.Command, args []string) error {
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	initFarmAccounts(clientManager.GetClient())

	for {
		time.Sleep(time.Second * 5)

		cli := clientManager.GetClient()

		// 1. check the ratio of (our_total_locked_lpt / total_locked_lpt)
		requiredToken, err := checkLockedRatio(cli)
		if err != nil {
			continue
		}
		if !requiredToken.IsZero() { // 2.1 our_total_locked_lpt / total_locked_lpt < 81%, then promote the ratio over 85%
			// 2.2 pick one addr, then query its own account
			account:= pickOneAccount()
			accInfo, err := cli.Auth().QueryAccount(account.Address)
			if err != nil {
				continue
			}

			// 2.3 there is not enough lpt in this addr, then add-liquidity in swap
			if accInfo.GetCoins().AmountOf(lockSymbol).LT(requiredToken.Amount) {
				// todo
			}

			// 2.4 lock lpt in the farm pool
			// todo

			// 2.5 update accounts
			account.IsLocked = true
		} else { // 2.1 our_total_locked_lpt / total_locked_lpt > 81%, then do nothing

		}
	}
}
