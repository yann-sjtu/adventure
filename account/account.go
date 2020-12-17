package account

import (
	"github.com/spf13/cobra"
)

const (
	accountCmdName = "account"
)

var accountCmd = &cobra.Command{
	Use:   accountCmdName,
	Short: "create account",
	Long:  "create user account",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func Cmd() *cobra.Command {
	// add sub command
	accountCmd.AddCommand(getMnemonicCmd())
	accountCmd.AddCommand(getAddressCmd())
	accountCmd.AddCommand(transferTokenCmd())
	//accountCmd.AddCommand(addressCmd())
	return accountCmd
}
