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
)

func replenishLockedToken(cli *gosdk.Client, requiredToken types.DecCoin) {
	fmt.Printf("======> [Phase2 Replenish] start, require %s \n", requiredToken.String())
	remainToken, totalNewLockedToken := requiredToken, zeroLpt

	index := pickRandomIndex()
	// loop[index:100]
	for i := index; i < len(accounts); i++ {
		accInfo, err := cli.Auth().QueryAccount(accounts[i].Address)
		if err != nil {
			log.Printf("[%d] %s failed to query its own account: %s\n", accounts[i].Index, accounts[i].Address, err)
			continue
		}

		accNum, seq := accInfo.GetAccountNumber(), accInfo.GetSequence()
		// if there is not enough lpt in this addr, then add-liquidity in swap
		lptToken := types.NewDecCoinFromDec(lockSymbol, accInfo.GetCoins().AmountOf(lockSymbol))
		if lptToken.IsLT(minLpt) {
			addLiquidityMsg := newMsgAddLiquidity(accNum, seq, minLptDec, defaultMaxBaseAmount, defaultQuoteAmount, getDeadline(), accounts[i].Address)
			err = common.SendMsg(common.Farm, addLiquidityMsg, accounts[i].Index)
			if err != nil {
				log.Printf("[%d] %s failed to add liquidity: %s\n", accounts[i].Index, accounts[i].Address, err)
				continue
			}
			log.Printf("[%d] %s send add-liquidity msg: %+v\n", accounts[i].Index, accounts[i].Address, addLiquidityMsg.Msgs[0])
			lptToken = minLpt
		}

		// 2.4 lock lpt in the farm pool
		lockMsg := newMsgLock(accNum, seq, lptToken, accounts[i].Address)
		err = common.SendMsg(common.Farmlp, lockMsg, accounts[i].Index)
		if err != nil {
			log.Printf("[%d] %s failed to add liquidity: %s\n", accounts[i].Index, accounts[i].Address, err)
			continue
		}
		log.Printf("[%d] %s send lock-pool msg: %+v\n", accounts[i].Index, accounts[i].Address, lockMsg.Msgs[0])

		// 2.5 update accounts
		accounts[i].LockedCoin = accounts[i].LockedCoin.Add(lptToken)

		if remainToken.IsLT(lptToken) {
			break
		}
		remainToken = remainToken.Sub(lptToken)
	}

	//todo: there need another loop[0:index]

	if !remainToken.IsZero() {
		fmt.Printf("%s is still remain, replenish it in next round\n", remainToken)
	}
}
