package maker

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/order/market/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	tinyMaker = false
)

func MakerCmd() *cobra.Command {
	makerCmd := &cobra.Command{
		Use: "maker",
		Run: runMaker,
	}
	flags := makerCmd.Flags()
	flags.StringVarP(&types.Product, "product", "p", "", "set coin product name")
	flags.StringVarP(&types.QueryProduct, "query_product", "q", "", "set query product name")
	flags.StringVarP(&types.Mnemonic, "mnemonic", "m", "", "set account mnemonic")
	flags.BoolVarP(&tinyMaker, "tiny_maker", "t", false, "set account mnemonic")
	return makerCmd
}

func runMaker(cmd *cobra.Command, args []string) {
	// init price
	priceList := initList()
	initPrice(priceList)

	// init maker account info
	info, _, err := utils.CreateAccountWithMnemo(types.Mnemonic, types.Name, types.PassWd)
	if err != nil {
		panic(err)
	}

	// init clients
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	// start
	t := time.NewTicker(time.Second * 15)
	for {
		<-t.C
		fmt.Println()
		log.Println("please wait 15 seconds")

		orders, ids := createMakerOrders(info, priceList)
		cli := clientManager.GetClient()
		types.CancelOrders(cli, info, ids)
		types.PlaceOrders(cli, info, &orders)
		setOrderIdInList(priceList, orders)
	}
}

var (
	basePrice  = 0.0
	levelPrice = 0.0

	length       = 20
	baseLevel    = 1000 //            base:1000
	perLevel     = 1000 //  buy: 0<-1000,   sell:10001->2000
	safeLevel    = 5
	safePercent  = 0.005
	levelPercent = 0.001
)

func initPrice(l *List) {
	basePrice := types.QueryOneTickerPrice(types.QueryProduct)
	if basePrice == 0.0 {
		panic("init price is zero")
	}
	levelPrice = basePrice * 0.001

	l.Insert(&level{basePrice, ""}, baseLevel)
	for i := 1; i <= perLevel; i++ {
		l.Insert(&level{basePrice * (1.0 + levelPercent*float64(i)), ""}, baseLevel+i)
		l.Insert(&level{basePrice * (1.0 - levelPercent*float64(i)), ""}, baseLevel-i)
	}
	l.sellFront, l.sellRear = baseLevel+safeLevel, baseLevel+safeLevel+length-1
	l.buyFront, l.buyRear = baseLevel-safeLevel, baseLevel-safeLevel-length+1

	log.Println("initPrice:", basePrice)
	fmt.Printf("buyRear:%d buyFront:%d sellFront:%d sellRear:%d\n", l.buyRear, l.buyFront, l.sellFront, l.sellRear)
	for i := l.sellRear; i >= l.sellFront; i-- {
		fmt.Printf("[init] sell_level:%d price:%f\n", i, l.list[i].price)
	}
	for i := l.buyFront; i >= l.buyRear; i-- {
		fmt.Printf("[init] buy_level:%d price:%f\n", i, l.list[i].price)
	}
}

func createMakerOrders(info keys.Info, l *List) ([]types.Order, string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常 when create orders:", err)
		}
	}()

	//1. calculate currentprice
	currentPrice := types.QueryOneTickerPrice(types.QueryProduct)
	if currentPrice == 0.0 {
		fmt.Println("current price is nil")
		return nil, ""
	}
	log.Println("base pirce:", basePrice, "| current price:", currentPrice)
	//printList(l)

	//2. calculate the level of price
	if l.list[baseLevel+1].price <= currentPrice { //价格阶梯往上
		levels := int((currentPrice - basePrice) / levelPrice)
		l.UpLevel(levels) //撤最低价卖单,最低价买单
		baseLevel += levels
	} else if currentPrice <= l.list[baseLevel-1].price { //价格阶梯往下
		levels := int((basePrice - currentPrice) / levelPrice)
		l.DownLevel(levels) //撤最高价卖单,最高价买单
		baseLevel -= levels
	}
	basePrice = l.list[baseLevel].price
	fmt.Println("[afer calculating] baseLevel:", baseLevel, " | basePrice:", basePrice)
	fmt.Printf("[afer calculating] buyRear:%d buyFront:%d sellFront:%d sellRear:%d\n", l.buyRear, l.buyFront, l.sellFront, l.sellRear)

	//3.1 return orders that will be created
	//3.2 return order ids that will be canceled
	return updateOrderList(info, l)
}

func updateOrderList(info keys.Info, l *List) ([]types.Order, string) {
	rand.Seed(time.Now().UnixNano())

	msg := types.QueryOrders(types.Product, info.GetAddress().String())
	if msg == nil {
		fmt.Printf("failed to query the depthbook of product %s from www.okex.me\n", types.Product)
		return nil, ""
	}

	orders := msg.Data.Data
	//create the orders not in the depthbook
	var makerOrderSlice []types.Order
	for i := l.buyFront; i >= l.buyRear; i-- {
		if !checkOrderExist(orders, l.list[i]) {
			requiredQuantity := float64(rand.Intn(9000000)+1000000) * 0.0001
			if tinyMaker {
				requiredQuantity = requiredQuantity / 1000.0
			}
			bidOrder := types.Order{OrderType: types.BUY, Price: l.list[i].price, Quantity: requiredQuantity, Level: i}
			makerOrderSlice = append(makerOrderSlice, bidOrder)
		}
	}
	for i := l.sellFront; i <= l.sellRear; i++ {
		if !checkOrderExist(orders, l.list[i]) {
			requiredQuantity := float64(rand.Intn(9000000)+1000000) * 0.0001
			if tinyMaker {
				requiredQuantity = requiredQuantity / 100.0
			}
			askOrder := types.Order{OrderType: types.SELL, Price: l.list[i].price, Quantity: requiredQuantity, Level: i}
			makerOrderSlice = append(makerOrderSlice, askOrder)
		}
	}

	//get the orders who don't belong to the list anymore
	var ids string
	for i := 0; i < len(orders) && i < 200; i++ {
		if !orders[i].IsFind {
			ids += orders[i].OrderID + ","
		}
	}
	return makerOrderSlice, strings.TrimRight(ids, ",")
}

func checkOrderExist(orders []types.OrderData, l *level) bool {
	if l.id == "" {
		return false
	}

	for i := 0; i < len(orders); i++ {
		if strings.Compare(orders[i].OrderID, l.id) == 0 {
			orders[i].IsFind = true
			return true
		}
	}
	l.id = ""
	return false
}

func setOrderIdInList(l *List, orders []types.Order) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常 when set order id in list:", err)
		}
	}()
	for _, order := range orders {
		l.list[order.Level].id = order.Id
	}
	printList(l)
}

func printList(l *List) {
	for i := l.sellRear; i >= l.sellFront; i-- {
		fmt.Printf("[afer order] sell_level:%d, %+v\n", i, l.list[i])
	}
	for i := l.buyFront; i >= l.buyRear; i-- {
		fmt.Printf("[afer order] buy_level:%d, %+v\n", i, l.list[i])
	}
}
