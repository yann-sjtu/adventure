package client

import (
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/spf13/cobra"
)

var (
	CliManager *common.ClientManager
	LineBreak  = &cobra.Command{Run: func(*cobra.Command, []string) {}}
)

func init() {
	CliManager = common.NewClientManager(config.Cfg.Hosts, "0.03"+common.NativeToken, 3000000)
}
