package client

import (
	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

var (
	CliManager *common.ClientManager
	LineBreak  = &cobra.Command{Run: func(*cobra.Command, []string) {}}
)

func init() {
	CliManager = common.NewClientManager(common.Cfg.Hosts, common.AUTO)
}
