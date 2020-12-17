package types

import (
	"fmt"
	"sync"

	"github.com/okex/adventure/x/strategy/farm/utils"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	// Singleton Pattern
	once          sync.Once
	lockerManager LockerManager
)

type LockerManager map[string]*Locker

func GetLockerManager(mnemoPath string) LockerManager {
	once.Do(func() {
		fmt.Printf("Loading all lockers from %s ...\n", mnemoPath)
		// 1.get locker keys
		lockerAccs, err := utils.GetTestAccountsFromFile(mnemoPath)
		if err != nil {
			panic(err)
		}

		lockerManager = make(LockerManager, len(lockerAccs))

		// 2.build pooler
		for i, lockerAcc := range lockerAccs {
			lockerManager[lockerAcc.GetAddress().String()] = NewLocker(lockerAcc, i)
		}
	})

	return lockerManager
}

func (lm LockerManager) UnlockAllFromOnePool(wg *sync.WaitGroup, senderAddrStr, poolName string) {
	if !lm.isLockerExisted(senderAddrStr) {
		return
	}

	lm[senderAddrStr].UnlockFromOnePool(wg, poolName, true)
}

func (lm LockerManager) UnlockPartFromOnePool(wg *sync.WaitGroup, senderAddrStr, poolName string) {
	if !lm.isLockerExisted(senderAddrStr) {
		return
	}

	lm[senderAddrStr].UnlockFromOnePool(wg, poolName, false)
}

func (lm LockerManager) LockOnOnePoolByOneLocker(wg *sync.WaitGroup, senderAddrStr string, pool gosdk.FarmPool) {
	if !lm.isLockerExisted(senderAddrStr) {
		return
	}

	lm[senderAddrStr].LockOnOnePoolByOneLocker(wg, pool)
}

func (lm LockerManager) isLockerExisted(lockerAddrStr string) bool {
	if _, ok := lm[lockerAddrStr]; ok {
		return true
	}

	return false
}
