package cval_control

import (
	"github.com/okex/adventure/x/monitor/cval_control/create_cvals"
	"github.com/spf13/cobra"
)

func CValControlCmd() *cobra.Command {
	cValControlCmd := &cobra.Command{
		Use:   "cval-control",
		Short: "candidate validators control",
	}

	cValControlCmd.AddCommand(
		create_cvals.CreateCValsCmd(),
	)

	return cValControlCmd
}
