package pooler

import (
	"sync"

	"github.com/okex/adventure/x/strategy/farm/pooler/types"
	"github.com/spf13/cobra"
)

func createSwapPairCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pair",
		Short: "create swap pair with the token issued and okt",
		Args:  cobra.NoArgs,
		RunE:  runCreateSwapPair,
	}

	flags := cmd.Flags()
	flags.StringP(FlagIssuerFilePath, "p", "", "the file path of pooler mnemonics")

	return cmd
}

func runCreateSwapPair(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(FlagIssuerFilePath)
	if err != nil {
		return err
	}
	issuerManager := types.GetPoolerManager(path)

	// create swap pair
	var wg sync.WaitGroup
	for _, issuer := range issuerManager {
		wg.Add(1)
		go issuer.CreateSwapPairWithUSDK(&wg)
	}
	wg.Wait()

	return nil
}
