package farm_control

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/exchain-go-sdk"
)

var (
	limitRatio = types.MustNewDecFromStr("0.70")
	numerator = types.MustNewDecFromStr("3.0")
	denominator = types.MustNewDecFromStr("10.0")
)

func calculateReuiredAmount(cli *gosdk.Client, accs []common.Account) (types.DecCoin, error) {
	// 1. query how many lpt locked on a farm pool
	totaLockedAmount, err := queryFarmPool(cli)
	if err != nil {
		return types.DecCoin{}, err
	}

	// 2. statistics how many lpt from our accounts locked on a farm pool
	ourTotalLockedAmount, err := statisticsOurLockedCoinInPool(cli, accs)
	if err != nil {
		return types.DecCoin{}, err
	}

	// 3. calculate the ratio ourTotalLockedAmount to totaLockedAmount
	ratio := ourTotalLockedAmount.Quo(totaLockedAmount)
	fmt.Printf("current ratio: %s, limit ratio: %s\n", ratio.String(), limitRatio)
	if ratio.GT(limitRatio) {
		return types.NewDecCoinFromDec(lockSymbol, types.ZeroDec()), nil
	}

	//   ourTotalLockedAmount + requiredAmount
	//   ——————————————————————————————————————  = 0.7
	//   totaLockedAmount     + requiredAmount
	requiredAmount := limitRatio.Mul(totaLockedAmount).Sub(ourTotalLockedAmount).Mul(denominator).Quo(numerator)
	return types.NewDecCoinFromDec(lockSymbol, requiredAmount), nil
}

const errMsg = "hasn't locked"

func statisticsOurLockedCoinInPool(cli *gosdk.Client, accs []common.Account) (types.Dec, error) {
	totalAmount := types.NewDecCoinFromDec(lockSymbol, types.ZeroDec())
	for _, acc := range accs {
		lockInfo, err := cli.Farm().QueryLockInfo(poolName, acc.Address)
		if err != nil {
			if strings.Contains(err.Error(), errMsg) {
				continue
			} else {
				return types.ZeroDec(), fmt.Errorf("failed to query %s lock-info: %s", acc.Address, err.Error())
			}
		}
		totalAmount = totalAmount.Add(lockInfo.Amount)
	}
	fmt.Println("  our total locked:", totalAmount)
	return totalAmount.Amount, nil
}

func queryFarmPool(cli *gosdk.Client) (types.Dec, error) {
	pool, err := cli.Farm().QueryPool(poolName)
	if err != nil {
		return types.ZeroDec(), err
	}
	fmt.Println("whole total locked:", pool.TotalValueLocked)
	return pool.TotalValueLocked.Amount, nil
}
