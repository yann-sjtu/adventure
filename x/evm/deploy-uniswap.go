package evm

import "github.com/spf13/cobra"

const (
	flagMnemonic = "MnemonicPath"
)

func deployUniswapv2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-erc20-tokens",
		Short: "arbitrage token from swap and orderdepthbook",
		Args:  cobra.NoArgs,
		RunE:   deployUniswapv2,
	}
	cmd.Flags().StringP(flagMnemonic, "m", "", "the MnemonicPath of an account to deploy uniswap v2")

	return cmd
}

func deployUniswapv2(cmd *cobra.Command, args []string) error {
	//init flag
	_, err := cmd.Flags().GetString(flagMnemonic)
	if err != nil {
		return err
	}

	// change

	return nil
}