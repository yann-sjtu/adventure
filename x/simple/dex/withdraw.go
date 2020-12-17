package dex

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

func Withdraw(cli *gosdk.Client, info keys.Info) {
	tokenPairs, err := cli.Dex().QueryProducts(info.GetAddress().String(), 1, 300)
	if err != nil || len(tokenPairs) == 0 {
		common.PrintQueryProductsError(err, common.Withdraw, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, common.Withdraw, info)
		return
	}

	tokenPair := tokenPairs[rand.Intn(len(tokenPairs))]
	product := tokenPair.BaseAssetSymbol + "_" + tokenPair.QuoteAssetSymbol
	_, err = cli.Dex().Withdraw(info, common.PassWord,
		product, "1.0"+common.NativeToken,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, common.Withdraw, info)
		return
	}
	common.PrintExecuteTxSuccess(common.Withdraw, info)
}
