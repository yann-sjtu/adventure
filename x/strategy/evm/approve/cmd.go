package approve

import (
	"log"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/okexchain-go-sdk/types"
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
	clis := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.003okt", types.BroadcastSync,300000)
	cli := clis.GetClient()

	accs := common.GetAccountManagerFromFile(MnemonicPath)
	for _, info := range accs.GetInfos() {
		acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
		if err != nil {
			return err
		}
		ethAddr := utils.GetEthAddressStrFromCosmosAddr(info.GetAddress())

		payload := UniswapV2.BuildApprovePayload(PermittedAddr, 100000)
		res, err := cli.Evm().SendTx(info, common.PassWord, FromAddr, "", ethcommon.Bytes2Hex(payload), "", acc.GetAccountNumber(), acc.GetSequence())
		if err != nil {
			log.Printf("%s fail to approve 100000coin on %s from %s: %s \n", ethAddr, PermittedAddr, FromAddr, err)
			return err
		}
		log.Printf("[%s]%s send approve 100000coin on %s from %s\n", res.TxHash, ethAddr, PermittedAddr, FromAddr)
	}

	return nil
}