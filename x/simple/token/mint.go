package token

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const (
	coinForMint = "100.0"
)

func Mint(cli *gosdk.Client, info keys.Info) {
	tokens, err := cli.Token().QueryTokenInfo(info.GetAddress().String(), "")
	if err != nil || len(tokens) == 0 {
		logger.PrintQueryTokensError(err, common.Mint, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, common.Mint, info)
		return
	}

	symbol := tokens[rand.Intn(len(tokens))].Symbol
	_, err = cli.Token().Mint(info, common.PassWord,
		coinForMint+symbol,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, common.Mint, info)
		return
	}
	logger.PrintExecuteTxSuccess(common.Mint, info)
}
