package farm

import (
	"github.com/okex/adventure/x/farm/utils"
	"github.com/spf13/cobra"
)

const (
	flagTargetAddrsFilePath = "path"
	flagAddrsNumOneTime     = "num"
)

func allocateTokensFromRicherCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocate-tokens [coinsStr]",
		Short: "allocate tokens to all addresses in a file from the richer",
		Args:  cobra.ExactArgs(1),
		RunE:  runAllocateTokensFromRicher,
	}

	flags := cmd.Flags()
	flags.StringP(flagTargetAddrsFilePath, "p", "", "the file path of target addresses")
	flags.IntP(flagAddrsNumOneTime, "n", 50, "num of addresses in one group multi-send")

	return cmd
}

func runAllocateTokensFromRicher(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(flagTargetAddrsFilePath)
	if err != nil {
		return err
	}

	targetAddrsStrs, err := utils.GetTestAddrsFromFile(path)
	if err != nil {
		panic(err)
	}

	targetAccAddrs, err := utils.ParseAccAddrsFromAddrsStr(targetAddrsStrs)
	if err != nil {
		panic(err)
	}

	addrNumOneTime, err := cmd.Flags().GetInt(flagAddrsNumOneTime)
	if err != nil {
		return err
	}

	return utils.MultiSendByGroup(nil, utils.GetRicherKeyInfo(), args[0], targetAccAddrs, addrNumOneTime)
}
