package transfer

import "github.com/spf13/cobra"

func TransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "send native token to address",
		Run:   transfer,
	}

	cmd.Flags().BoolVarP(&fixed, "fixed", "f", false, "if true, transfer to one address; otherwise, transfer to a random address")
	return cmd
}
