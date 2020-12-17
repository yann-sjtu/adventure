package dex

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const registerOperator = common.RegisterOperator

func RegisterOperator(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, registerOperator, info)
		return
	}

	_, err = cli.Dex().RegisterDexOperator(info, common.PassWord,
		info.GetAddress().String(), "http://xxx/operator.json",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, registerOperator, info)
		return
	}
	common.PrintExecuteTxSuccess(registerOperator, info)
}
