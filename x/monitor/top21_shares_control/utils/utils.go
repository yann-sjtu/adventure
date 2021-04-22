package utils

import (
	"fmt"
	stakingtypes "github.com/okex/exchain/x/staking/types"
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

	weightByDec, sdkErr := sdk.NewDecFromStr(fmt.Sprintf("%."+precision+"f", weight))
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

	weightByDec, sdkErr := sdk.NewDecFromStr(fmt.Sprintf("%."+precision+"f", weight))
	if sdkErr == nil {
		token := shares.Quo(weightByDec)
		return token
	}
	return sdk.ZeroDec()
}

func BuildFilter(addrs []string) map[string]struct{} {
	filter := make(map[string]struct{})

	for _, addr := range addrs {
		filter[addr] = struct{}{}
	}

	return filter
}

func ConvertValAddrsStr2AccAddrsStr(valAddrsStr []string) (accAddrsStr []string, err error) {
	for _, valAddrStr := range valAddrsStr {
		valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			return nil, err
		}

		accAddrsStr = append(accAddrsStr, sdk.AccAddress(valAddr).String())
	}

	return
}

func GetValAddrsStrFromVals(vals []stakingtypes.Validator) []string {
	var valAddrsStr []string
	for _, val := range vals {
		valAddrsStr = append(valAddrsStr, val.OperatorAddress.String())
	}

	return valAddrsStr
}

func RemoveDuplicate(list []string) []string {
	filter := make(map[string]struct{})
	for _, item := range list {
		filter[item] = struct{}{}
	}

	var res []string
	for key := range filter {
		res = append(res, key)
	}

	return res
}
