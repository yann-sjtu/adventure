package dex

import (
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
)

const editOperator = common.EditOperator

func EditOperator(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, editOperator, info)
		return
	}

	_, err = cli.Dex().EditDexOperator(info, common.PassWord,
		info.GetAddress().String(), "http://"+strconv.Itoa(rand.Int())+"/operator.json",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, editOperator, info)
		return
	}
	common.PrintExecuteTxSuccess(editOperator, info)
}
