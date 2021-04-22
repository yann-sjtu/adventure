package strategy

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	mnemonic = ""

	issue_token_num = 0
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create token from swap and orderdepthbook",
		Args:  cobra.NoArgs,
		Run:   createLoop,
	}

	flags := cmd.Flags()
	flags.StringVarP(&mnemonic, "mnemonic", "m", "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool", "set account mnemonic")
	flags.IntVarP(&issue_token_num, "issue_num", "n", 1000, "set num of issusing token")

	return cmd
}

//nohup adventure order maker -p="btc-8bb_usdk-739" -q="btc_usdt"  -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -t >> ~/btc-8bb_usdk-739-maker.log 2>&1 &
//okchaincli  tx dex list --from captain --gas-prices="0.00000001okt" --gas "400000"  --base-asset btc-8bb --quote-asset usdk-739 -y -b block
func createLoop(cmd *cobra.Command, args []string) {
	info, _, err := utils.CreateAccountWithMnemo(mnemonic, fmt.Sprintf("acc"), "12345678")
	if err != nil {
		return
	}

	clis := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	cli := clis.GetClient()
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	accNum, seqNum := accInfo.GetAccountNumber(), accInfo.GetSequence()
	fmt.Println("accNum", accNum, "seqNum", seqNum)
	if err != nil {
		fmt.Println(err, common.CreateExchange)
		return
	}

	// TODO: issue tokens

	acc, _ := cli.Auth().QueryAccount(info.GetAddress().String())
	coins := acc.GetCoins()
	for i, token := range coins {
		var token2 sdk.Coin
		for {
			token2 = coins[rand.Intn(len(coins))]
			if token.Denom != token2.Denom {
				break
			}
		}
		_, err = cli.AmmSwap().CreateExchange(info, common.PassWord,
			token.Denom, token2.Denom,
			"", accNum, seqNum+uint64(i))
		if err != nil {
			fmt.Println(err, common.CreateExchange, info.GetAddress().String())
			return
		}
		fmt.Println(i, common.CreateExchange, token.Denom+"_"+token2.Denom, " done")
	}

}
