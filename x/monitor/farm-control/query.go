package farm_control

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/types"
	gosdk "github.com/okex/okexchain-go-sdk"
	farmTypes "github.com/okex/okexchain-go-sdk/module/farm/types"
)

const (
	poolName   = "1st_pool_okt_usdt"
	lockSymbol = "ammswap_okt_usdt-a2b"

	baseCoin = "okt"
	quoteCoin = "usdt-a2b"
)

var (
	limitRatio = types.NewDecWithPrec(81, 2)
	lockedRatio = types.NewDecWithPrec(85, 2)
)

func checkLockedRatio(cli *gosdk.Client) (types.SysCoin, error)  {
	totalOurValueLocked, err := queryAccountInPool(cli)
	if err != nil {
		return types.SysCoin{}, err
	}

	totalValueLocked, err := queryFarmPool(cli)
	if err != nil {
		return types.SysCoin{}, err
	}

	// todo
	ratio := totalOurValueLocked.Amount.Quo(totalValueLocked.Amount)
	log.Printf("total our value locked: %s, total value locked: %s, ratio: %s \n", totalOurValueLocked.String(), totalValueLocked.String(), ratio.String())
	if ratio.LTE(limitRatio) {
		totalRequiredAmount := totalValueLocked.Amount.Mul(lockedRatio)
		requiredAmount := totalRequiredAmount.Sub(totalOurValueLocked.Amount)
		return types.NewCoin(totalOurValueLocked.Denom, requiredAmount), nil
	}

	return types.NewCoin(totalOurValueLocked.Denom, types.ZeroDec()), nil
}

func queryAccountInPool(cli *gosdk.Client) (types.SysCoin, error) {
	totalAmount := types.NewCoin(lockSymbol, types.ZeroDec())
	for _, account := range accounts {
		if !account.IsLocked {
			continue
		}
		addr := account.Address

		var lockInfo farmTypes.LockInfo
		var err error
		for i := 0; i < 10; i++ {
			lockInfo, err = cli.Farm().QueryLockInfo(poolName, addr)
			if err != nil {
				continue
			}
			if lockInfo.Amount.Denom != lockSymbol {
				return types.SysCoin{}, fmt.Errorf("%s from account %s locked in pool %s is not equal with %s",
					lockInfo.Amount.String(), addr, poolName, lockSymbol)
			}
			fmt.Printf("%s[%d] locked lpt: %s", addr, account.Index, lockInfo.Amount.String())
			totalAmount = totalAmount.Add(lockInfo.Amount)
			break
		}

		if err != nil {
			return types.SysCoin{}, err
		}
	}

	return totalAmount, nil
}

func queryFarmPool(cli *gosdk.Client) (types.SysCoin, error) {
	pool, err := cli.Farm().QueryPool(poolName)
	if err != nil {
		return types.SysCoin{}, err
	}

	return pool.TotalValueLocked, nil
}


