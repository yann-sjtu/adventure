package ammswap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/okex/adventure/common"
	"github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
)

const password = "12345678"
const coin = "usdk-739"
const host = "http://127.0.0.1:26657"

func Test_Create(t *testing.T) {
	mnemonic := "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	info, _, _ := utils.CreateAccountWithMnemo(mnemonic, "test", password)

	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := cli.AmmSwap().CreateExchange(info, password, "usdk-179", common.NativeToken, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func Test_AddLiquidity(t *testing.T) {
	mnemonic := "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	info, _, _ := utils.CreateAccountWithMnemo(mnemonic, "test", password)

	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		fmt.Println(err)
		return
	}

	//res, err := cli.AmmSwap().AddLiquidity(info, common.PassWord,
	//	"1", "1usdk-739", "100000okt", "10m",
	//	"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(res)
	res, err := cli.AmmSwap().AddLiquidity(info, common.PassWord,
		"0.000001", "100"+common.NativeToken, "0.001usdk-179", "10m",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	fmt.Println(res)
}

func Test_DeleteLiquidity(t *testing.T) {
	mnemonic := "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	info, _, _ := utils.CreateAccountWithMnemo(mnemonic, "test", password)

	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := cli.AmmSwap().RemoveLiquidity(info, common.PassWord,
		"0.000001", "0.0000"+common.NativeToken, "0usdk-179", "10m",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func Test_Swap(t *testing.T) {
	mnemonic := "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	info, _, _ := utils.CreateAccountWithMnemo(mnemonic, "test", password)

	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := cli.AmmSwap().TokenSwap(info, common.PassWord,
		"1"+common.NativeToken, "0.001usdk-179", info.GetAddress().String(), "1m",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func Test_QuerySwap(t *testing.T) {
	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	res, err := cli.AmmSwap().QuerySwapTokenPair("btc-a12_xxb-060")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

	res1, err1 := cli.AmmSwap().QuerySwapTokenPairs()
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println(res1)
}

func Test_QueryBuyAmount(t *testing.T) {
	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	res, err := cli.AmmSwap().QueryBuyAmount("btc-a12", "1.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func Test_QueryToken(t *testing.T) {
	host := "http://10.0.240.23:20057"
	// pick a client randomly
	cfg, _ := types.NewClientConfig(host, "okexchain", types.BroadcastBlock, "", 200000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	//okexchain1lgwsujv4efrsf8wsdkz4ggnq0qnnjeqkgwk9yy
	acc, _ := cli.Auth().QueryAccount("okexchain1lgwsujv4efrsf8wsdkz4ggnq0qnnjeqkgwk9yy")
	fmt.Println(len(acc.GetCoins()))

	pairs, _ := cli.AmmSwap().QuerySwapTokenPairs()
	fmt.Println(len(pairs))
	fmt.Println(pairs[len(pairs)-1])

	//res, err := cli.AmmSwap().QuerySwapTokenPair(coin + "_okt")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
}

func Test_Json(t *testing.T) {
	type AutoGenerated struct {
		AppState struct {
			Ammswap struct {
				Params struct {
					FeeRate string `json:"fee_rate"`
				} `json:"params"`
				SwapTokenPairRecords []struct {
					BasePooledCoin struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"base_pooled_coin"`
					PoolTokenName   string `json:"pool_token_name"`
					QuotePooledCoin struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"quote_pooled_coin"`
				} `json:"swap_token_pair_records"`
			} `json:"ammswap"`
		} `json:"app_state"`
	}

	data, err := ioutil.ReadFile("/Users/green/project/okex/adventure/genesis.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var a AutoGenerated

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, &a)
	if err != nil {
		fmt.Println(err)
		return
	}
	k := len(a.AppState.Ammswap.SwapTokenPairRecords)
	fmt.Println(k)
	fmt.Println(a.AppState.Ammswap.SwapTokenPairRecords[0])
	fmt.Println(a.AppState.Ammswap.SwapTokenPairRecords[k-1])
}

//okexchain1zux0wa4xzc67chnf089gmx67nxmqdw9eepj8jk
//okexchain12upy8deq4645wlvva7zq604jt63wfl28vsstru
//okexchain1kpnsm0w9c65zgwhdzamhacp8fryxangs0heqke
//okexchain1qlsqszsghed2v7u4069hmsvfj69cxy57kcu6uv
//okexchain1asnrfkwmafpxqh8k3s6xvu0a6f2yk3nn4vd25g
//okexchain1tfs6xp9aceypkjd9gqfhrx8wjk0r8s32q4smjq
//okexchain1vt6y9yuwm5ka4hxwlzr06j38gjeyentupul8dz
//okexchain1hshp4zhq8q2qfcqmeqkk2gll3wms7af6qctvzm
//okexchain1pst8ck3w20qh0y9qw6j8mne4gf7z4v4pfmldlx
//okexchain1t4d6z5qev0mfdj07wjfp7ckshd2suynz46pwl6
