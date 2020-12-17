package farm

import (
	"github.com/okex/adventure/x/strategy/farm/locker"
	"github.com/okex/adventure/x/strategy/farm/pooler"
	"github.com/spf13/cobra"
)

func FarmCmd() *cobra.Command {
	var farmCmd = &cobra.Command{
		Use:   "farm",
		Short: "farm cli for system test",
	}

	farmCmd.AddCommand(
		pooler.PoolerCmd(),
		locker.LockerCmd(),
		allocateTokensFromRicherCmd(),
		issueTokensByNumCmd(),
	)

	return farmCmd
}
