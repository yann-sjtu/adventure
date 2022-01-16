package batch_transfer

import (
	"github.com/spf13/cobra"
)

const (
	FlagPriavteKey  = "private-key"
	FlagAddressFile = "address-file"
)

func BatchTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-transfer",
		Short: "used for transfer native token to accounts (default to 2000 fixed addresses)",
		Long: `
used for transfer native token to accounts (default to 2000 fixed addresses)

Example:
  $ adventure evm batch-transfer 10 -i ${ip} -s ${private_key} -a ${address_file}
				`,
		Args: cobra.ExactArgs(1),
		Run:  batchTransfer,
	}

	cmd.Flags().StringVarP(&privateKey, FlagPriavteKey, "s", "", "its private key should be imported as a rich account")
	cmd.Flags().StringVarP(&addressFile, FlagAddressFile, "a", "", "the path of ethereum-format address file")
	return cmd
}
