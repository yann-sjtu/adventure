package pooler

import (
	"github.com/okex/adventure/x/strategy/farm/client"
	"github.com/spf13/cobra"
)

const (
	FlagIssuerFilePath = "path"
)

func PoolerCmd() *cobra.Command {
	poolerCmd := &cobra.Command{
		Use:   "pooler",
		Short: "cases of poolers",
	}

	poolerCmd.AddCommand(
		issueTokensCmd(),
		createSwapPairCmd(),
		addLiquidityCmd(),
		client.LineBreak,
		strategyPoolerCmd(),
	)

	return poolerCmd
}
