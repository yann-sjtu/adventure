package farm_control

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	"github.com/spf13/cobra"
)

var (
	sleepTime int
)

func FarmControlCmd() *cobra.Command {
	farmControlCmd := &cobra.Command{
		Use:   "farm-control",
		Short: "farm controlling in the first miner on a farm pool",
		Args:  cobra.NoArgs,
		RunE:  runFarmControlCmd,
	}

	flags := farmControlCmd.Flags()
	flags.IntVarP(&sleepTime, "sleep_time", "s", 30, "")
	flags.IntVarP(&startIndex, "start_index", "i", 901, "")
	flags.StringVar(&poolName, "pool_name", "1st_pool_okt_usdt", "")
	flags.StringVar(&lockSymbol, "lock_symbol", "ammswap-okt_usdt-a2b", "")
	flags.StringVar(&baseCoin, "base_coin", "okt", "")
	flags.StringVar(&quoteCoin, "quote_coin", "usdt-a2b", "")
	return farmControlCmd
}

var (
	startIndex = 0

	poolName   = ""
	lockSymbol = ""

	baseCoin  = ""
	quoteCoin = ""
)

func initGlobalParam() {
	limitRatio = types.MustNewDecFromStr("0.70")
	numerator = types.MustNewDecFromStr("3.0")
	denominator = types.MustNewDecFromStr("10.0")

	zeroLpt = types.NewDecCoinFromDec(lockSymbol, types.ZeroDec())
	zeroQuoteAmount = types.NewDecCoinFromDec(quoteCoin, types.ZeroDec())
}

func runFarmControlCmd(cmd *cobra.Command, args []string) error {
	initGlobalParam()
	accounts := monitorcommon.AddrsBook[startIndex/100]
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	for _, account := range accounts {
		fmt.Println()
		log.Printf("======================================== %+v ========================================\n", account)
		cli := clientManager.GetClient()

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
		} else {
			// 2.1 our_total_locked_lpt / total_locked_lpt < 80%, then promote the ratio over 81%
			err = replenishLockedToken(cli, account, requiredToken)
			if err != nil {
				fmt.Printf("[Phase2 Replenish] failed: %s\n", err.Error())
			}
		}
	}

	return nil
}
