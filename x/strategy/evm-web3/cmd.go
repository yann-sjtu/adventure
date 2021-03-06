package evm_web3

import (
	"github.com/okex/adventure/x/strategy/evm-web3/uniswap"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	var evmCmd = &cobra.Command{
		Use:   "evm-web3",
		Short: "evm web3 cli of test strategy",
	}

	evmCmd.AddCommand(
		uniswap.Cmd(),
	)
	return evmCmd
}