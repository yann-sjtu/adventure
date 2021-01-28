package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"log"
	"math/rand"
	"time"
)

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

func (k *Keeper) GetTheHighestShares(valAddrsStr []string) sdk.Dec {
	highestShares := sdk.ZeroDec()
	for _, addr := range valAddrsStr {
		shares, ok := k.data.Top21SharesMap[addr]
		if !ok {
			log.Panicf("intruder %s ain't found in k.data.Top21SharesMap[addr]", addr)
		}

		if shares.GT(highestShares) {
			highestShares = shares
		}
	}

	return highestShares
}
