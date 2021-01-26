package farm_control

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	common2 "github.com/okex/adventure/common"
	"github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	zeroQuoteAmount types.DecCoin

	k = 0
)

func replenishLockedToken(cli *gosdk.Client, requiredToken types.DecCoin) {
	fmt.Printf("======> [Phase2 Replenish] start, require %s \n", requiredToken.String())
	remainToken, totalNewLockedToken, totalNewQuoteToken := requiredToken, zeroLpt, zeroQuoteAmount

	// loop[index:100]
	for r := 0; r < 1; r++ {
		i := (k*1 + r) % 100
		if k%100 == 0 && k != 0 {
			time.Sleep(time.Second * time.Duration(sleepTime))
		}
		index, addr := accounts[i].Index, accounts[i].Address

		// 1. query account
		accInfo, err := cli.Auth().QueryAccount(addr)
		if err != nil {
			log.Printf("[%d] %s failed to query its own account: %s\n", index, addr, err)
			continue
		}

		accNum, seq := accInfo.GetAccountNumber(), accInfo.GetSequence()
		// 2. if there is not enough lpt in this addr, then add-liquidity in swap
		lptToken := types.NewDecCoinFromDec(lockSymbol, accInfo.GetCoins().AmountOf(lockSymbol))
		if lptToken.IsZero() {
			// 3.1 query the account balance
			ownBaseAmount, ownQuoteAmount, err := getOwnBaseCoinAndQuoteCoin(accInfo.GetCoins())
			if err != nil {
				log.Printf("[%d] %s %s\n", index, addr, err.Error())
				continue
			}

			// 3.2 query & calculate how okt could be bought with the number of usdt
			toBaseCoin, toQuoteCoin, err := calculateBaseCoinAndQuoteCoin(cli, ownBaseAmount, ownQuoteAmount)
			if err != nil {
				log.Printf("[%d] %s failed to calculate max-base-coin & quote-coin: %s\n", index, addr, err.Error())
				continue
			}

			// 3.3 add okt & usdt to get lpt
			addLiquidityMsg := newMsgAddLiquidity(accNum, seq, types.ZeroDec(), toBaseCoin, toQuoteCoin, getDeadline(), addr)
			err = common.SendMsg(common.Farm, addLiquidityMsg, index)
			if err != nil {
				log.Printf("[%d] %s failed to add-liquidity: %s\n", index, addr, err)
				continue
			}
			log.Printf("[%d] %s send add-liquidity msg: %+v\n", index, addr, addLiquidityMsg.Msgs[0])
			totalNewQuoteToken = totalNewQuoteToken.Add(toQuoteCoin)
		} else {
			// 3. lock lpt in the farm pool
			lockMsg := newMsgLock(accNum, seq, lptToken, addr)
			err = common.SendMsg(common.Farmlp, lockMsg, index)
			if err != nil {
				log.Printf("[%d] %s failed to lock: %s\n", index, addr, err)
				continue
			}
			log.Printf("[%d] %s send lock msg: %+v\n", index, addr, lockMsg.Msgs[0])

			// 4. update statistics data
			//accounts[i].LockedCoin = accounts[i].LockedCoin.Add(lptToken)
			totalNewLockedToken = totalNewLockedToken.Add(lptToken)
			if remainToken.IsLT(lptToken) {
				remainToken = zeroLpt
				break
			}
			remainToken = remainToken.Sub(lptToken)
		}
	}
	k++
	fmt.Printf("%s is locked in farm, %s is added in swap\n", totalNewLockedToken, totalNewQuoteToken)
	if !remainToken.IsZero() {
		fmt.Printf("%s remainning still have to be replenished\n", remainToken)
	}
}

func getOwnBaseCoinAndQuoteCoin(coins types.DecCoins) (ownBaseAmount, ownQuoteAmount types.DecCoin, err error) {
	ownBaseAmount = types.NewDecCoinFromDec(baseCoin, coins.AmountOf(baseCoin))
	if ownBaseAmount.Amount.LT(types.ZeroDec()) {
		return ownBaseAmount, ownQuoteAmount, fmt.Errorf("has no %s, balance[%s]", baseCoin, coins.String())
	}
	if strings.Compare(ownBaseAmount.Denom, common2.NativeToken) == 0 && ownBaseAmount.Amount.LTE(types.OneDec()) {
		err = fmt.Errorf("has less than 1 %s, so not to add-liquidity. balance[%s]", common2.NativeToken, coins.String())
		return
	}
	ownQuoteAmount = types.NewDecCoinFromDec(quoteCoin, coins.AmountOf(quoteCoin))
	if ownQuoteAmount.Amount.LT(types.ZeroDec()) {
		return ownBaseAmount, ownQuoteAmount, fmt.Errorf("has no %s, balance[%s]", quoteCoin, coins.String())
	}
	return
}

func calculateBaseCoinAndQuoteCoin(cli *gosdk.Client, ownBaseAmount, ownQuoteAmount types.DecCoin) (types.DecCoin, types.DecCoin, error) {
	quotePerPrice, err := cli.AmmSwap().QueryBuyAmount(types.NewDecCoinFromDec(baseCoin, types.OneDec()).String(), quoteCoin)
	if err != nil {
		return types.DecCoin{}, types.DecCoin{}, err
	}
	log.Printf("balance[%s, %s] perPrice:%s \n", ownBaseAmount, ownQuoteAmount, quotePerPrice)
	if ownBaseAmount.Amount.Mul(quotePerPrice).GT(ownQuoteAmount.Amount) {
		// all in quote coin
		toBaseAmount := types.NewDecCoinFromDec(baseCoin, ownQuoteAmount.Amount.Quo(quotePerPrice).Mul(types.MustNewDecFromStr("1.01")))
		return toBaseAmount, ownQuoteAmount, nil
	} else {
		// all in base coin
		if strings.Compare(ownBaseAmount.Denom, common2.NativeToken) == 0 {
			ownBaseAmount.Amount = ownBaseAmount.Amount.Sub(types.OneDec())
		}
		toQuoteAmount := types.NewDecCoinFromDec(quoteCoin, ownBaseAmount.Amount.Mul(quotePerPrice))
		return ownBaseAmount, toQuoteAmount, nil
	}
}

func generateRandomQuoteCoin() types.DecCoin {
	rand.Seed(time.Now().UnixNano())
	// 9000.00~15000.00usdt -> 140~230okt -> 0.8~1.3lpt
	numInt := rand.Intn(600000) + 900000
	numFloat := float64(numInt) / 100.0
	numStr := strconv.FormatFloat(numFloat, 'f', 4, 64)
	return types.NewDecCoinFromDec(quoteCoin, types.MustNewDecFromStr(numStr))
}

func generateRandomMaxBaseCoin() types.DecCoin {
	rand.Seed(time.Now().UnixNano())
	// 240~300okt
	numInt := rand.Intn(6000) + 24000
	numFloat := float64(numInt) / 100.0
	numStr := strconv.FormatFloat(numFloat, 'f', 4, 64)
	return types.NewDecCoinFromDec(baseCoin, types.MustNewDecFromStr(numStr))
}
