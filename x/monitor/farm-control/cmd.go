package farm_control

import (
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	"github.com/spf13/cobra"
)

func FarmControlCmd() *cobra.Command {
	farmControlCmd := &cobra.Command{
		Use:   "farm-control",
		Short: "farm controlling in the first miner on a farm pool",
		Args:  cobra.NoArgs,
		RunE:  runFarmControlCmd,
	}

	flags := farmControlCmd.Flags()
	flags.IntVarP(&sleepTime, "sleep_time", "s", 0, "sleep time of add-liquidity msg and lock msg")
	flags.IntVarP(&startIndex, "start_index", "i", 0, "account index")
	flags.StringVar(&poolName, "pool_name", "", "farm pool name")
	flags.StringVar(&lockSymbol, "lock_symbol", "", "token name used for locking into farm pool")
	flags.StringVar(&baseCoin, "base_coin", "", "base coin name in swap pool")
	flags.StringVar(&quoteCoin, "quote_coin", "", "quote coin name in swap pool")
	flags.BoolVarP(&toAddLiquidity, "to_add_liquidity", "a", true, "decide to add liquidity or not")
	flags.BoolVarP(&toLock, "to_lock", "l", true, "decide to lock lpt or not")

	return farmControlCmd
}

var (
	sleepTime = 0

	startIndex = 0

	poolName   = ""
	lockSymbol = ""
	baseCoin   = ""
	quoteCoin  = ""

	toAddLiquidity = true
	toLock         = true
)

func runFarmControlCmd(cmd *cobra.Command, args []string) error {
	accounts := monitorcommon.AddrsBook[startIndex/100]
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	for _, account := range accounts {
		cli := clientManager.GetClient()
		if accInfo, err := cli.Auth().QueryAccount(account.Address); err != nil {
			continue
		} else if accInfo.GetCoins().AmountOf(baseCoin).LT(types.MustNewDecFromStr("2")) {
			continue
		} else if accInfo.GetCoins().AmountOf(quoteCoin).LT(types.MustNewDecFromStr("2")) {
			continue
		}
		fmt.Println()
		log.Printf("=================== %+v ===================\n", account)

		// 1. check the ratio of (our_total_locked_lpt / total_locked_lpt), then return how many lpt to be replenished
		requiredToken, err := calculateReuiredAmount(cli, accounts)
		if err != nil {
			fmt.Printf("[Phase1 Calculate] failed: %s\n", err.Error())
			continue
		}

		// 2. judge if the requiredToken is zero or not
		if requiredToken.IsZero() {
			// 2.1 our_total_locked_lpt / total_locked_lpt > 80%, then do nothing
			fmt.Printf("This Round doesn't need to lock more %s \n", lockSymbol)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			fmt.Printf("there is %s to be replenished\n", requiredToken)
			// 2.1 our_total_locked_lpt / total_locked_lpt < 80%, then promote the ratio over 81%
			err = replenishLockedToken(cli, account.Index, account.Address)
			if err != nil {
				fmt.Printf("[Phase2 Replenish] failed: %s\n", err.Error())
				continue
			}
		}
	}

	return nil
}
