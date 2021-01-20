package locker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/x/strategy/farm/emitter"
	"github.com/okex/adventure/x/strategy/farm/locker/types"
	"github.com/spf13/cobra"
)

func lockToFarmPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock-to [pool name] [lock amount]",
		Short: "lockers lock specific amount of tokens to a farm pool",
		Args:  cobra.ExactArgs(2),
		RunE:  runLockToFarmPoolCmd,
	}

	flags := cmd.Flags()
	flags.StringP(flagLockerFilePath, "p", "", "the file path of locker mnemonics")

	return cmd
}

func runLockToFarmPoolCmd(cmd *cobra.Command, args []string) error {
	lockerPath, err := cmd.Flags().GetString(flagLockerFilePath)
	if err != nil {
		return err
	}

	lockAmount, err := sdk.ParseDecCoin(args[1])
	if err != nil {
		return err
	}

	emt := emitter.NewEmitter(nil, types.GetLockerManager(lockerPath))
	emt.LockToFarmPool(args[0], lockAmount)

	return nil
}
