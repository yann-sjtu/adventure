package pooler

import (
	"github.com/okex/adventure/x/farm/pooler/types"
	"github.com/spf13/cobra"
	"sync"
)

func addLiquidityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity",
		Short: "add liquidity to the swap pool and get the liquidity pool token",
		Args:  cobra.NoArgs,
		RunE:  runAddLiquidity,
	}

	flags := cmd.Flags()
	flags.StringP(FlagIssuerFilePath, "p", "", "the file path of pooler mnemonics")

	return cmd
}

func runAddLiquidity(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(FlagIssuerFilePath)
	if err != nil {
		return err
	}
	issuerManager := types.GetPoolerManager(path)

	// add liquidity
	var wg sync.WaitGroup
	for _, issuer := range issuerManager {
		wg.Add(1)
		go issuer.AddLiquidity(&wg)
	}
	wg.Wait()

	return nil
}
