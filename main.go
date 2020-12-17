package main

import (
	"log"
	"os"

	"github.com/okex/adventure/account"
	"github.com/okex/adventure/version"
	"github.com/okex/adventure/x/strategy/ammswap/strategy"
	"github.com/okex/adventure/x/strategy/evm"
	"github.com/okex/adventure/x/strategy/farm"
	"github.com/okex/adventure/x/strategy/order/market"
	"github.com/okex/adventure/x/strategy/staking/validators"
	tokenCmd "github.com/okex/adventure/x/strategy/token/cmd"
	"github.com/spf13/cobra"
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}

func main() {
	cobra.EnableCommandSorting = false
	mainCmd := &cobra.Command{
		Use:   "adventure",
		Short: "A client tool for okchain",
		Long: `⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏
 .----------------.  .----------------.  .----------------.  .----------------.  .-----------------. .----------------.  .----------------.  .----------------.  .----------------. 
| .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. |
| |      __      | || |  ________    | || | ____   ____  | || |  _________   | || | ____  _____  | || |  _________   | || | _____  _____ | || |  _______     | || |  _________   | |
| |     /  \     | || | |_   ___ '.  | || ||_  _| |_  _| | || | |_   ___  |  | || ||_   \|_   _| | || | |  _   _  |  | || ||_   _||_   _|| || | |_   __ \    | || | |_   ___  |  | |
| |    / /\ \    | || |   | |   '. \ | || |  \ \   / /   | || |   | |_  \_|  | || |  |   \ | |   | || | |_/ | | \_|  | || |  | |    | |  | || |   | |__) |   | || |   | |_  \_|  | |
| |   / ____ \   | || |   | |    | | | || |   \ \ / /    | || |   |  _|  _   | || |  | |\ \| |   | || |     | |      | || |  | '    ' |  | || |   |  __ /    | || |   |  _|  _   | |
| | _/ /    \ \_ | || |  _| |___.' / | || |    \ ' /     | || |  _| |___/ |  | || | _| |_\   |_  | || |    _| |_     | || |   \ '--' /   | || |  _| |  \ \_  | || |  _| |___/ |  | |
| ||____|  |____|| || | |________.'  | || |     \_/      | || | |_________|  | || ||_____|\____| | || |   |_____|    | || |    '.__.'    | || | |____| |___| | || | |_________|  | |
| |              | || |              | || |              | || |              | || |              | || |              | || |              | || |              | || |              | |
| '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' |
'----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'
⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏ ⛏

adventure is a very powerful cli tool for OKChain. It supports JSON-file and Sub-command to stimulate transactions.
`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			return
		},
	}

	mainCmd.AddCommand(account.Cmd())
	mainCmd.AddCommand(TxCmd())
	mainCmd.AddCommand(validators.StakingCmd())
	mainCmd.AddCommand(market.OrderMarketCmd())
	mainCmd.AddCommand(strategy.SwapCmd())
	mainCmd.AddCommand(tokenCmd.TokenCmd())
	mainCmd.AddCommand(farm.FarmCmd())
	mainCmd.AddCommand(evm.EvmCmd())
	mainCmd.AddCommand(version.Cmd)
	if err := mainCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
