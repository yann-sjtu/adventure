package evm

import (
	"github.com/okex/adventure/x/evm/dyf"
	"github.com/okex/adventure/x/evm/mint"
	"github.com/okex/adventure/x/evm/template/DYF"
	"github.com/okex/adventure/x/evm/template/ERC721"
	"github.com/okex/adventure/x/evm/template/TTotken"
	"github.com/okex/adventure/x/evm/template/USDT"
	"github.com/okex/adventure/x/evm/template/UniswapV2"
	"github.com/okex/adventure/x/evm/template/UniswapV2Staker"
	"github.com/okex/adventure/x/evm/uniswap-mining"
	"github.com/okex/adventure/x/evm/uniswap-swap"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	InitTemplate()

	var evmCmd = &cobra.Command{
		Use:   "evm",
		Short: "evm cli of test strategy",
	}

	evmCmd.AddCommand(
		mint.MintCmd(),
		uniswap_mining.Cmd(),
		uniswap_swap.Cmd(),
		dyf.Cmd(),
	)
	return evmCmd
}

func InitTemplate() {
	UniswapV2.Init()
	ERC721.Init()
	USDT.Init()
	UniswapV2Staker.Init()
	TTotken.Init()
	DYF.Init()
}
