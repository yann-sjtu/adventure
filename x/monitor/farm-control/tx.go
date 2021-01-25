package farm_control

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	minLptDec = types.MustNewDecFromStr("0.0000000000000")
	minLpt = types.NewDecCoinFromDec(lockSymbol, minLptDec)

	defaultMaxBaseAmount = types.NewDecCoinFromDec(baseCoin, types.MustNewDecFromStr("350"))
	defaultQuoteAmount = types.NewDecCoinFromDec(quoteCoin, types.MustNewDecFromStr("200"))
	zeroQuoteAmount = types.NewDecCoinFromDec(quoteCoin, types.ZeroDec())

	bloom = make([]int, len(accounts), len(accounts))
	k = 0
)


func replenishLockedToken(cli *gosdk.Client, requiredToken types.DecCoin) {
	fmt.Printf("======> [Phase2 Replenish] start, require %s \n", requiredToken.String())
	remainToken, totalNewLockedToken, totalNewQuoteToken := requiredToken, zeroLpt, zeroQuoteAmount

	// loop[index:100]
	for r := 0; r < 1; r++ {
		i := (k*1+r)%100
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
			//toQuoteAmount := generateRandomQuoteCoin()
			// 3.1 query the account balance
			ownQuoteAmount := types.NewDecCoinFromDec(quoteCoin,  accInfo.GetCoins().AmountOf(quoteCoin))
			if ownQuoteAmount.Amount.IsZero() {
				log.Printf("[%d] %s has no %s, balance: %s\n", index, addr, quoteCoin, accInfo.GetCoins().String())
				continue
			}
			//if ownQuoteAmount.Amount.LT(toQuoteAmount.Amount) {
			//	toQuoteAmount = ownQuoteAmount
			//}

			// 3.2 query & calculate how okt could be bought with the number of usdt
			baseCoinPrice, err := cli.AmmSwap().QueryBuyAmount(ownQuoteAmount.String(), baseCoin)
			if err != nil {
				log.Printf("[%d] %s failed to query base coin price: %s\n", index, addr, err.Error())
				continue
			}
			maxBaseAmount := types.NewDecCoinFromDec(baseCoin, baseCoinPrice)
			log.Printf("[%d] %s could add-liquiqity %s with %s \n", index, addr, ownQuoteAmount, maxBaseAmount)
			maxBaseAmount.Amount = maxBaseAmount.Amount.Add(types.MustNewDecFromStr("5"))

			// 3.3 add okt & usdt to get lpt
			addLiquidityMsg := newMsgAddLiquidity(accNum, seq, types.ZeroDec(), maxBaseAmount, ownQuoteAmount, getDeadline(), addr)
			err = common.SendMsg(common.Farm, addLiquidityMsg, index)
			if err != nil {
				log.Printf("[%d] %s failed to add-liquidity: %s\n", index, addr, err)
				continue
			}
			log.Printf("[%d] %s send add-liquidity msg: %+v\n", index, addr, addLiquidityMsg.Msgs[0])
			totalNewQuoteToken = totalNewQuoteToken.Add(ownQuoteAmount)
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

		bloom[i] = 1
	}
	k++

	fmt.Printf("%s is locked in farm, %s is added in swap\n", totalNewLockedToken, totalNewQuoteToken)
	if !remainToken.IsZero() {
		fmt.Printf("%s remainning still have to be replenished\n", remainToken)
	}
}

func generateRandomQuoteCoin() types.DecCoin {
	rand.Seed(time.Now().UnixNano())
	// 9000.00~15000.00usdt -> 140~230okt -> 0.8~1.3lpt
	numInt := rand.Intn(600000) + 900000
	numFloat := float64(numInt)/100.0
	numStr := strconv.FormatFloat(numFloat, 'f', 4, 64)
	return types.NewDecCoinFromDec(quoteCoin, types.MustNewDecFromStr(numStr))
}

func generateRandomMaxBaseCoin() types.DecCoin {
	rand.Seed(time.Now().UnixNano())
	// 240~300okt
	numInt := rand.Intn(6000) + 24000
	numFloat := float64(numInt)/100.0
	numStr := strconv.FormatFloat(numFloat, 'f', 4, 64)
	return types.NewDecCoinFromDec(baseCoin, types.MustNewDecFromStr(numStr))
}