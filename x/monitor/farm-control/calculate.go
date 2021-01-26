package farm_control

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	limitRatio  types.Dec
	//lockedRatio = types.NewDecWithPrec(81, 2)
	numerator  types.Dec
	denominator types.Dec

	zeroLpt types.DecCoin
)

func calculateReuiredAmount(cli *gosdk.Client) (types.DecCoin, error) {
	// 1. query how many lpt locked on a farm pool
	totaLockedAmount, err := queryFarmPool(cli)
	if err != nil {
		return zeroLpt, err
	}

	// 2. statistics how many lpt from our accounts locked on a farm pool
	ourTotalLockedAmount := statisticsOurLockedCoinInPool()

	// 3. calculate the ratio ourTotalLockedAmount to totaLockedAmount
	ratio := ourTotalLockedAmount.Quo(totaLockedAmount)
	fmt.Printf("current ratio: %s, limit ratio: %s\n", ratio.String(), limitRatio)
	if ratio.GT(limitRatio) {
		return zeroLpt, nil
	}

	//   ourTotalLockedAmount + requiredAmount
	//   ——————————————————————————————————————  = 0.7
	//   totaLockedAmount     + requiredAmount
	requiredAmount := limitRatio.Mul(totaLockedAmount).Sub(ourTotalLockedAmount).Mul(denominator).Quo(numerator)
	return types.NewDecCoinFromDec(lockSymbol, requiredAmount), nil
}

func statisticsOurLockedCoinInPool() types.Dec {
	totalAmount := zeroLpt
	for i := 0; i < len(accounts); i++ {
		if accounts[i].LockedCoin.IsZero() {
			continue
		}

		totalAmount = totalAmount.Add(accounts[i].LockedCoin)
	}
	fmt.Println("  our total locked:", totalAmount)
	return totalAmount.Amount
}

func queryFarmPool(cli *gosdk.Client) (types.Dec, error) {
	pool, err := cli.Farm().QueryPool(poolName)
	if err != nil {
		return types.ZeroDec(), err
	}
	fmt.Println("whole total locked:", pool.TotalValueLocked)
	return pool.TotalValueLocked.Amount, nil
}
