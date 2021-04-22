package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gosdk "github.com/okex/exchain-go-sdk"
)

var (
	dec21     = sdk.NewDec(21)
	percent25 = sdk.MustNewDecFromStr("0.25")
	percent75 = sdk.MustNewDecFromStr("0.75")
)

func (k *Keeper) sumShares(vals []gosdk.Validator) (targetTotal, globalTotal, bondedTotal sdk.Dec) {
	targetTotal, globalTotal, bondedTotal = sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
	for _, val := range vals {
		globalTotal = globalTotal.Add(val.DelegatorShares)
		// check whether target
		if _, ok := k.targetValsFilter[val.OperatorAddress.String()]; ok {
			targetTotal = targetTotal.Add(val.DelegatorShares)
		}
		// check whether boned
		if val.Status.Equal(sdk.Bonded) {
			bondedTotal = bondedTotal.Add(val.DelegatorShares)
		}
	}

	return
}

func (k *Keeper) calculatePercentToPlunder(vals []gosdk.Validator, targetTotal, globalTotal sdk.Dec) sdk.Dec {
	partIn25 := k.expectedParams.GetExpectedValNumberInTop21().Quo(dec21).Mul(percent25)
	partIn75 := targetTotal.Quo(globalTotal).Mul(percent75)
	return partIn25.Add(partIn75)
}
