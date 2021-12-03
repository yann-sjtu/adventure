package evm

import (
	"github.com/okex/adventure/x/evm/template/UniswapV2"
	"github.com/okex/adventure/x/evm/template/UniswapV2Staker"
	"github.com/okex/adventure/x/evm/uniswap-mining"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	InitTemplate()

	var evmCmd = &cobra.Command{
		Use:   "evm",
		Short: "evm cli of test strategy",
	}

	evmCmd.AddCommand(
		uniswap_mining.Cmd(),
	)
	return evmCmd
}

func InitTemplate() {
	UniswapV2.Init()
	UniswapV2Staker.Init()
}
