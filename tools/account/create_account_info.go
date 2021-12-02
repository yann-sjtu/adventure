package account

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

const (
	mnemoLabel      = "mnemonics"
	privKeyLabel    = "privkey"
	bech32AddrLabel = "bech32_addr"
	hexAddrLabel    = "hex_addr"
)

func createAccInfoCmd() *cobra.Command {
	// add flags
	addressCmd := &cobra.Command{
		Use:   "create-account-info",
		Short: "create account info with 4 kinds files",
		Long:  "create account info with output of mnemonics/privkeys/bech32 addresses/hex addresses files",
		RunE:  runCreateAccInfoCmd,
	}
	flags := addressCmd.Flags()
	flags.IntP(flagAccountNum, "x", 0, "the number of accounts to create")
	flags.StringP(flagOutputFolder, "f", "", "the output folder path of 4 kinds of files")
	return addressCmd
}

func runCreateAccInfoCmd(cmd *cobra.Command, args []string) error {
	n, err := cmd.Flags().GetInt(flagAccountNum)
	if err != nil {
		return err
	}

	dirPath, err := cmd.Flags().GetString(flagOutputFolder)
	if err != nil {
		return err
	}

	mnemonics := make([]string, n)
	infos := make([]keys.Info, n)
	fmt.Printf("creating %d accounts ...\n", n)
	for i := 0; i < n; i++ {
		accInfo, mnemo, err := utils.CreateAccount(fmt.Sprintf("acc%d", i), DefaultPassphrase)
		if err != nil {
			panic(err)
		}

		mnemonics[i] = mnemo
		infos[i] = accInfo
	}

	var wg sync.WaitGroup
	for _, label := range []string{mnemoLabel, privKeyLabel, bech32AddrLabel, hexAddrLabel} {
		wg.Add(1)
		go writeFile(&wg, label, dirPath, infos, mnemonics)
	}
	wg.Wait()

	return nil
}

func writeFile(wg *sync.WaitGroup, label, path string, accInfos []keys.Info, mnemonics []string) {
	defer wg.Done()
	switch label {
	case mnemoLabel:
		writeMnemonicFile(path, mnemonics)
	case privKeyLabel:
		writePrivKeyFile(path, mnemonics)
	case bech32AddrLabel:
		writeBech32AddrFile(path, accInfos)
	case hexAddrLabel:
		writeHexAddrFile(path, accInfos)
	default:
		panic("unsupported label")
	}
}

func writeHexAddrFile(path string, accInfos []keys.Info) {
	n := len(accInfos)
	hexAddrFilePath := filepath.Join(path, fmt.Sprintf("%s_%d.txt", hexAddrLabel, n))

	f, err := os.OpenFile(hexAddrFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	for i := 0; i < n; i++ {
		_, _ = w.WriteString(ethcmn.BytesToAddress(accInfos[i].GetAddress()).Hex())
		_ = w.WriteByte('\n')
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("done[%s file]\n", hexAddrLabel)
}

func writeBech32AddrFile(path string, accInfos []keys.Info) {
	n := len(accInfos)
	bech32AddrFilePath := filepath.Join(path, fmt.Sprintf("%s_%d.txt", bech32AddrLabel, n))

	f, err := os.OpenFile(bech32AddrFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	for i := 0; i < n; i++ {
		_, _ = w.WriteString(accInfos[i].GetAddress().String())
		_ = w.WriteByte('\n')
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("done[%s file]\n", bech32AddrLabel)
}

func writePrivKeyFile(path string, mnemonics []string) {
	n := len(mnemonics)
	mnemoFilePath := filepath.Join(path, fmt.Sprintf("%s_%d.txt", privKeyLabel, n))

	f, err := os.OpenFile(mnemoFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	for i := 0; i < n; i++ {
		privKey, err := utils.GeneratePrivateKeyFromMnemo(mnemonics[i])
		if err != nil {
			panic(err)
		}
		_, _ = w.WriteString(privKey)
		_ = w.WriteByte('\n')
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("done[%s file]\n", privKeyLabel)
}

func writeMnemonicFile(path string, mnemonics []string) {
	n := len(mnemonics)
	mnemoFilePath := filepath.Join(path, fmt.Sprintf("%s_%d.txt", mnemoLabel, n))

	f, err := os.OpenFile(mnemoFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	for i := 0; i < n; i++ {
		_, _ = w.WriteString(mnemonics[i])
		_ = w.WriteByte('\n')
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("done[%s file]\n", mnemoLabel)
}
