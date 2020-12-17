package token

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const (
	coinForBurn = "100.0"
)

func Burn(cli *gosdk.Client, info keys.Info) {
	tokens, err := cli.Token().QueryTokenInfo(info.GetAddress().String(), "")
	if err != nil || len(tokens) == 0 {
		common.PrintQueryTokensError(err, common.Burn, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, common.Burn, info)
		return
	}

	symbol := tokens[rand.Intn(len(tokens))].Symbol
	_, err = cli.Token().Burn(info, common.PassWord,
		coinForBurn+symbol,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, common.Burn, info)
		return
	}
	common.PrintExecuteTxSuccess(common.Burn, info)
}
