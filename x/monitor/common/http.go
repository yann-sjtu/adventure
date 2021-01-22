package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	Vote       = 1 //投票
	Undelegate = 2 //赎回
	Staking    = 3 //抵押

	Farm       = 4 //添加流动性
	Undelefarm = 5 //删除流动性

	Farmlp   = 6 //抵押LP
	Unfarmlp = 7 //删除LP
)

var (
	secret = ""

	serverUrl = "https://www.okex.com"
	ctx       = "/vault/api/v2/okpool/vote"

	key       = ""
)

func SendMsg(txType int, addresIndex int, msg types.StdSignMsg) error {
	// 0.1 time now
	timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 0.2
	msgWithIndex := NewMsgWithIndex(msg, addresIndex)
	msgStr, err := json.Marshal(msgWithIndex)
	if err != nil {
		return err
	}

	object := NewObject(string(msgStr), timeStr, txType)
	err = DoPost(timeStr, object)
	if err != nil {
		return err
	}
	return nil
}

func DoPost(timeStr string, object Object) error {
	reqStr, err := json.Marshal(object)
	if err != nil {
		return err
	}
	nonce := timeStr
	// 0.3 sign with hmac-sha256
	signature := hmacSha256(secret, nonce+ctx+string(reqStr))

	// 1.1 init new request
	req, err := http.NewRequest(http.MethodPost, serverUrl, nil)
	if err != nil {
		panic(err)
	}
	// 1.2 set params in header
	req.Header.Set("KEY", key)
	req.Header.Set("SIGNATURE", signature)
	req.Header.Set("NONCE", nonce)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("resp failed: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read failed: %s", err.Error())
	}

	fmt.Println(string(body))
	return nil
}
