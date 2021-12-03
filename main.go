package main

import (
	"log"
	"os"
	"path"

	"github.com/okex/adventure/bench"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/tools/account"
	"github.com/okex/adventure/tools/version"
	"github.com/okex/adventure/x/evm"
	"github.com/okex/adventure/x/evm/evm-transfer"
	"github.com/okex/adventure/x/evm/query"
	"github.com/okex/adventure/x/staking/validators"
	"github.com/spf13/cobra"
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

			return
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			common.InitConfig(common.ConfigPath)
		},
	}

	mainCmd.AddCommand(
		evm.EvmCmd(),
		query.BenchQueryCmd(),
		evm_transfer.TransferCmd(),

		version.Cmd,

		account.Cmd(),
		validators.StakingCmd(),
		bench.BenchCmd(),
	)

	mainCmd.PersistentFlags().StringVarP(&common.ConfigPath, ConfigFlag, "c", path.Join(os.ExpandEnv("$HOME/.adventure"), "config.toml"), "setting of config path")
	mainCmd.PersistentFlags().StringVarP(&common.NetworkType, NetworkFlag, "n", "", "setting of network type")

	if err := mainCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
