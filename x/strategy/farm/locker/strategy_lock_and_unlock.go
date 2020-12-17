package locker

import (
	"log"
	"time"

	"github.com/okex/adventure/x/strategy/farm/constants"
	"github.com/okex/adventure/x/strategy/farm/emitter"
	lockertypes "github.com/okex/adventure/x/strategy/farm/locker/types"
	"github.com/okex/adventure/x/strategy/farm/utils"
	"github.com/spf13/cobra"
)

const (
	flagLockerFilePath = "path"
)

func strategyLockAndUnlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy-lock-unlock",
		Short: "lockers' strategy that they can lock and unlock tokens to a random pool",
		Args:  cobra.NoArgs,
		RunE:  runStrategyLockAndUnlockCmd,
	}

	flags := cmd.Flags()
	flags.StringP(flagLockerFilePath, "p", "", "the file path of locker mnemonics")

	return cmd
}

func runStrategyLockAndUnlockCmd(cmd *cobra.Command, args []string) error {
	// load locker manager
	lockerPath, err := cmd.Flags().GetString(flagLockerFilePath)
	if err != nil {
		return err
	}

	emt := emitter.NewEmitter(nil, lockertypes.GetLockerManager(lockerPath))

	for {
		// 1.unlock all token from expired pools
		if err = emt.UnlockFromExpiredPools(); err != nil {
			log.Println(err)
		}

		// 2.pick lockers randomly this round
		emt.PickLockersRandomly()

		// 3.pick target pools randomly this round
		targetPools, err := utils.GetPoolsRandomly()
		if err != nil {
			log.Println(err)
			continue
		}

		if len(targetPools) != 0 {
			// 4.all picked lockers lock some token to the pools picked this round
			emt.LockToRandomPoolsByRandomLockers(targetPools)
		}

		// 5.all picked lockers unlock some token from the pools that they have locked on before this round
		emt.UnlockFromPoolsLockedBeforeByRandomLockers()

		time.Sleep(constants.SleepSecondPerRoundStrategyLockAndUnlock)
	}

}
