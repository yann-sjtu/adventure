package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/adventure/common"
	"github.com/okex/okexchain-go-sdk/utils"
	ammswaptypes "github.com/okex/okexchain/x/ammswap/types"
)

//  coinType: 1571
//  hex:
//    msg = "{"chain_id":"okexchain-65","account_number":1,"sequence":1,"fee":{"amount":[{"denom":"okt","amount":"0.020000000000000000"}],"gas":200000},"msgs":[{"sold_token_amount":{"denom":"okt","amount":"1.000000000000000000"},"min_bought_token_amount":{"denom":"usdk-179","amount":"0.001000000000000000"},"deadline":1611313285,"recipient":"okexchain1pt7xrmxul7sx54ml44lvv403r06clrdkgmvr9g","sender":"okexchain1pt7xrmxul7sx54ml44lvv403r06clrdkgmvr9g"}],"memo":""}"
//    accindex = 1
//  relatedId: time
//  txType: 1
func TestStr(t *testing.T) {
	testObjectMarshal()
}

func testObjectMarshal() {
	msg := makeMsg()

	timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	obj, err := NewObject(msg, 1, timeStr, Staking)
	if err != nil {
		panic(err)
	}
	str, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))
}

func testMsgWithIndexMarshal() {
	signMsg := makeMsg()
	msgWithIndex := NewMsgWithIndex(signMsg, 1)
	str, err := json.Marshal(msgWithIndex)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))
}

func testMsgMarshal() {
	signMsg := makeMsg()
	str, err := json.Marshal(signMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(str))
	return
}

func makeMsg() (signMsg authtypes.StdSignMsg) {
	soldTokenDecCoin, err := types.ParseDecCoin("1"+common.NativeToken)
	if err != nil {
		return
	}

	minBoughtTokenDecCoin, err := types.ParseDecCoin("0.001usdk-179")
	if err != nil {
		return
	}

	duration, err := time.ParseDuration("1m")
	if err != nil {
		return
	}
	deadline := time.Now().Add(duration).Unix()
	mnemonic := "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	info, _, _ := utils.CreateAccountWithMnemo(mnemonic, "test", "12345678")

	msg := ammswaptypes.NewMsgTokenToToken(soldTokenDecCoin, minBoughtTokenDecCoin, deadline, info.GetAddress(), info.GetAddress())
	msgs := []types.Msg{msg}

	signMsg = authtypes.StdSignMsg{
		ChainID:       "okexchain-65",
		AccountNumber: 1,
		Sequence:      2,
		Memo:          "hahahah",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, types.NewDecCoinsFromDec(common.NativeToken, types.OneDec())),
	}

	return
}

func makeMsgWithIndex(index int) (msg MsgWithIndex) {
	signMsg := makeMsg()
	return NewMsgWithIndex(signMsg, index)
}