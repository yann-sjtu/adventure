package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

	Send = 8 //转账
)

var (
	key    = "273186b97e6741a5a8fe68383d4949c8"
	secret = "091b30bfb6604917945752de9cb87609"

	serverUrl = "https://www.okex.com"
	ctx       = "/vault/api/v2/okpool/voteOKT"
)

func SendMsg(txType int, msg types.StdSignMsg, addresIndex int) error {
	// 0.1 time now
	timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 0.2
	object, err := NewObject(msg, addresIndex, timeStr, txType)
	if err != nil {
		return err
	}
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
	//sss := string(reqStr)
	//fmt.Println(sss)
	nonce := timeStr
	// 0 sign with hmac-sha256
	signature := hmacSha256(secret, nonce+ctx+string(reqStr))

	// 1.1 init new request
	req, err := http.NewRequest(http.MethodPost, serverUrl+ctx, bytes.NewBuffer(reqStr))
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

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "success") {
		return errors.New(bodyStr)
	}

	return nil
}
