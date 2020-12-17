package order

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const order = common.Order

func Orders(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, order, info)
		return
	}

	cfg := config.Cfg.Order.NewConfig
	product := cfg.Products[rand.Intn(len(cfg.Products))]
	_, err = cli.Order().NewOrders(info, common.PassWord,
		product+","+product, "BUY,SELL", cfg.BuyPrice+","+cfg.SellPrice, cfg.BuyQuantity+","+cfg.SellQuantity,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, order, info)
		return
	}
	logger.PrintExecuteTxSuccess(order, info)
}
