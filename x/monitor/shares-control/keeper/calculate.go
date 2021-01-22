package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gosdk "github.com/okex/okexchain-go-sdk"
)

func (k *Keeper) sumShares(vals []gosdk.Validator) (targetTotal, globalTotal sdk.Dec) {
	targetTotal, globalTotal = sdk.ZeroDec(), sdk.ZeroDec()
	for _, val := range vals {
		globalTotal = globalTotal.Add(val.DelegatorShares)
		// check whether target
		if _, ok := k.targetValsFilter[val.OperatorAddress.String()]; ok {
			targetTotal = targetTotal.Add(val.DelegatorShares)
		}
	}

	return
}
