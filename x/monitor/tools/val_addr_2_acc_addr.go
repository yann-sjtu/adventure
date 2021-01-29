package tools

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func valAddrToAccAddrCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "v-a",
		Short: "convert validator addr to account validator",
		Args:  cobra.ExactArgs(1),
		RunE:  runValAddrToAccAddrCmd,
	}
}

func runValAddrToAccAddrCmd(cmd *cobra.Command, args []string) error {
	valAddr, err := sdk.ValAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	fmt.Println(sdk.AccAddress(valAddr).String())
	return nil
}
