package evm

import (
	"github.com/okex/adventure/x/strategy/evm/template/ERC721"
	"github.com/okex/adventure/x/strategy/evm/template/USDT"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2Staker"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	InitTemplate()

	var evmCmd = &cobra.Command{
		Use:   "evm",
		Short: "evm cli of test strategy",
	}

	evmCmd.AddCommand(deployErc20Cmd(), uniswapTestCmd())
	return evmCmd
}

func InitTemplate() {
	UniswapV2.Init()
	ERC721.Init()
	USDT.Init()
	UniswapV2Staker.Init()
}