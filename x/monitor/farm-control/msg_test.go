package farm_control

import (
	"fmt"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/okex/okexchain/app/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"

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

var addrs = common.Addrs901To1000

func TestDeposit(t *testing.T) {
	addr := "okexchain1v9asy9x82lk7hfw27kq3pzeg2rgeeg6t5u27uv"

	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}

	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	coins, err := sdk.ParseDecCoin("5.5555okt")
	if err != nil {
		panic(err)
	}

	msg := newMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), coins, addr)
	err = common.SendMsg(common.Staking, msg, 801)
	if err != nil {
		fmt.Println("failed:", err)
	}
}

func TestWithdraw(t *testing.T) {
	addr := "okexchain1v9asy9x82lk7hfw27kq3pzeg2rgeeg6t5u27uv"

	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}

	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	coins, err := sdk.ParseDecCoin("1.1111okt")
	if err != nil {
		panic(err)
	}

	msg := newMsgWithdraw(accInfo.GetAccountNumber(), accInfo.GetSequence(), coins, addr)
	fmt.Printf("%+v \n", msg.Msgs[0])
	err = common.SendMsg(common.Undelegate, msg, 801)
	if err != nil {
		fmt.Println("failed:", err)
	}
}

func TestAddShares(t *testing.T) {
	addr := "okexchain1v9asy9x82lk7hfw27kq3pzeg2rgeeg6t5u27uv"

	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}

	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	valAddrs, err := utils.ParseValAddresses([]string{"okexchainvaloper18au05qx485u2qcw2gvqsrfh29evq77lm45mf4h", "okexchainvaloper1s6nfs7mlj7ewsskkrmekqhpq2w234fczy53wqq"})
	if err != nil {
		panic(err)
	}

	msg := newMsgAddShares(accInfo.GetAccountNumber(), accInfo.GetSequence(), valAddrs, addr)
	err = common.SendMsg(common.Vote, msg, 801)
	if err != nil {
		fmt.Println("failed:", err)
	}
}

func TestSendTx(t *testing.T) {
	addr := addrs[1]
	index := startIndex + 1
	fmt.Println(addr, index)

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
	msg := newMsgLock(accInfo.GetAccountNumber(), accInfo.GetSequence(), sdk.NewDecCoinFromDec(lockSymbol, sdk.MustNewDecFromStr("0.128147322151976840")), addr)
	err = common.SendMsg(common.Farmlp, msg, index)
	if err != nil {
		panic(fmt.Errorf("failed send msg: %s", err.Error()))
	}

	// TEST 删除LP
	//msg2 := newMsgUnLock(accInfo.GetAccountNumber(), accInfo.GetSequence(), sdk.NewDecCoin(lockSymbol, sdk.NewIntWithDecimal(1,8)), addr)
	//err = common.SendMsg(common.Unfarmlp, msg2, index)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}

	//deadline := getDeadline()
	// TEST 添加流动性
	//maxOkt := sdk.NewDecCoinFromDec(baseCoin, sdk.MustNewDecFromStr("4"))
	//usdt := sdk.NewDecCoinFromDec(quoteCoin, sdk.MustNewDecFromStr("200"))
	//
	//msg3 := newMsgAddLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
	//	sdk.MustNewDecFromStr("0.00000001"), maxOkt, usdt, deadline,
	//	addr)
	//err = common.SendMsg(common.Farm, msg3, index)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}
	//
	//// TEST 删除流动性
	//msg4 := newMsgRemoveLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
	//	sdk.NewDecWithPrec(1,4), sdk.NewDecCoin(baseCoin, sdk.NewIntWithDecimal(1, 1)), sdk.NewDecCoin(quoteCoin, sdk.NewIntWithDecimal(1, 2)), deadline,
	//	addr)
	//err = common.SendMsg(common.Undelefarm, msg4, index)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}

	fmt.Println("success")
}

func TestQueryAmount(t *testing.T) {
	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}
	cli := gosdk.NewClient(cfg)
	res, err := cli.AmmSwap().QueryBuyAmount("10"+quoteCoin, baseCoin)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestQueryLockInfo(t *testing.T) {
	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}
	cli := gosdk.NewClient(cfg)
	res, err := cli.Farm().QueryLockInfo(poolName, addrs[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}


func newMsgDeposit(accNum, seqNum uint64, amount sdk.SysCoin, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := stakingtypes.NewMsgDeposit(cosmosAddr, amount)
	msgs := []sdk.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(types.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}

func newMsgWithdraw(accNum, seqNum uint64, amount sdk.SysCoin, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := stakingtypes.NewMsgWithdraw(cosmosAddr, amount)
	msgs := []sdk.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(types.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}

func newMsgAddShares(accNum, seqNum uint64, valAddrs []sdk.ValAddress, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := stakingtypes.NewMsgAddShares(cosmosAddr, valAddrs)
	msgs := []sdk.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(types.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}