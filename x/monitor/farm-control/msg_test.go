package farm_control

import (
	"fmt"
	"testing"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/okex/okexchain/app/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
	tokentypes "github.com/okex/okexchain/x/token/types"

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
	index := 801

	cfg, err := gosdk.NewClientConfig("http://10.0.240.38:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}

	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	coins := sdk.NewDecCoinFromDec("okt", sdk.MustNewDecFromStr("100"))

	msg := newMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), coins, addr)
	err = common.SendMsg(common.Staking, msg, index)
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
	//okexchain1v9asy9x82lk7hfw27kq3pzeg2rgeeg6t5u27uv 801
	//okexchain16p9xfn6437rz898a89zq09755ng2fdt5pqp70j 802
	//okexchain1gln5srut8yr4da5czc6rrvwsa8t0nqr0j8py6j 803
	//okexchain1p4fk4kehstehj9k4r0qgqmnpukcxx20c9txvg3 804
	//okexchain1vy0d5a4rh5l42dhs39e8zhjyjapa6ym32g4f6z 805
	addr := "okexchain1p4fk4kehstehj9k4r0qgqmnpukcxx20c9txvg3"
	index := 804

	cfg, err := gosdk.NewClientConfig("http://10.0.240.38:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}

	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	valAddrs, err := utils.ParseValAddresses([]string{
		"okexchainvaloper18au05qx485u2qcw2gvqsrfh29evq77lm45mf4h", //1
		"okexchainvaloper1s6nfs7mlj7ewsskkrmekqhpq2w234fczy53wqq", //2
		"okexchainvaloper1zxthrcdcecfe5ss4tal0tq30hzel2lks2fp8v0", //3
		"okexchainvaloper1q9nct2gska2yutx24starv6s63xz022faxunec", //4
		"okexchainvaloper195ez67wmhprwrru34gvttyd8ttpl7edxpfhu8f", //5
		"okexchainvaloper19wln93k3faq7vkqzlc9gljr3ey5fljt9p6cats", //6
		"okexchainvaloper1qva0ejf0t943x6rt824gwmvtjgec9cjrvr94gn", //7
		"okexchainvaloper1ka92ujcwh6hyyeu4tymzy3dedgxplt4dmcj9ar", //8
		"okexchainvaloper1zza3jrylyecrtuh0p9ts2xauzsefuvwa9h5jtj", //9
		"okexchainvaloper1ja9xngm4zh0t442mse73ll30p7dczd49q0kg3j", //10
		"okexchainvaloper1c34s7lc7ec8gs9xrtxeh0j2wjaam25c3c8ta69", //11
		"okexchainvaloper1ygcvtcqxl82xvzrq25dymam434k3nnc8xxacd0", //12
		"okexchainvaloper1m569cfenudxemegcf4mmykhugnslhdv0klarfw", //13
		"okexchainvaloper1fymxn4gazxzjdfvwvr0ccnrnjpwmj0r9uxqs3s", //14
		"okexchainvaloper1xkl5agjzqnjnptyat2dng2asmx8g5kllckhxqc", //15
		"okexchainvaloper1tkwxgcpvptua0q0h5tn0at58ufnjdue7kf5fvp", //16
		"okexchainvaloper1508d7eq592kg2lh9d46xvv3r4sm7gm8wlmjzfz", //17
		"okexchainvaloper18v23ln9ycrtg0mrwsm004sh4tdknudtddffjr5", //18
		"okexchainvaloper1ucmx6vvtrwam9pg20fnwmy9z80uhchyxsmt945", //19
		"okexchainvaloper1g3a6vtau2k93n4tqgqnrggeu3qa4x20rccyawy", //20
		"okexchainvaloper19e6edpu97d6w2t5dlp7lph2fkdja0lvlz0zndm", //21

		//"okexchainvaloper1tat4lam8wjqmeax9mv4s584vu2mp7c0ccywxft", //22
		//"okexchainvaloper1mlmwvdprn8dj6g45vdxkjsjgu4ntu9j7amrdl7", //23
		//"okexchainvaloper1rz7frqz9ky52qqjwlpawfe5hz6plcrmm3lv56j", //24
		//"okexchainvaloper1w3ptfgekjgdvwkqmdepdeyvuxqmcplfsjhn2f0", //25
		//"okexchainvaloper1v4kagglr3vq82vqywqd8quhsuarkm4kf6mnu0h", //26
		//"okexchainvaloper1rmrx7wp60almzvghx2820aamjfd4kgwlgw9w34", //27
		//"okexchainvaloper13mayrjzsrp976y0ae0qw8sjan3qg2xfdfgkhqr", //28
		//"okexchainvaloper104y8sy0r6fke4a9qr8u05j6v5y68gkh4uedk7l", //29
		//"okexchainvaloper14zgafe7cynlpuhpfpqpxu2gntzhq6tteagj8px", //30
		//"okexchainvaloper1rv8tjxp8d8ucuak8c7svewwugzfdjwf9dtr80x", //31
	})
	if err != nil {
		panic(err)
	}

	msg := newMsgAddShares(accInfo.GetAccountNumber(), accInfo.GetSequence(), valAddrs, addr)
	err = common.SendMsg(common.Vote, msg, index)
	if err != nil {
		fmt.Println("failed:", err)
	}
}

