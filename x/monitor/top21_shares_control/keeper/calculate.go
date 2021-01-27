package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) SumShares() (enemyTotalShares, tarValsTotalShares sdk.Dec, err error) {
	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
	if err != nil {
		return
	}
	enemyTotalShares, tarValsTotalShares = sdk.ZeroDec(), sdk.ZeroDec()

	var counter int
	for _, val := range vals {
		if _, ok := k.targetValsFilter[sdk.AccAddress(val.OperatorAddress).String()]; ok {
			tarValsTotalShares = tarValsTotalShares.Add(val.DelegatorShares)
			counter++
		} else {
			enemyTotalShares = enemyTotalShares.Add(val.DelegatorShares)
		}
	}

	fmt.Printf("%d target vals [%s], %d enemy vals [%s]\n",
		counter, tarValsTotalShares.String(), 21-counter, enemyTotalShares.String())

	return
}

func (k *Keeper) CalculateHowMuchToDeposit(enemyTotalShares, tarValsCurrentTotalShares sdk.Dec) error {
	expectedTotalShares := enemyTotalShares.Quo(sdk.OneDec().Sub(k.dominationPct))
	//expectedTarValsTotalShares := expectedTotalShares.Mul(k.dominationPct)
	_ = expectedTotalShares.Mul(k.dominationPct)
	// pick worker randomly
	return nil

}
