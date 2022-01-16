package wmt

import (
	"github.com/spf13/cobra"
)

func WMTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wmt",
		Short: "only used on OEC testnet",
		Run:   wmt,
	}

	return cmd
}
