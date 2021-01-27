package utils

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// UTC Time: 2000/1/1 00:00:00
	blockTimestampEpoch = int64(946684800)
	secondsPerWeek      = int64(60 * 60 * 24 * 7)
	weeksPerYear        = float64(52)
)

// copy from okexchain, don't use
func calculateWeight(nowTime int64, tokens sdk.Dec) (shares sdk.Dec, sdkErr error) {
	nowWeek := (nowTime - blockTimestampEpoch) / secondsPerWeek
	rate := float64(nowWeek) / weeksPerYear
	weight := math.Pow(float64(2), rate)

	precision := fmt.Sprintf("%d", sdk.Precision)

	weightByDec, sdkErr := sdk.NewDecFromStr(fmt.Sprintf("%." + precision + "f", weight))
	if sdkErr == nil {
		shares = tokens.Mul(weightByDec)
	}
	return
}

// use this to reverse how much token to be deposited with the specific number of share
func ReverseSharesIntoToken(shares sdk.Dec, nowTime int64) sdk.Dec {
	nowWeek := (nowTime - blockTimestampEpoch) / secondsPerWeek
	rate := float64(nowWeek) / weeksPerYear
	weight := math.Pow(float64(2), rate)
	precision := fmt.Sprintf("%d", sdk.Precision)

	weightByDec, sdkErr := sdk.NewDecFromStr(fmt.Sprintf("%." + precision + "f", weight))
	if sdkErr == nil {
		token := shares.Quo(weightByDec)
		return token
	}
	return sdk.ZeroDec()
}