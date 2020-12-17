package account

import (
	"fmt"

	"github.com/bartekn/go-bip39"
	"github.com/spf13/cobra"
)

const (
	mnemonicsCmdName    = "mnemonics"
	mnemonicEntropySize = 128
)

var (
	mnemonicsNumber uint
)

func getMnemonicCmd() *cobra.Command {
	// add flags
	mnemonicsCmd := &cobra.Command{
		Use:   mnemonicsCmdName,
		Short: "create mnemonics",
		Long:  "create the account mnemonics",
		RunE:  generateMnemonics,
	}
	flags := mnemonicsCmd.Flags()
	flags.UintVarP(&mnemonicsNumber, "number", "n", 1, "the number of the user's mnemonics")
	return mnemonicsCmd
}

func generateMnemonics(cmd *cobra.Command, args []string) error {
	var i uint
	for i = 0; i < mnemonicsNumber; i++ {
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			return err
		}

		mnemonic, err := bip39.NewMnemonic(entropySeed[:])
		if err != nil {
			return err
		}

		fmt.Println(mnemonic)
	}
	return nil
}
