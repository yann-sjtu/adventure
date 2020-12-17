package dex

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

func List(cli *gosdk.Client, info keys.Info) {
	time.Sleep(time.Duration(20) * time.Second)

	tokens, err := cli.Token().QueryTokenInfo(info.GetAddress().String(), "")
	if err != nil || len(tokens) == 0 {
		logger.PrintQueryTokensError(err, common.List, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, common.List, info)
		return
	}

	token := tokens[rand.Intn(len(tokens))]
	_, err = cli.Dex().List(info, common.PassWord,
		token.Symbol, common.NativeToken, "1.0",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, common.List, info)
		return
	}
	logger.PrintExecuteTxSuccess(common.List, info)
}
