package approve_one_to_all

import (
	"fmt"
	"log"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func ApproveTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-one-to-all",
		Short: "",
		Args:  cobra.NoArgs,
		RunE:  approveCoins,
	}

	flags := cmd.Flags()
	flags.StringVarP(&FromContract, "from-contract", "f", "", "set the address of erc20")
	flags.StringVarP(&ToAddrPath, "to-addrs", "o", "", "set the address of approval contract")
	flags.StringVarP(&Mnemonic, "mnemonic", "m", "", "set the mnemonic to be transfered coins")

	return cmd
}

var (
	FromContract      = ""
	ToAddrPath = ""
	Mnemonic  = ""
)

func approveCoins(cmd *cobra.Command, args []string) error {
	clis := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.003okt", types.BroadcastSync, 300000)
	cli := clis.GetClient()

	info, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(Mnemonic), fmt.Sprintf("acc%d", 1), "12345678")
	if err != nil {
		log.Fatalln(err.Error())
	}
	acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		return err
	}
	accNum, seqNum := acc.GetAccountNumber(), acc.GetSequence()


	addrs := common.GetAccountAddressFromFile(ToAddrPath)
	for i, addr := range addrs {
		cosmosaddr, err := utils.ToCosmosAddress(addr)
		if err != nil {
			panic(err)
		}

		ethAddr := utils.GetEthAddressStrFromCosmosAddr(cosmosaddr)

		payload := UniswapV2.BuildApprovePayload(addr, 10000000000000000)
		res, err := cli.Evm().SendTx(info, common.PassWord, FromContract, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum+uint64(i))
		if err != nil {
			log.Printf("%s fail to approve sending 100000coin to %s from contract %s: %s \n", acc.GetAddress().String(), ethAddr, FromContract, err)
			return err
		}
		log.Printf("[%s]%s approve sending 100000coin to %s from contract %s\n", res.TxHash, acc.GetAddress().String(), ethAddr, FromContract)
	}

	return nil
}
