package approve

import (
	"log"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func ApproveTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "",
		Args:  cobra.NoArgs,
		RunE:  approveCoins,
	}

	flags := cmd.Flags()
	flags.StringVarP(&FromAddr, "from-addr", "f", "", "set the address of erc20")
	flags.StringVarP(&PermittedAddr, "permitted-addr", "o", "", "set the address of approval contract")
	flags.StringVarP(&MnemonicPath, "mnemonic-path", "p", "", "set the addresses path to be transfered coins")

	return cmd
}

var (
	FromAddr      = ""
	PermittedAddr = ""
	MnemonicPath  = ""
)

func approveCoins(cmd *cobra.Command, args []string) error {
	clis := common.NewClientManager(common.Cfg.Hosts, "0.003okt", 300000)
	cli := clis.GetClient()

	accs := common.GetAccountManagerFromFile(MnemonicPath)
	for i, info := range accs.GetInfos() {
		acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
		if err != nil {
			return err
		}
		ethAddr := utils.GetEthAddressStrFromCosmosAddr(info.GetAddress())

		payload := UniswapV2.BuildApprovePayload(PermittedAddr, 1000000000000)
		res, err := cli.Evm().SendTx(info, common.PassWord, FromAddr, "", ethcommon.Bytes2Hex(payload), "", acc.GetAccountNumber(), acc.GetSequence()+uint64(i))
		if err != nil {
			log.Printf("%s fail to approve 1000000000000coin on %s from %s: %s \n", ethAddr, PermittedAddr, FromAddr, err)
			return err
		}
		log.Printf("[%s]%s send approve 1000000000000coin on %s from %s\n", res.TxHash, ethAddr, PermittedAddr, FromAddr)
	}

	return nil
}