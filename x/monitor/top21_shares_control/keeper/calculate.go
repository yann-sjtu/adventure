package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/top21_shares_control/utils"
	"math/rand"
	"time"
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

func (k *Keeper) CalculateHowMuchToDeposit(enemyTotalShares, tarValsCurrentTotalShares sdk.Dec) (worker mntcmn.Worker, depositCoin sdk.SysCoin, err error) {
	expectedTotalShares := enemyTotalShares.Quo(sdk.OneDec().Sub(k.dominationPct))
	expectedTarValsTotalShares := expectedTotalShares.Mul(k.dominationPct)

	// pick worker
	worker, valNum, err := k.pickWorker()
	if err != nil {
		return
	}

	// get coin for required shares
	sharesRequired := expectedTarValsTotalShares.Sub(tarValsCurrentTotalShares)
	if !sharesRequired.IsPositive() {
		// TODO: add withdraw delegation function
		return worker, depositCoin, fmt.Errorf("required shares is not positive, target vals shares: current[%s] expected[%s]", tarValsCurrentTotalShares.String(), expectedTarValsTotalShares.String())
	}

	depositCoin = sdk.NewDecCoinFromDec("okt", utils.ReverseSharesIntoToken(sharesRequired.QuoInt64(int64(valNum)), time.Now().Unix()))
	return
}

func (k *Keeper) pickWorker() (worker mntcmn.Worker, valNum int, err error) {
	// pick worker randomly
	rand.Seed(time.Now().UnixNano())
	worker = k.workers[rand.Intn(len(k.workers))]
	delegator, err := k.cliManager.GetClient().Staking().QueryDelegator(worker.GetAccAddr().String())
	if err != nil {
		return
	}

	valNum = len(delegator.ValidatorAddresses)
	if valNum == 0 {
		return worker, valNum, fmt.Errorf("worker [%s] hasn't voted yet", delegator.DelegatorAddress.String())
	}

	return
}

func (k *Keeper) HandleStrategy() (sdk.DecCoin, []string, error) {
	// pick worker randomly
	rand.Seed(time.Now().UnixNano())
	worker := k.workers[rand.Intn(len(k.workers))]
	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(worker.GetAccAddr().String())
	if err != nil {
		return sdk.DecCoin{}, []string{}, err
	}

	_ = accInfo

	return sdk.DecCoin{}, []string{}, nil
}
