package token

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const (
	testCoinName = "abc"
	totalSupply  = "100000000.00000000"
	mintable     = true
)

func Issue(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, common.Issue, info)
		return
	}

	_, err = cli.Token().Issue(info, common.PassWord,
		testCoinName, testCoinName, totalSupply, "Used for test "+testCoinName,
		"", mintable, accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, common.Issue, info)
		return
	}
	logger.PrintExecuteTxSuccess(common.Issue, info)
}
