package ammswap

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
	swapTypes "github.com/okex/okexchain/x/ammswap/types"
)

const removeLiquidity = common.RemoveLiquidity

func RemoveLiquidity(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, removeLiquidity, info)
		return
	}

	rand.Seed(time.Now().Unix() + rand.Int63n(100))
	tokens := accInfo.GetCoins()
	var token1, token2 types.DecCoin
	for i := 0; i < 5; i++ {
		token1, token2 = tokens[rand.Intn(len(tokens))], tokens[rand.Intn(len(tokens))]
		if strings.Contains(token1.Denom, "swap") == true {
			continue
		}
		if strings.Contains(token2.Denom, "swap") == true {
			continue
		}
	}

	tokenPairName := swapTypes.GetSwapTokenPairName(token1.Denom, token2.Denom)
	swap, err := cli.AmmSwap().QuerySwapTokenPair(tokenPairName)
	if err != nil || swap.PoolTokenName == "" {
		var t1, t2 string
		if token1.Denom < token2.Denom {
			t1, t2 = token1.Denom, token2.Denom
		} else {
			t1, t2 = token2.Denom, token1.Denom
		}

		_, err = cli.AmmSwap().CreateExchange(info, common.PassWord,
			t1, t2,
			"", accInfo.GetAccountNumber(), accInfo.GetSequence())
		if err != nil {
			common.PrintExecuteTxError(err, createExchange, info)
			return
		}
		common.PrintExecuteTxSuccess(createExchange, info)
		return
	}

	_, err = cli.AmmSwap().RemoveLiquidity(info, common.PassWord,
		"0.001", "0"+token1.Denom, "0"+token2.Denom, "10m",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, removeLiquidity+" "+token1.Denom+"_"+token2.Denom, info)
		return
	}
	common.PrintExecuteTxSuccess(removeLiquidity+" "+token1.Denom+"_"+token2.Denom, info)
}
