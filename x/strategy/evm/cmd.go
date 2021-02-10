package evm

import (
	"github.com/okex/adventure/x/strategy/evm/approve"
	"github.com/okex/adventure/x/strategy/evm/deploy-contracts"
	"github.com/okex/adventure/x/strategy/evm/template/ERC721"
	"github.com/okex/adventure/x/strategy/evm/template/USDT"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2Staker"
	"github.com/okex/adventure/x/strategy/evm/transfer"
	"github.com/okex/adventure/x/strategy/evm/uniswap-operate"
	"github.com/spf13/cobra"
)

func EvmCmd() *cobra.Command {
	InitTemplate()

	var evmCmd = &cobra.Command{
		Use:   "evm",
		Short: "evm cli of test strategy",
	}

	evmCmd.AddCommand(deploy_contracts.DeployErc20Cmd(), uniswap_operate.UniswapTestCmd(), transfer.TransferErc20Cmd(), approve.ApproveTokenCmd())
	return evmCmd
}

func InitTemplate() {
	UniswapV2.Init()
	ERC721.Init()
	USDT.Init()
	UniswapV2Staker.Init()
}