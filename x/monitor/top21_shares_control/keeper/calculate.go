package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/top21_shares_control/constant"
	utils "github.com/okex/adventure/x/monitor/top21_shares_control/utils"
	"log"
	"time"
)

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

func (k *Keeper) GetTargetValAddrsStrToPromote(limitShares sdk.Dec) []string {
	var targetValAddrsStrToPromote []string

	for targetValAddr, shares := range k.data.TargetValSharesMap {
		if shares.LTE(limitShares) {
			targetValAddrsStrToPromote = append(targetValAddrsStrToPromote, targetValAddr)
		}
	}

	n := len(targetValAddrsStrToPromote)
	if n != 0 {
		log.Printf("%d target vals %s need to be promoted\n", n, targetValAddrsStrToPromote)
	}

	return targetValAddrsStrToPromote
}

func (k *Keeper) GetSharesToPromote(valAddrsStrToPromote []string, limitShares sdk.Dec) sdk.Dec {
	// 1. get the lowest shares in valAddrsStrToPromote
	var lowestShares sdk.Dec
	for i, valAddrStr := range valAddrsStrToPromote {
		shares := k.data.TargetValSharesMap[valAddrStr]
		if i == 0 {
			lowestShares = shares
		} else if shares.LT(lowestShares) {
			lowestShares = shares
		}
	}

	// 2. get the shares to promote
	sharesToPromte := limitShares.Sub(lowestShares).Add(constant.RedundancySharesToPromote)

	log.Printf("lowest shares of target validator: [%s]     shares to promote: [%s]\n", lowestShares, sharesToPromte)
	return sharesToPromte
}

func (k *Keeper) PickWorker(valAddrsStrToPromote []string) []mntcmn.Worker {
	var workersList []string
	for _, valAddr := range valAddrsStrToPromote {
		workersList = append(workersList, k.workersSchedule[valAddr])
	}

	workersList = utils.RemoveDuplicate(workersList)
	var workers []mntcmn.Worker
	for _, workerAddr := range workersList {
		workers = append(workers, k.workers[workerAddr])
	}

	log.Printf("%d worker %s is ready\n", len(workersList), workersList)
	return workers
}

func (k *Keeper) CalculateTokenToDeposit(shares sdk.Dec) sdk.SysCoin {
	token := sdk.SysCoin{
		Denom:  common.NativeToken,
		Amount: utils.ReverseSharesIntoToken(shares, time.Now().Unix()),
	}

	log.Printf("it's expected to deposit %s to promote the target validators\n", token.String())
	return token
}
