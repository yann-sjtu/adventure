package evm_web3

import (
	"github.com/okex/adventure/x/strategy/evm-web3/dyf"
	"github.com/okex/adventure/x/strategy/evm-web3/uniswap-mining"
	"github.com/okex/adventure/x/strategy/evm-web3/uniswap-swap"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	var evmCmd = &cobra.Command{
		Use:   "evm-web3",
		Short: "evm web3 cli of test strategy",
	}

	evmCmd.AddCommand(
		uniswap_mining.Cmd(),
		uniswap_swap.Cmd(),
		dyf.Cmd(),
	)
	return evmCmd
}