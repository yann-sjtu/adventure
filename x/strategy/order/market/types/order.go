package types

import (
	"fmt"
	"log"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/utils"
)

type Order struct {
	OrderType string
	Price     float64
	Quantity  float64
	Level     int
	Id        string
}

func PlaceOrders(cli *gosdk.Client, info keys.Info, orders *[]Order) {
	if len(*orders) == 0 {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("捕获异常 when place orders:", err)
		}
	}()

	var products, orderTypes, prices, quantities string
	for _, order := range *orders {
		fmt.Printf("create order %+v\n", order)
		products += Product + ","
		orderTypes += order.OrderType + ","
		prices += fmt.Sprintf("%.4f", order.Price) + ","
		quantities += fmt.Sprintf("%.4f", order.Quantity) + ","
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	res, err := cli.Order().NewOrders(info, PassWd,
		strings.TrimRight(products, ","), strings.TrimRight(orderTypes, ","),
		strings.TrimRight(prices, ","), strings.TrimRight(quantities, ","),
		"new orders", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		panic(err)
	}

	ids, err := utils.GetOrderIDsFromResponse(&res)
	if err != nil {
		panic(err)
	}
	for i, id := range ids {
		(*orders)[i].Id = id
		log.Printf("%s successfully submit new order %+v\n", info.GetAddress().String(), (*orders)[i])
	}
}

func CancelOrders(cli *gosdk.Client, info keys.Info, ids string) {
	if ids == "" {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("failed to cancel orders", ids, ". error:", err)
		}
	}()

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	_, err = cli.Order().CancelOrders(info, PassWd, ids,
		"cancel orders", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		fmt.Println("failed to cancel orders", ids, ":", err)
		return
	}
	fmt.Println("successfully cancel orders:", ids)
}
