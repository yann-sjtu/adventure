package locker

import (
	"github.com/okex/adventure/x/strategy/farm/constants"
	"github.com/okex/adventure/x/strategy/farm/emitter"
	"github.com/okex/adventure/x/strategy/farm/locker/types"
	"github.com/okex/adventure/x/strategy/farm/utils"
	"github.com/spf13/cobra"
	"time"
)

func strategySinglePoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy-single-pool",
		Short: "lockers' strategy that they can lock and unlock tokens to a single pool",
		Args:  cobra.ExactArgs(1),
		RunE:  runStrategySinglePoolCmd,
	}

	flags := cmd.Flags()
	flags.StringP(flagLockerFilePath, "p", "", "the file path of locker mnemonics")

	return cmd
}

func runStrategySinglePoolCmd(cmd *cobra.Command, args []string) error {
	// load locker manager
	lockerPath, err := cmd.Flags().GetString(flagLockerFilePath)
	if err != nil {
		return err
	}

	emt := emitter.NewEmitter(nil, types.GetLockerManager(lockerPath))
	targetPoolName := args[0]
	var counter int

	for {
		// 1. get target pool
		targetPool, err := utils.QueryTargetPool(targetPoolName)
		if err != nil {
			continue
		}

		// 2. pick lockers randomly with fixed number this round
		emt.PickLockersRandomly()

		// 3. lock/unlock from the selected lockers
		if counter%2 == 0 {
			emt.LockToTargetPool(targetPool)
		} else {
			emt.UnlockFromPoolsLockedBeforeByRandomLockers()
		}

		counter++
		time.Sleep(constants.SleepSecondPerRoundStrategyLockAndUnlock)
	}
}