func TestSendTx(t *testing.T) {
	addr := "okexchain1hw5h3genzr4ftc2q7zepyxk84qmz54lph5mqny"
	index := 706
	fmt.Println(addr, index)
	poolName := "okb_usdt"
	ammswapTokenName := "ammswap_okb-c4d_usdt-a2b"

	cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	if err != nil {
		panic(err)
	}
	cli := gosdk.NewClient(cfg)
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	// TEST 添加流动性
	baseCoin = "okb-c4d"
	quoteCoin = "usdt-a2b"
	ownBaseAmount, ownQuoteAmount, err := getOwnBaseCoinAndQuoteCoin(accInfo.GetCoins())
	if err != nil {
		panic(fmt.Errorf("[%d] %s %s\n", index, addr, err.Error()))
	}

	// 3.2 query & calculate how okt could be bought with the number of usdt
	toBaseCoin, toQuoteCoin, err := calculateBaseCoinAndQuoteCoin(&cli, ownBaseAmount, ownQuoteAmount)
	if err != nil {
		panic(fmt.Errorf("[%d] %s failed to calculate max-base-coin & quote-coin: %s\n", index, addr, err.Error()))
	}
	fmt.Println(toBaseCoin, toQuoteCoin)

	msg3 := newMsgAddLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
		sdk.MustNewDecFromStr("0.00000001"), toBaseCoin, toQuoteCoin, getDeadline(),
		addr)
	err = common.SendMsg(common.Farm, msg3, index)
	if err != nil {
		fmt.Println("failed:", err)
	}
	//
	//// TEST 删除流动性
	//msg4 := newMsgRemoveLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
	//	sdk.NewDecWithPrec(1,4), sdk.NewDecCoin(baseCoin, sdk.NewIntWithDecimal(1, 1)), sdk.NewDecCoin(quoteCoin, sdk.NewIntWithDecimal(1, 2)), deadline,
	//	addr)
	//err = common.SendMsg(common.Undelefarm, msg4, index)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}

	time.Sleep(time.Minute)
	accInfo, err = cli.Auth().QueryAccount(addr)
	if err != nil {
		panic(err)
	}
	// TEST 抵押LP
	msg := newMsgLock(accInfo.GetAccountNumber(), accInfo.GetSequence(),
		poolName,
		sdk.NewDecCoinFromDec(ammswapTokenName, accInfo.GetCoins().AmountOf(ammswapTokenName)),
		addr)
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
		Fee:           authtypes.NewStdFee(5000000, sdk.NewDecCoinsFromDec(types.NativeToken, sdk.MustNewDecFromStr("0.005"))),
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
		Fee:           authtypes.NewStdFee(1100000, sdk.NewDecCoinsFromDec(types.NativeToken, sdk.MustNewDecFromStr("0.011"))),
	}

	return signMsg
}

