package token

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
)

const (
	isDescEdit      = true
	isWholeNameEdit = true
)

func Edit(cli *gosdk.Client, info keys.Info) {
	tokens, err := cli.Token().QueryTokenInfo(info.GetAddress().String(), "")
	if err != nil || len(tokens) == 0 {
		logger.PrintQueryTokensError(err, common.Edit, info)
		return
	}

	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, common.Edit, info)
		return
	}

	symbol := tokens[rand.Intn(len(tokens))].Symbol
	newWholeName := getRandomString(3)
	_, err = cli.Token().Edit(info, common.PassWord,
		symbol, time.Now().String()+" "+newWholeName, newWholeName,
		"", isDescEdit, isWholeNameEdit, accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, common.Edit, info)
		return
	}
	logger.PrintExecuteTxSuccess(common.Edit, info)
}

func getRandomString(length int) string {
	if length == 0 {
		return ""
	}
	clen := len(stdChars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for getRandomString()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = stdChars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
