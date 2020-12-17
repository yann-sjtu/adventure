package dex

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const list = common.List

func List(cli *gosdk.Client, info keys.Info) {
	time.Sleep(time.Duration(20) * time.Second)

	tokens, err := cli.Token().QueryTokenInfo(info.GetAddress().String(), "")
	if err != nil || len(tokens) == 0 {
		common.PrintQueryTokensError(err, list, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, list, info)
		return
	}

	token := tokens[rand.Intn(len(tokens))]
	_, err = cli.Dex().List(info, common.PassWord,
		token.Symbol, common.NativeToken, "1.0",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, list, info)
		return
	}
	common.PrintExecuteTxSuccess(list, info)
}
