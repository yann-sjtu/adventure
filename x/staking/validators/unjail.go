package validators

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

const flagAccountsFilePath = "path"

func getUnjailCmd() *cobra.Command {
	unbondCmd := &cobra.Command{
		Use:   "unjail",
		Short: "unjail validators",
		RunE:  runUnjailScript,
	}
	unbondCmd.Flags().StringP(flagAccountsFilePath, "p", "", "the file path of validators mnemonic")
	return unbondCmd
}

func runUnjailScript(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString(flagAccountsFilePath)
	if err != nil {
		return err
	}
	valManager := common.GetAccountManagerFromFile(path)
	valInfos := valManager.GetInfos()

	//create rpc client
	clientManager := common.NewClientManager(common.GlobalConfig.Networks[""].Hosts, common.AUTO)
	for i, info := range valInfos {
		client := clientManager.GetClient()
		addr := info.GetAddress().String()
		accInfo, err := client.Auth().QueryAccount(addr)
		if err != nil {
			log.Printf("account[%d]%s failed to get info. error:%s\n", i, addr, err.Error())
			continue
		}
		operAddr := sdk.ValAddress(info.GetAddress()).String()
		val, err := client.Staking().QueryValidator(operAddr)
		if err != nil {
			log.Printf("account[%d]%s failed to get val info. error:%s\n", i, addr, err.Error())
			continue
		}
		if val.Jailed == true {
			res, err := client.Slashing().Unjail(info, common.PassWord,
				"", accInfo.GetAccountNumber(), accInfo.GetSequence())
			if err != nil {
				log.Printf("account[%d]%s failed to send tx. error:%s\n", i, addr, err.Error())
				continue
			}
			log.Println(res)
		}
	}

	return nil
}
