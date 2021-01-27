package create_cvals

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/cval_control/constant"
	"github.com/spf13/cobra"
)

func CreateCValsCmd() *cobra.Command {
	createCValsCmd := &cobra.Command{
		Use:   "create-cvals",
		Short: "create candidate validators",
		Args:  cobra.NoArgs,
		RunE:  runCreateCValsCmd,
	}

	flags := createCValsCmd.Flags()
	flags.StringP(constant.FlagTomlFilePath, "p", "./config.toml", "the file path of config.toml")

	return createCValsCmd
}

func runCreateCValsCmd(cmd *cobra.Command, args []string) error {
	fmt.Println(123)
	return nil
}
