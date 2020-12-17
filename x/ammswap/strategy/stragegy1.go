package strategy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/okex/okexchain-go-sdk/utils"
	orderkeeper "github.com/okex/okexchain/x/order/keeper"
	"github.com/spf13/cobra"
)

func arbitrageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "arbitrage",
		Short: "arbitrage token from swap and orderdepthbook",
		Args:  cobra.NoArgs,
		Run:   arbitrageLoop,
	}

	flags := cmd.Flags()
	flags.StringVarP(&btc, "btc_name", "b", "", "set btc_name")
	flags.StringVarP(&usdk, "usdk_name", "u", "", "set usdk_name")
	flags.StringVarP(&mnemonic, "mnemonic", "m", "skate tomato unusual mixed sunset network razor buyer donate much tuition maple", "set account mnemonic")
	//flags.Uint64VarP(&num, "num", "n", 1000, "set num of issusing token")

	return cmd
}

var (
	btc  = "btc-8bb"
	usdk = "usdk-739"

	mnemonic = ""
)

//nohup adventure order maker -p="btc-8bb_usdk-739" -q="btc_usdt"  -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -t >> ~/btc-8bb_usdk-739-maker.log 2>&1 &
//okchaincli  tx dex list --from captain --gas-prices="0.00000001okt" --gas "400000"  --base-asset btc-8bb --quote-asset usdk-739 -y -b block
func arbitrageLoop(cmd *cobra.Command, args []string) {
	//clientManager := common.NewClientManager(config.Cfg.Hosts, config.AUTO)
	clientManager := common.NewClientManager([]string{"http://127.0.0.1:26657"}, config.AUTO)
	info, _, _ := utils.CreateAccountWithMnemo(mnemonic, "test", "12345678")

	for {
		cli := clientManager.GetClient()

		// calculate how many okt requires for buying 1 btc
		lqBtc, _ := cli.AmmSwap().QuerySwapTokenPair(btc)
		lqBtcNum, _ := strconv.ParseFloat(lqBtc.BasePooledCoin.Amount.String(), 64)
		lqOkt1Num, _ := strconv.ParseFloat(lqBtc.QuotePooledCoin.Amount.String(), 64)
		oktNum := calculateSellTokenNum(1, lqOkt1Num, lqBtcNum)
		fmt.Printf("it requires %.4f okt for buying %.4f btc in ammswap\n", oktNum, 1.0)

		// calculate how many usdk requires for buying a specific number of okt
		lqUsdk, _ := cli.AmmSwap().QuerySwapTokenPair(usdk)
		lqUsdkNum, _ := strconv.ParseFloat(lqUsdk.BasePooledCoin.Amount.String(), 64)
		lqOkt2Num, _ := strconv.ParseFloat(lqUsdk.QuotePooledCoin.Amount.String(), 64)
		sellNum := calculateSellTokenNum(oktNum, lqUsdkNum, lqOkt2Num)
		fmt.Printf("it requires %.4f usdk for buying %.4f okt in ammswap\n", sellNum, oktNum)
		// sellNum := calculateSellTokenNum(1, lqUsdkNum, lqBtcNum)

		// query depth book about btc_usdk
		depthBook, _ := cli.Order().QueryDepthBook(btc + "_" + usdk)
		buyNum, lastIndex, lastNum := buy(depthBook.Bids, 1)
		fmt.Printf("it will get %.4f usdk for selling %.4f okt in ammswap\n", buyNum, 1.0)

		profit := buyNum - sellNum
		if profit > 0 {
			// sell usdk, then buy btc in swap
			accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
			if err != nil {
				fmt.Println(err)
				continue
			}

			sellNumStr := strconv.FormatFloat(sellNum, 'f', 4, 64)
			res, err := cli.AmmSwap().TokenSwap(info, "12345678", sellNumStr+usdk, "0.99"+btc, info.GetAddress().String(), "10s", "", accInfo.GetAccountNumber(), accInfo.GetSequence())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}

			// sell btc, then buy usdk in order book
			var products, orderTypes, prices, quantities string
			for i := 0; i <= lastIndex; i++ {
				products += btc + "_" + usdk + ","
				orderTypes += "BUY" + ","
				prices += depthBook.Bids[i].Price + ","
				if i != lastIndex {
					quantities += depthBook.Bids[i].Quantity + ","
				} else {
					lastNumStr := strconv.FormatFloat(lastNum, 'f', 4, 64)
					quantities += lastNumStr + ","
				}
			}
			res, err = cli.Order().NewOrders(info, "12345678",
				strings.TrimRight(products, ","), strings.TrimRight(orderTypes, ","),
				strings.TrimRight(prices, ","), strings.TrimRight(quantities, ","),
				"new orders", accInfo.GetAccountNumber()+1, accInfo.GetSequence()+1)
			if err != nil {
				fmt.Println(err)
				continue
			}

			ids, err := utils.GetOrderIDsFromResponse(&res)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("place orders:", ids)
		}
	}
}

func calculateBuyTokenNum(sellTokenNum, sellTokenLP, buyTokenLP float64) float64 {
	return buyTokenLP - buyTokenLP*sellTokenLP/(sellTokenLP+0.997*sellTokenNum)
}

func calculateSellTokenNum(buyTokenNum, sellTokenLP, buyTokenLP float64) float64 {
	return (buyTokenLP*sellTokenLP/(buyTokenLP-buyTokenNum) - sellTokenLP) / 0.997
}

func buy(bids []orderkeeper.BookResItem, vIn float64) (float64, int, float64) {
	var profit float64
	var lastIndex int
	for index, item := range bids {
		price, _ := strconv.ParseFloat(item.Price, 64)
		volume, _ := strconv.ParseFloat(item.Quantity, 64)
		if vIn > volume {
			vIn -= volume
			profit += price * volume
		} else {
			profit += price * vIn
			lastIndex = index
		}
	}
	return profit, lastIndex, vIn
}
