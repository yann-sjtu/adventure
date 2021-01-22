package farm_control

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)


//Vote       = 1 //投票
//Undelegate = 2 //赎回
//Staking    = 3 //抵押
//
//Farm       = 4 //添加流动性
//Undelefarm = 5 //删除流动性
//
//Farmlp   = 6 //抵押LP
//Unfarmlp = 7 //删除LP

func TestSendTx(t *testing.T) {
	addr := addrs[0]
	fmt.Println(addr, startIndex)

	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}
	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	// TEST 抵押LP
	msg := newMsgLock(accInfo.GetAccountNumber(), accInfo.GetSequence(), sdk.NewDecCoin(lockSymbol, sdk.NewIntWithDecimal(1,4)), addr)
	err = common.SendMsg(common.Farmlp, msg, startIndex)
	if err != nil {
		fmt.Println("failed:", err)
	}

	// TEST 删除LP
	//msg2 := newMsgUnLock(accInfo.GetAccountNumber(), accInfo.GetSequence(), sdk.NewDecCoin(lockSymbol, sdk.NewIntWithDecimal(1,8)), addr)
	//err = common.SendMsg(common.Unfarmlp, msg2, startIndex)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}


	//duration, err := time.ParseDuration("1m")
	//if err != nil {
	//	return
	//}
	//deadline := time.Now().Add(duration).Unix()
	//// TEST 添加流动性
	//msg3 := newMsgAddLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
	//	sdk.NewDecWithPrec(1,4), sdk.NewDecCoin(baseCoin, sdk.NewIntWithDecimal(1, 1)), sdk.NewDecCoin(quoteCoin, sdk.NewIntWithDecimal(1, 2)), deadline,
	//	addr)
	//err = common.SendMsg(common.Farm, msg3, startIndex)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}
	//
	//// TEST 添加流动性
	//msg4 := newMsgRemoveLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
	//	sdk.NewDecWithPrec(1,4), sdk.NewDecCoin(baseCoin, sdk.NewIntWithDecimal(1, 1)), sdk.NewDecCoin(quoteCoin, sdk.NewIntWithDecimal(1, 2)), deadline,
	//	addr)
	//err = common.SendMsg(common.Undelefarm, msg4, startIndex)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}
}
