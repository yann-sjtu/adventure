package main

import (
	"log"
	"os"
	"path"

	"github.com/okex/adventure/query"
	"github.com/okex/adventure/tools/account"
	"github.com/okex/adventure/tools/version"
	"github.com/okex/adventure/x/monitor"
	"github.com/okex/adventure/x/simple"
	"github.com/okex/adventure/x/strategy/ammswap/strategy"
	"github.com/okex/adventure/x/strategy/evm"
	evmweb3 "github.com/okex/adventure/x/strategy/evm-web3"
	"github.com/okex/adventure/x/strategy/farm"
	"github.com/okex/adventure/x/strategy/farm/client"
	"github.com/okex/adventure/x/strategy/order/market"
	"github.com/okex/adventure/x/strategy/staking/validators"
	"github.com/okex/adventure/x/strategy/token"
	"github.com/spf13/cobra"
)

const (
	ConfigFlag  = "config"
	NetworkFlag = "network"
)

var (
	defaultConfigPath = os.ExpandEnv("$HOME/.adventure")
)

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

	mainCmd.AddCommand(
		monitor.MonitorCmd(),
		client.LineBreak,
		account.Cmd(),
		simple.TxCmd(),
		validators.StakingCmd(),
		market.OrderMarketCmd(),
		strategy.SwapCmd(),
		token.TokenCmd(),
		farm.FarmCmd(),
		evm.EvmCmd(),
		evmweb3.EvmCmd(),
		version.Cmd,
		query.BenchQueryCmd(),
	)

	mainCmd.PersistentFlags().StringP(ConfigFlag, "c", path.Join(defaultConfigPath, "config.toml"),"setting of config path")
	mainCmd.PersistentFlags().StringP(NetworkFlag, "n", "","setting of network type")

	if err := mainCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
