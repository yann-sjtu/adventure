package dex

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

func Deposit(cli *gosdk.Client, info keys.Info) {
	tokenPairs, err := cli.Dex().QueryProducts(info.GetAddress().String(), 1, 300)
	if err != nil || len(tokenPairs) == 0 {
		common.PrintQueryProductsError(err, common.Deposit, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, common.Deposit, info)
		return
	}

	tokenPair := tokenPairs[rand.Intn(len(tokenPairs))]
	product := tokenPair.BaseAssetSymbol + "_" + tokenPair.QuoteAssetSymbol
	_, err = cli.Dex().Deposit(info, common.PassWord,
		product, "10"+common.NativeToken,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, common.Deposit, info)
		return
	}
	common.PrintExecuteTxSuccess(common.Deposit, info)
}
