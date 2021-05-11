package evm

import (
	"github.com/okex/adventure/x/strategy/evm-unfinish/uniswap-operate"
	"github.com/okex/adventure/x/strategy/evm/approve-all-to-one"
	approve_one_to_all "github.com/okex/adventure/x/strategy/evm/approve-one-to-all"
	"github.com/okex/adventure/x/strategy/evm/mint"
	"github.com/okex/adventure/x/strategy/evm/template/DYF"
	"github.com/okex/adventure/x/strategy/evm/template/ERC721"
	"github.com/okex/adventure/x/strategy/evm/template/TTotken"
	"github.com/okex/adventure/x/strategy/evm/template/USDT"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2Staker"
	"github.com/okex/adventure/x/strategy/evm/transfer"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	InitTemplate()

	var evmCmd = &cobra.Command{
		Use:   "evm",
		Short: "evm cli of test strategy",
	}

	evmCmd.AddCommand(
		uniswap_operate.UniswapTestCmd(),
		transfer.TransferErc20Cmd(),
		approve_all_to_one.ApproveTokenCmd(),
		approve_one_to_all.ApproveTokenCmd(),
		mint.MintCmd(),
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
