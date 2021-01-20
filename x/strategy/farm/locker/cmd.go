package locker

import (
	"github.com/spf13/cobra"
)

func LockerCmd() *cobra.Command {
	lockerCmd := &cobra.Command{
		Use:   "locker",
		Short: "cases of lockers",
	}

	lockerCmd.AddCommand(
		allocateTokensToAllLockersFromAllPoolersCmd(),
		lockToFarmPoolCmd(),
		strategyLockAndUnlockCmd(),
		strategySinglePoolCmd(),
	)

	return lockerCmd
}
