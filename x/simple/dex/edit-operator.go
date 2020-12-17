package dex

import (
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const editOperator = common.EditOperator

func EditOperator(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, editOperator, info)
		return
	}

	_, err = cli.Dex().EditDexOperator(info, common.PassWord,
		info.GetAddress().String(), "http://"+strconv.Itoa(rand.Int())+"/operator.json",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, editOperator, info)
		return
	}
	logger.PrintExecuteTxSuccess(editOperator, info)
}
