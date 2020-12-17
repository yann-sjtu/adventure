package account

import (
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	gosdk "github.com/okex/okexchain-go-sdk"
	tokenTypes "github.com/okex/okexchain-go-sdk/module/token/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

const (
	flagMnemoFilePath    = "mnemo_file"
	flagAccountsFilePath = "acc_path"
	flagAmount           = "init_amount"
	flagRich             = "rich_mnemonic"
)

func transferTokenCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:  "send",
		Long: "send tokens to test accounts",
		RunE: transferTokenScript1,
	}
	sendCmd.Flags().StringP(flagMnemoFilePath, "m", "", "the file path of mnemo to test")
	sendCmd.Flags().StringP(flagAccountsFilePath, "p", "", "the file path of account to test")
	sendCmd.Flags().StringP(flagAmount, "a", "0.1"+common.NativeToken, "send the initilize fund to test account")
	sendCmd.Flags().StringP(flagRich, "r", "", "send the initilize fund to test account")
	return sendCmd
}

func transferTokenScript(cmd *cobra.Command, args []string) error {
	//init flag
	path, err := cmd.Flags().GetString(flagAccountsFilePath)
	if err != nil {
		return err
	}

	richMnemonic, err := cmd.Flags().GetString(flagRich)
	if err != nil {
		return err
	}
	rich, _, _ := utils.CreateAccountWithMnemo(richMnemonic, "captain", common.PassWord)
	log.Println("rich addr:", rich.GetAddress().String())

	// create addrs
	addrs := common.GetAccountAddressFromFile(path)
	sum := len(addrs)

	//create rpc client
	clients := common.NewClientManager(config.Cfg.Hosts, "0.08"+common.NativeToken, 7952591)

	// query acc
	cli := clients.GetClient()
	accInfo, err := cli.Auth().QueryAccount(rich.GetAddress().String())
	if err != nil {
		return err
	}

	// query swap pairs
	swapPairs, err := cli.AmmSwap().QuerySwapTokenPairs()
	if err != nil {
		return err
	}

	for i, swapPair := range swapPairs {
		addr := addrs[i%sum]
		name1, name2 := swapPair.BasePooledCoin.Denom, swapPair.QuotePooledCoin.Denom

		coinStr := "10000" + name1
		// assume a successful transfer
		if _, err := cli.Token().Send(rich, common.PassWord, addr, coinStr, "", accInfo.GetAccountNumber(), accInfo.GetSequence()+uint64(i*3)); err != nil {
			fmt.Printf("[%d] multi send %s to %s fail\n", i, coinStr, addr)
			fmt.Println(err)
		} else {
			fmt.Printf("[%d] multi send %s to %s successfully\n", i, coinStr, addr)
		}

		coinStr = "10000" + name2
		// assume a successful transfer
		if _, err := cli.Token().Send(rich, common.PassWord, addr, coinStr, "", accInfo.GetAccountNumber(), accInfo.GetSequence()+uint64(i*3+1)); err != nil {
			fmt.Printf("[%d] multi send %s to %s fail\n", i, coinStr, addr)
			fmt.Println(err)
		} else {
			fmt.Printf("[%d] multi send %s to %s successfully\n", i, coinStr, addr)
		}

		coinStr = "10000" + common.NativeToken
		// assume a successful transfer
		if _, err := cli.Token().Send(rich, common.PassWord, addr, coinStr, "", accInfo.GetAccountNumber(), accInfo.GetSequence()+uint64(i*3+2)); err != nil {
			fmt.Printf("[%d] multi send %s to %s fail\n", i, coinStr, addr)
			fmt.Println(err)
		} else {
			fmt.Printf("[%d] multi send %s to %s successfully\n", i, coinStr, addr)
		}
	}

	return nil
}

func transferTokenScript1(cmd *cobra.Command, args []string) error {
	// get address from mnemo file
	var addrs []string
	path, err := cmd.Flags().GetString(flagMnemoFilePath)
	if err != nil {
		return err
	}

	if len(path) != 0 {
		addrs = common.GetAccountAddressFromMnemoFile(path)
	} else {
		// get address from addrs file
		path, err = cmd.Flags().GetString(flagAccountsFilePath)
		if err != nil {
			return err
		}
		addrs = common.GetAccountAddressFromFile(path)
	}

	coinStr, err := cmd.Flags().GetString(flagAmount)
	if err != nil {
		return err
	}

	richMnemonic, err := cmd.Flags().GetString(flagRich)
	if err != nil {
		return err
	}

	sum := len(addrs)

	//create rpc client
	clients := common.NewClientManager(config.Cfg.Hosts, config.AUTO)
	group := sum / 1000
	for i := 0; i < group; i++ {
		log.Printf("prepare to multi send %s to account[%d,%d]\n", coinStr, i*1000, (i+1)*1000-1)
		err = topUp(addrs[i*1000:(i+1)*1000], coinStr, clients.GetRandomClient(), richMnemonic)
		if err != nil {
			return err
		}
		log.Printf("multi send %s to account[%d,%d] successfully\n", coinStr, i*1000, (i+1)*1000-1)
		time.Sleep(2 * time.Second)
	}
	left := sum % 1000
	if left != 0 {
		startIndex := sum / 1000 * 1000
		log.Printf("prepare to multi send %s to account[%d,%d]\n", coinStr, startIndex, startIndex+left-1)
		err = topUp(addrs[startIndex:startIndex+left], coinStr, clients.GetRandomClient(), richMnemonic)
		if err != nil {
			return err
		}
		log.Printf("multi send %s to account[%d,%d] successfully\n", coinStr, startIndex, startIndex+left-1)
	}

	return nil
}

func topUp(accs []string, coinStr string, cli *gosdk.Client, mnemonic string) error {
	transferUnit, err := makeTransferUnit(accs, coinStr)
	if err != nil {
		return err
	}

	if mnemonic == "" {
		mnemonic = captainMnemonic
	}
	rich, _, _ := utils.CreateAccountWithMnemo(mnemonic, "captain", common.PassWord)
	log.Println("rich addr:", rich.GetAddress().String())
	accInfo, err := cli.Auth().QueryAccount(rich.GetAddress().String())
	if err != nil {
		return err
	}

	// assume a successful transfer
	if _, err := cli.Token().MultiSend(rich, common.PassWord, transferUnit, "", accInfo.GetAccountNumber(), accInfo.GetSequence()); err != nil {
		return err
	}
	return nil
}

func makeTransferUnit(accs []string, coinStr string) ([]tokenTypes.TransferUnit, error) {
	coins, err := types.ParseDecCoins(coinStr)
	if err != nil {
		return nil, err
	}

	accLen := len(accs)
	transferUnits := make([]tokenTypes.TransferUnit, accLen)
	for i := 0; i < accLen; i++ {
		accAddr, _ := types.AccAddressFromBech32(accs[i])
		transferUnits[i] = tokenTypes.TransferUnit{To: accAddr, Coins: coins}
	}

	return transferUnits, nil
}
