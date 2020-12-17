package cmd

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

func issueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "issue token",
		Args:  cobra.NoArgs,
		Run:   issueTokenLoop,
	}
	flags := cmd.Flags()
	flags.StringVarP(&testCoinName, "token_name", "t", "btc", "set token pre name")
	flags.StringVarP(&host, "host", "u", "http://127.0.0.1:26657", "set host")
	flags.StringVarP(&mnemonicPath, "mnemonic_path", "p", "", "set account mnemonic")
	flags.Uint64VarP(&num, "num", "n", 1000, "set num of issusing token")

	return cmd
}

var (
	testCoinName = ""
	totalSupply  = "9000000000.00000000"

	host         = ""
	mnemonicPath = ""

	num = uint64(0)

	stdChars = []byte("abcdefghijklmnopqrstuvwxyz")
)

func issueTokenLoop(cmd *cobra.Command, args []string) {
	clis := common.NewClientManager(common.Cfg.Hosts, "auto")
	//cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastSync, "", 400000, 1.1, "0.00000001"+common.NativeToken)

	accs := common.GetAccountManagerFromFile(mnemonicPath)

	for _, info := range accs.GetInfos() {
		go func(info keys.Info) {
			cli := clis.GetClient()
			accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
			if err != nil {
				fmt.Println(err, common.Issue)
				return
			}
			accNum, seqNum := accInfo.GetAccountNumber(), accInfo.GetSequence()
			fmt.Println("accNum", accNum, "seqNum", seqNum)
			for i := uint64(0); i < num; i++ {
				name := getRandomString(3) + strconv.Itoa(rand.Intn(10))
				_, err := cli.Token().Issue(info, common.PassWord,
					name, name, totalSupply, "Used for test "+name+" "+strconv.Itoa(int(seqNum+i)),
					"", true, accNum, seqNum+i)
				//time.Sleep(10*time.Millisecond)
				if err != nil {
					fmt.Println(err, common.Issue, info)
					return
				}
				fmt.Println(i, common.Issue, name, " done")
			}
		}(info)
	}

	select {}
}

func getRandomString(length int) string {
	if length == 0 {
		return ""
	}
	clen := len(stdChars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for getRandomString()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = stdChars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
