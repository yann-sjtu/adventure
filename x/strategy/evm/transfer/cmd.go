package transfer

import (
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/USDT"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func TransferErc20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "arbitrage token from swap and orderdepthbook",
		Args:  cobra.NoArgs,
		RunE:   transferCoins,
	}

	flags := cmd.Flags()
	flags.StringVarP(&ErcAddr, "token-addr", "d", "", "set the address of erc20")
	flags.StringVarP(&Owner, "owner", "o", "", "set the mnemonic of erc20 owner")
	flags.StringVarP(&AccountsPath, "accounts-path", "p", "", "set the addresses path to be transfered coins")

	return cmd
}

var (
	ErcAddr = ""
	Owner   = ""
	AccountsPath = ""
)

func transferCoins(cmd *cobra.Command, args []string) error {
	clis := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.003okt", types.BroadcastSync, 300000)
	cli := clis.GetClient()

	ownerInfo, _, err := utils.CreateAccountWithMnemo(Owner, fmt.Sprintf("acc%d", 1), common.PassWord)
	if err != nil {
		panic(err)
	}
	acc, err := cli.Auth().QueryAccount(ownerInfo.GetAddress().String())
	if err != nil {
		return err
	}
	ownerAccAddress, err := sdk.AccAddressFromBech32(ownerInfo.GetAddress().String())
	if err != nil {
		panic(err)
	}
	ownerEthAddr, _ := utils.ToHexAddress(ownerAccAddress.String())

	addrs := common.GetAccountAddressFromFile(AccountsPath)
	ethAddrs := convertCosmosAddrsToEthAddrs(addrs)
	fmt.Printf("start send coins from contract %s\n", ErcAddr)
	for i, ethAddr := range ethAddrs {
		payload := USDT.BuildUSDTTransferPayload(ethAddr, 100)
		res, err := cli.Evm().SendTxEthereum("", ErcAddr, "", hexutil.Encode(payload), 300000, acc.GetSequence()+uint64(i))
		if err != nil {
			log.Printf("%s fail to send erc20 coins to %s: %s \n", ownerEthAddr, ethAddr, err)
			return err
		}
		log.Printf("[%s]%s send erc20 coins to %s\n", res.TxHash, ownerEthAddr, ethAddr)
	}

	return nil
}

func convertCosmosAddrsToEthAddrs(addrStrs []string) []string {
	ethAddrs := make([]string, len(addrStrs))
	for i, addrStr := range addrStrs {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			panic(err)
		}

		ethAddr, _ := utils.ToHexAddress(addr.String())
		ethAddrs[i] = ethAddr.String()
	}
	return ethAddrs
}