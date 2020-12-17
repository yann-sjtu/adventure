package pooler

import (
	"github.com/okex/adventure/x/farm/pooler/types"
	"github.com/okex/adventure/x/farm/utils"
	"github.com/spf13/cobra"
	"sync"
)

func issueTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "issue tokens",
		Args:  cobra.NoArgs,
		RunE:  runIssueTokens,
	}

	flags := cmd.Flags()
	flags.StringP(FlagIssuerFilePath, "p", "", "the file path of pooler mnemonics")

	return cmd
}

func runIssueTokens(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(FlagIssuerFilePath)
	if err != nil {
		return err
	}
	issuerManager := types.GetPoolerManager(path)

	// issue usdk by richer
	if err = utils.IssueStableCoin(); err != nil {
		panic(err)
	}

	// issue token
	var wg sync.WaitGroup
	for _, issuer := range issuerManager {
		wg.Add(1)
		go issuer.IssueToken(&wg)
	}
	wg.Wait()

	return nil
}
