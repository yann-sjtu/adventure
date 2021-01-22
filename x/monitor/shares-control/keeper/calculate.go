package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gosdk "github.com/okex/okexchain-go-sdk"
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
