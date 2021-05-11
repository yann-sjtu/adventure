package token

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
)

const (
	testCoinName = "abc"
	totalSupply  = "100000000.00000000"
	mintable     = true
)

func Issue(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, common.Issue, info)
		return
	}

	_, err = cli.Token().Issue(info, common.PassWord,
		testCoinName, testCoinName, totalSupply, "Used for test "+testCoinName,
		"", mintable, accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, common.Issue, info)
		return
	}
	common.PrintExecuteTxSuccess(common.Issue, info)
}
