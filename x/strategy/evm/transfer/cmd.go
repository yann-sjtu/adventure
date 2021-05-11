package transfer

import (
	"fmt"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/USDT"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
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
	hosts := common.GlobalConfig.Networks[common.NetworkType].Hosts
	clis := common.NewClientManagerWithMode(hosts, "0.003okt", types.BroadcastSync, 300000)
	cli := clis.GetClient()

	privkey,err := utils.GeneratePrivateKeyFromMnemo(strings.TrimSpace(Owner))
	if err != nil {
		return err
	}

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
		payload := UniswapV2.BuildApprovePayload(ethAddr, 100)
		res, err := cli.Evm().SendTxEthereum(privkey, ErcAddr, "", ethcommon.Bytes2Hex(payload), 300000, acc.GetSequence()+uint64(i*2))
		if err != nil {
			log.Printf("%s fail to approve sending 100000coin to %s from contract %s: %s \n", ownerEthAddr, ethAddr, ErcAddr, err)
			return err
		}
		log.Printf("[%s]%s approve sending 100000coin to %s from contract %s\n", res.TxHash, ownerEthAddr, ethAddr, ErcAddr)


		payload = USDT.BuildUSDTTransferPayload(ethAddr, 1)
		res, err = cli.Evm().SendTxEthereum(privkey, ErcAddr, "", hexutil.Encode(payload), 300000, acc.GetSequence()+uint64(i*2+1))
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