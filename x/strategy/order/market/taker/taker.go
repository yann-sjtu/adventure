package taker

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/okex/adventure/x/strategy/order/market/types"
	"github.com/spf13/cobra"
)

var (
	accountPath = ""

	sleepTime = 0
)

func TakerCmd() *cobra.Command {
	takerCmd := &cobra.Command{
		Use: "taker",
		Run: runTaker,
	}
	flags := takerCmd.Flags()
	flags.StringVarP(&types.Product, "product", "p", "", "set order product name")
	flags.StringVarP(&types.QueryProduct, "query_product", "q", "", "set query product name")
	flags.StringVarP(&accountPath, "account_path", "a", "", "set account mnemonic path")
	flags.IntVarP(&sleepTime, "sleep_time", "t", 0, "set sleep time")
	return takerCmd
}

type activeAccount struct {
	info      keys.Info
	buyOrder  types.Order
	sellOrder types.Order
}

func runTaker(cmd *cobra.Command, args []string) {
	// init account
	accountManager := common.GetAccountManagerFromFile(accountPath)
	infos := accountManager.GetInfos()
	activeAccounts := initActiveAccount(infos)

	// init clients
	clientManager := common.NewClientManager(config.Cfg.Hosts, config.AUTO)

	// init price
	latestPrice := types.QueryOneTickerPrice(types.QueryProduct)
	if latestPrice == 0.0 {
		panic("the price fetched from okex is zero!")
	}

	// start multi taker
	for {
		//get price
		price := types.QueryOneTickerPrice(types.QueryProduct)
		if price == 0.0 {
			price = latestPrice
		}

		//place orders
		placeOrders(clientManager, activeAccounts, price)

		latestPrice = price
		fmt.Println()

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func initActiveAccount(infos []keys.Info) []activeAccount {
	var activeAccounts []activeAccount
	for _, info := range infos {
		tmpAccount := activeAccount{info: info, buyOrder: types.Order{OrderType: "BUY"}, sellOrder: types.Order{OrderType: "SELL"}}
		activeAccounts = append(activeAccounts, tmpAccount)
	}
	return activeAccounts
}

func placeOrders(clientManager *common.ClientManager, activeAccounts []activeAccount, price float64) {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(2) {
	case 0:
		price += float64(rand.Intn(6)) / 10000.0
	case 1:
		price -= float64(rand.Intn(6)) / 10000.0
	}
	if price < 0.0001 {
		price = 0.0001
	}

	var wg sync.WaitGroup
	for i, _ := range activeAccounts {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			account := activeAccounts[i]

			cli := clientManager.GetClient()
			// cancel the old orders
			types.CancelOrders(cli, account.info, getUselessIds(account.info.GetAddress().String()))
			// sumbit new orders
			orders := createTakerOrders(price)
			types.PlaceOrders(cli, account.info, &orders)
			activeAccounts[i].buyOrder = orders[0]
			activeAccounts[i].sellOrder = orders[1]
		}(i)
	}
	wg.Wait()
}

func getUselessIds(addr string) string {
	msg := types.QueryOrders(types.Product, addr)
	if msg == nil || msg.Data.Data == nil {
		return ""
	}

	var ids string
	i := 0
	for _, order := range msg.Data.Data {
		if i >= 30 {
			break
		}
		ids += order.OrderID + ","
		i++
	}

	return strings.TrimRight(ids, ",")
}

func getOldOrderIds(account activeAccount) string {
	ids := account.buyOrder.Id + "," + account.sellOrder.Id
	return ids
}

func createTakerOrders(price float64) []types.Order {
	rand.Seed(time.Now().Unix())
	quanties := float64(rand.Intn(10000)+1) * 0.0001
	takerOrderSlice := make([]types.Order, 2)
	takerOrderSlice[0] = types.Order{OrderType: types.BUY, Price: price, Quantity: quanties}
	takerOrderSlice[1] = types.Order{OrderType: types.SELL, Price: price, Quantity: quanties}
	return takerOrderSlice
}
