package ammswap

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const createExchange = common.CreateExchange

func CreateExchange(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, createExchange, info)
		return
	}

	acc, _ := cli.Auth().QueryAccount("okexchain1lgwsujv4efrsf8wsdkz4ggnq0qnnjeqkgwk9yy")
	tokens := acc.GetCoins()
	rand.Seed(time.Now().Unix() + rand.Int63n(100))
	token1 := tokens[rand.Intn(len(tokens))].Denom
	token2 := tokens[rand.Intn(len(tokens))].Denom
	var t1, t2 string
	if token1 < token2 {
		t1, t2 = token1, token2
	} else {
		t1, t2 = token2, token1
	}

	_, err = cli.AmmSwap().CreateExchange(info, common.PassWord,
		t1, t2,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, createExchange, info)
		return
	}
	logger.PrintExecuteTxSuccess(createExchange, info)
}
