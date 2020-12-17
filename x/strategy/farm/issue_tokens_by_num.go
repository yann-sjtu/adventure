package farm

import (
	"sync"

	"github.com/okex/adventure/x/strategy/farm/pooler"
	"github.com/okex/adventure/x/strategy/farm/pooler/types"
	"github.com/spf13/cobra"
)

const (
	flagNumber = "number"
)

func issueTokensByNumCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-tokens",
		Short: "issue tokens with an input number",
		Args:  cobra.NoArgs,
		RunE:  runIssueTokensByNum,
	}

	flags := cmd.Flags()
	flags.StringP(pooler.FlagIssuerFilePath, "p", "", "the file path of issuers mnemonics")
	flags.IntP(flagNumber, "n", -1, "the number of accounts to import")
	return cmd
}

func runIssueTokensByNum(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(pooler.FlagIssuerFilePath)
	if err != nil {
		return err
	}

	num, err := cmd.Flags().GetInt(flagNumber)
	if err != nil {
		return err
	}

	var issuerManager types.PoolerManager
	if num < 0 {
		// import all mnemonics
		issuerManager = types.GetPoolerManager(path)
	} else {
		// import mnemonics by number
		issuerManager = types.GetPoolerManager(path, num)
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
