package pooler

import (
	"fmt"
	"log"
	"time"

	"github.com/okex/adventure/x/strategy/farm/constants"
	"github.com/okex/adventure/x/strategy/farm/emitter"
	poolertypes "github.com/okex/adventure/x/strategy/farm/pooler/types"
	"github.com/okex/adventure/x/strategy/farm/utils"
	"github.com/okex/okexchain-go-sdk"
	"github.com/spf13/cobra"
)

const (
	flagPoolerFilePath = "pooler"
)

func strategyPoolerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy-pooler",
		Short: "poolers' strategy that they can destroy, create and provide farm pools automatically",
		Args:  cobra.NoArgs,
		RunE:  runStrategyPoolerCmd,
	}

	flags := cmd.Flags()
	flags.StringP(flagPoolerFilePath, "p", "", "the file path of pooler mnemonics")

	return cmd
}

func runStrategyPoolerCmd(cmd *cobra.Command, args []string) error {
	// load pooler manager
	poolerPath, err := cmd.Flags().GetString(flagPoolerFilePath)
	if err != nil {
		return err
	}

	emt := emitter.NewEmitter(poolertypes.GetPoolerManager(poolerPath), nil)

	for {
		// 1.get the pools that are expired on current height
		expiredPools, currentHeight, err := utils.GetExpiredPoolsOnCurrentHeight()
		if err != nil {
			log.Println(err)
			continue
		}

		if len(expiredPools) != 0 {
			// if there r some pools expired
			logExpiredInfo(expiredPools, currentHeight)
			// 1.destroy or provide pool randomly
			if utils.GetRandomBool() {
				emt.DestroyPools(expiredPools)
			} else {
				emt.ProvidePools(expiredPools, currentHeight)
			}
			time.Sleep(constants.SleepSecondAfterOperationOfExpiredPools)
		}

		// 2.check the poolers which don't have a pool
		poolersWithoutPools, err := emt.GetPoolersWithoutPools()
		if err != nil {
			log.Println(err)
			continue
		}
		if len(poolersWithoutPools) != 0 {
			emt.CreateAndProvidePool(poolersWithoutPools, currentHeight)
		}

		time.Sleep(constants.SleepSecondPerRoundStrategyPooler)
	}
}

func logExpiredInfo(pools []gosdk.FarmPool, currentHeight int64) {
	fmt.Printf(`
============================================================
|          Pools are caught expired on height %d         |
============================================================

`, currentHeight)
	for i, pool := range pools {
		fmt.Println(formatPoolPrinter(i, pool))
	}
}

func formatPoolPrinter(index int, pool gosdk.FarmPool) string {
	return fmt.Sprintf(`FarmPool%d:
  Pool Name:		%s	
  Owner:		%s
  Min Lock Amount:	%s
  Deposit Amount:	%s
  Total Value Locked:	%s
  Yielded Token Infos:	%v
  Total Accumulated Rewards:	%s
`,
		index+1, pool.Name, pool.Owner, pool.MinLockAmount, pool.DepositAmount, pool.TotalValueLocked, pool.YieldedTokenInfos, pool.TotalAccumulatedRewards)
}
