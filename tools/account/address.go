package account

import (
	"fmt"
	"os"

	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
)

const (
	flagAddressFilePath = "file"
	addressCmdName      = "address"
	workerNum           = 100

	captainName       = "captain"
	captainMnemonic   = "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	DefaultPassphrase = "12345678"
)

func getAddressCmd() *cobra.Command {
	// add flags
	addressCmd := &cobra.Command{
		Use:   addressCmdName,
		Short: "create address",
		Long:  "create address",
		RunE:  generateAddress,
	}
	flags := addressCmd.Flags()
	flags.StringP(flagAccountsFilePath, "p", "", "the file path of mnemonic")
	flags.StringP(flagAddressFilePath, "f", "", "the output file path of address")
	return addressCmd
}

func generateAddress(cmd *cobra.Command, args []string) error {
	mnemonicPath, err := cmd.Flags().GetString(flagAccountsFilePath)
	if err != nil {
		return err
	}
	accManager := common.GetAccountManagerFromFile(mnemonicPath)

	addressPath, err := cmd.Flags().GetString(flagAddressFilePath)
	if err != nil {
		return err
	}
	file, err := os.Create(addressPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, info := range accManager.GetInfos() {
		file.WriteString(info.GetAddress().String() + "\n")
	}
	return nil
}

//var (
//	addressNumber uint32
//	showMnemonics bool
//)
//
//var accountAddressCmd = &cobra.Command{
//	Use:   addressCmdName,
//	Short: fmt.Sprint(cmdShortDes),
//	Long:  fmt.Sprint(longDes),
//	RunE:  generateAddr,
//}
//
//func addressCmd() *cobra.Command {
//	// add flags
//	flags := accountAddressCmd.Flags()
//	flags.Uint32VarP(&addressNumber, "number", "n", 1, "the number of the address")
//	flags.BoolVarP(&showMnemonics, "show_mnemonics", "s", false, "whether to show the mnemonics")
//	return accountAddressCmd
//}

//func generateAddr(kb keys.Keybase) error {
//	accs := GenerateAccount(kb, addressNumber)
//	for i, info := range accs {
//		fmt.Println(i, info.GetAddress().String())
//	}
//	return nil
//}

// GenerateAccount generate the account(keys) using worker pool
func GenerateAccount(kb keys.Keybase, addressNumber uint32) []keys.Info {
	//kb := keys.NewInMemory()

	index := make(chan uint32, 5000)
	infoChan := make(chan keys.Info, 5000)

	//var wg sync.WaitGroup
	for i := 0; i < workerNum; i++ {
		//wg.Add(1)
		go worker(kb, index, infoChan)
	}
	//fmt.Println("addressNumber: ", addressNumber)
	var addrIdx uint32
	for addrIdx = 0; addrIdx < addressNumber; addrIdx++ {
		//fmt.Println("input: ", addrIdx)
		index <- addrIdx
	}

	accs := make([]keys.Info, 0)
	for i := 0; i < int(addressNumber); i++ {
		info := <-infoChan
		//fmt.Println(i, info.GetAddress().String())
		accs = append(accs, info)
	}

	close(index)
	//wg.Wait()
	return accs
}

// worker is the worker which generate the address
func worker(kb keys.Keybase, idx chan uint32, accs chan keys.Info) error {
	//defer wg.Done()
	encryptPassword := "12345678"
	for i := range idx {
		name := fmt.Sprintf("test%d", i)
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			return err
		}

		mnemonic, err := bip39.NewMnemonic(entropySeed[:])
		if err != nil {
			return err
		}
		info, err := kb.CreateAccount(name, mnemonic, "", encryptPassword, "", keys.Secp256k1)
		if err != nil {
			return err
		}
		//fmt.Println(i, info.GetAddress().String())
		accs <- info
	}
	return nil
}

func CreateCaptain(kb keys.Keybase) (info keys.Info, err error) {
	return kb.CreateAccount(captainName, captainMnemonic, "", DefaultPassphrase, "", keys.Secp256k1)
}

func CreateKeyInfoByNameAndMnemonic(kb keys.Keybase, name, mnemonic, passphrase string) (info keys.Info, err error) {
	return kb.CreateAccount(name, mnemonic, "", passphrase, "", keys.Secp256k1)
}
