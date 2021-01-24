package farm_control

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	minLptDec = types.MustNewDecFromStr("0.001")
	minLpt = types.NewDecCoinFromDec(lockSymbol, minLptDec)

	defaultMaxBaseAmount = types.NewDecCoinFromDec(baseCoin, types.MustNewDecFromStr("0.25"))
	defaultQuoteAmount = types.NewDecCoinFromDec(quoteCoin, types.MustNewDecFromStr("13"))
	zeroQuoteAmount = types.NewDecCoinFromDec(quoteCoin, types.ZeroDec())
)

func replenishLockedToken(cli *gosdk.Client, requiredToken types.DecCoin) {
	fmt.Printf("======> [Phase2 Replenish] start, require %s \n", requiredToken.String())
	remainToken, totalNewLockedToken, totalNewQuoteToken := requiredToken, zeroLpt, zeroQuoteAmount

	// loop[index:100]
	bloom := make([]int, len(accounts), len(accounts))
	for i := 0; i < len(accounts)/5; i++ {
		r := pickRandomIndex()
		if bloom[r] == 1 {
			continue
		}
		bloom[r] = 1
		index, addr := accounts[r].Index, accounts[r].Address
		
		// 1. query account
		accInfo, err := cli.Auth().QueryAccount(addr)
		if err != nil {
			log.Printf("[%d] %s failed to query its own account: %s\n", index, addr, err)
			continue
		}

		accNum, seq := accInfo.GetAccountNumber(), accInfo.GetSequence()
		// 2. if there is not enough lpt in this addr, then add-liquidity in swap
		lptToken := types.NewDecCoinFromDec(lockSymbol, accInfo.GetCoins().AmountOf(lockSymbol))
		if lptToken.IsLT(minLpt) {
			addLiquidityMsg := newMsgAddLiquidity(accNum, seq, minLptDec, defaultMaxBaseAmount, defaultQuoteAmount, getDeadline(), addr)
			err = common.SendMsg(common.Farm, addLiquidityMsg, index)
			if err != nil {
				log.Printf("[%d] %s failed to add-liquidity: %s\n", index, addr, err)
				continue
			}
			log.Printf("[%d] %s send add-liquidity msg: %+v\n", index, addr, addLiquidityMsg.Msgs[0])
			totalNewQuoteToken = totalNewQuoteToken.Add(defaultQuoteAmount)
			lptToken = minLpt
		}

		// 3. lock lpt in the farm pool
		lockMsg := newMsgLock(accNum, seq+1, lptToken, addr)
		err = common.SendMsg(common.Farmlp, lockMsg, index)
		if err != nil {
			log.Printf("[%d] %s failed to lock: %s\n", index, addr, err)
			continue
		}
		log.Printf("[%d] %s send lock msg: %+v\n", index, addr, lockMsg.Msgs[0])

		// 4. update statistics data
		accounts[r].LockedCoin = accounts[r].LockedCoin.Add(lptToken)
		totalNewLockedToken = totalNewLockedToken.Add(lptToken)
		if remainToken.IsLT(lptToken) {
			remainToken = zeroLpt
			break
		}
		remainToken = remainToken.Sub(lptToken)
	}

	// todo: there need another loop[0:index]

	fmt.Printf("%s is locked in farm, %s is added in swap\n", totalNewLockedToken, totalNewQuoteToken)
	if !remainToken.IsZero() {
		fmt.Printf("%s remainning still have to be replenished\n", remainToken)
	}
}
