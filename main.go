package main

import (
	"log"
	"os"
	"path"

	"github.com/okex/adventure/common"
	"github.com/okex/adventure/query"
	"github.com/okex/adventure/tools/account"
	"github.com/okex/adventure/tools/version"
	"github.com/okex/adventure/x/monitor"
	"github.com/okex/adventure/x/simple"
	"github.com/okex/adventure/x/strategy/evm"
	evmweb3 "github.com/okex/adventure/x/strategy/evm-web3"
	"github.com/okex/adventure/x/strategy/staking/validators"
	"github.com/okex/adventure/x/strategy/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigFlag  = "config"
	NetworkFlag = "network"
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

			common.InitConfig(viper.GetString(ConfigFlag))
			return
		},
	}

	mainCmd.AddCommand(
		monitor.MonitorCmd(),
		account.Cmd(),
		validators.StakingCmd(),
		token.TokenCmd(),
		evm.EvmCmd(),
		evmweb3.EvmCmd(),
		version.Cmd,
		query.BenchQueryCmd(),

		//TODO:
		simple.TxCmd(),
	)

	mainCmd.PersistentFlags().StringP(ConfigFlag, "c", path.Join(os.ExpandEnv("$HOME/.adventure"), "config.toml"),"setting of config path")
	mainCmd.PersistentFlags().StringP(NetworkFlag, "n", "","setting of network type")

	if err := mainCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