func TestTokenTransfer(t *testing.T) {
	var addrs = []string{
		"okexchain140sy6yc4mx7equhrzavm20snzm5vd9va6f5rfw",
		"okexchain1sey3syrw0xsqsvkdlk86j72w0hfnq8j4yqk9e3",
		"okexchain1mpmw2jx0dfnrh73tj0xzcc79n5ddsp5s323y03",
		"okexchain183quv4zattal80tkq3ccnxntv3yrpt6yyjt6sa",
		"okexchain1zyufn2zm8az8664ed07muskktz2qmkuymf3ye3",
		"okexchain120rxcq6w9y46qa2ewyzuxvx0t74y2ch506kl2h",
		"okexchain1l94x89s6d65ffzzt2ns8lyr9d958u7dkm25zc0",
		"okexchain1tpg9mmvqw4j97z0kgwuz4xum7swva9s8j5qzlc",
		"okexchain1g82hlllygaf6rnnsaxqdl0xxmue2fwt2j9hdkf",
		"okexchain189sq8hphj3kzp8a302kk48r7m4f2kq4z2vu0u7",
		"okexchain1h4t9z7amss2tmy07efngjez3zrpe7zrg4k95kp",
		"okexchain17rkgqreruk9wchyf4a62n32g82sngnp6sjc0dz",
		"okexchain15v8k8gfp2paxrpaw98mnf9pfycgr4xard3u8yr",
		"okexchain1ah6fu38g6nm9rmksa7uc6hn4qyu4nah8335900",
		"okexchain146dh4fw7a9qycqhagd7zwkj2n833n0tx8gtwy9",
		"okexchain19z2jzft3y8dlkeaxpnraccrdfxn0uz079kwfvy",
		"okexchain157p3dta442g9cav3l0g5ws4rr79al3rpvls0ju",
		"okexchain1hg5synr7qxqsyc0gj2r0hvtdf0kntfsl73xp23",
		"okexchain1cv2vv36kk8adk2rve0766lwr6q50qsg7se2x03",
		"okexchain1yw6qx8dudxpkeghdh7n8z300e4svxtzrk2qc6j",
		"okexchain1rmpk5rmsyagakdxx7t8xny8eglu6lp5dvj8g4w",
	}
	for _,addr := range addrs {
		cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
		if err != nil {
			panic(err)
		}
		cli := gosdk.NewClient(cfg)
		accInfo, err := cli.Auth().QueryAccount(addr)
		if err != nil {
			panic(err)
		}
		fmt.Println(accInfo.GetCoins())
	}


	//index := 907
	//addr := "okexchain1g82hlllygaf6rnnsaxqdl0xxmue2fwt2j9hdkf"
	//to := "okexchain1d66fc3dsddvzyl2zhkd4002lw5jypcxx9u6fmr"
	//
	//cfg, err := gosdk.NewClientConfig("http://10.0.240.37:26657", "okexchain-66", gosdk.BroadcastBlock, "0.002okt", 200000, 0, "")
	//if err != nil {
	//	panic(err)
	//}
	//
	//cli := gosdk.NewClient(cfg)
	//accInfo, err := cli.Auth().QueryAccount(addr)
	//if err != nil {
	//	panic(err)
	//}
	//
	//amount := sdk.NewDecCoinsFromDec("usdt-a2b", accInfo.GetCoins().AmountOf("usdt-a2b"))
	//msg := newSendToken(accInfo.GetAccountNumber(), accInfo.GetSequence(),addr,to, amount)
	//fmt.Printf("%+v \n", msg.Msgs[0])
	//err = common.SendMsg(common.Send, msg, index)
	//if err != nil {
	//	fmt.Println("failed:", err)
	//}
}

func newSendToken(accNum, seqNum uint64, from, to string, amount sdk.DecCoins) authtypes.StdSignMsg {
	fromCosmosAddr, err := utils.ToCosmosAddress(from)
	if err != nil {
		panic(err)
	}
	toCosmosAddr, err := utils.ToCosmosAddress(to)
	if err != nil {
		panic(err)
	}

	msg := tokentypes.MsgSend{fromCosmosAddr, toCosmosAddr, amount}
	msgs := []sdk.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(500000, sdk.NewDecCoinsFromDec(types.NativeToken, sdk.MustNewDecFromStr("0.005"))),
	}

	return signMsg
}