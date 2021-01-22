package shares_control

import (
	"fmt"
	"github.com/okex/adventure/x/monitor/shares-control/keeper"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	flagTomlFilePath = "toml-path"
)

func SharesControlCmd() *cobra.Command {
	sharesControlCmd := &cobra.Command{
		Use:   "shares-control",
		Short: "shares controlling over all target validators",
		Args:  cobra.NoArgs,
		RunE:  runSharesControlCmd,
	}

	flags := sharesControlCmd.Flags()
	flags.StringP(flagTomlFilePath, "p", "./config.toml", "the file path of config.toml")

	return sharesControlCmd
}

func runSharesControlCmd(cmd *cobra.Command, args []string) error {

	path, err := cmd.Flags().GetString(flagTomlFilePath)
	if err != nil {
		return err
	}

	kp := keeper.NewKeeper()
	err = kp.Init(path)
	if err != nil {
		return err
	}

	var round int
	for {
		round++
		fmt.Printf("============================== Round %d ==============================\n", round)
		// analyse shares
		result, err := kp.AnalyseShares()
		if err != nil {
			log.Println(err)
			continue
		}

		switch result.GetCode() {
		case 0:
		case 1:
		case 2:
		}

		time.Sleep(roundInterval)
	}
}
