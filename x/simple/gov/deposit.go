package gov

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	gosdk "github.com/okex/exchain-go-sdk"
)

const deposit = "deposit"

func Deposit(cli *gosdk.Client, info keys.Info) {
	//accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	//if err != nil {
	//	logger.PrintQueryAccountError(err, deposit, info)
	//	return
	//}
	//cli.Governance().SubmitCommunityPoolSpendProposal()
}
