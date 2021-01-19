package emitter

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/okex/adventure/x/strategy/farm/client"
	"github.com/okex/adventure/x/strategy/farm/constants"
	lockertypes "github.com/okex/adventure/x/strategy/farm/locker/types"
	poolertypes "github.com/okex/adventure/x/strategy/farm/pooler/types"
	"github.com/okex/adventure/x/strategy/farm/utils"
	"github.com/okex/okexchain-go-sdk"
)

type Emitter struct {
	// PoolerManager and LockerManager: there can only be one in the Emitter
	// it will bring signature error with the existence of both Managers
	poolertypes.PoolerManager
	lockertypes.LockerManager
	selector
	txNumOneTime int
	sleepTime    time.Duration
}

func NewEmitter(pm poolertypes.PoolerManager, lm lockertypes.LockerManager) Emitter {
	if pm == nil && lm == nil {
		panic("no manager input")
	}

	if !(pm == nil || lm == nil) {
		panic("PoolerManager and LockerManager can't exist both")
	}

	emt := Emitter{
		PoolerManager: pm,
		LockerManager: lm,
		txNumOneTime:  constants.TxNumOneTime,
		sleepTime:     constants.SleepTimeBtwGroupBroadcast,
	}

	if lm != nil {
		// init selector
		lockerManagerLen := len(lm)
		if lockerManagerLen < constants.PickedNumOneRound {
			panic("init selector failed. picked number is larger than the length of LockerManager")
		}
		emt.pickedNumOneRound = constants.PickedNumOneRound
		emt.lockerAddrsList = make([]string, lockerManagerLen)
		emt.selectedAddrsList = make([]string, emt.pickedNumOneRound)
		var index int
		for addrStr := range lm {
			emt.lockerAddrsList[index] = addrStr
			index++
		}
	}

	return emt
}

func (emt *Emitter) PickLockersRandomly() {
	maxRange := len(emt.lockerAddrsList)
	rand.Seed(time.Now().UnixNano())
	luckyNum := rand.Intn(maxRange)
	// picked lockers randomly
	for i := 0; i < emt.pickedNumOneRound; i++ {
		// update selectedAddrsList
		emt.selectedAddrsList[i] = emt.lockerAddrsList[(luckyNum+i)%maxRange]
	}

	fmt.Printf(`
============================================================
|              %d lockers are picked randomly               |
============================================================

`, emt.pickedNumOneRound)
	for _, addrStr := range emt.selectedAddrsList {
		fmt.Println(addrStr)
	}

	return
}

func (emt *Emitter) GetPoolersWithoutPools() (accAddrStrs []string, err error) {
	pools, err := utils.QueryAllAdventurePools()
	if err != nil {
		return
	}

	filter := emt.PoolerManager.GetPoolerFilter()
	for _, pool := range pools {
		delete(filter, pool.Owner.String())
	}

	for addr := range filter {
		accAddrStrs = append(accAddrStrs, addr)
	}

	return
}

func (emt *Emitter) CreateAndProvidePool(poolerAddrList []string, latestHeight int64) {
	times := len(poolerAddrList)/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = len(poolerAddrList)
		}

		//fmt.Printf("Group %d: create and provide pools %d ~ %d ...\n", i, index1, index2)
		emt.createAndProvidePoolByGroup(poolerAddrList[index1:index2], latestHeight)
		time.Sleep(emt.sleepTime)
	}
}

func (emt *Emitter) createAndProvidePoolByGroup(poolerAddrList []string, latestHeight int64) {
	var wg sync.WaitGroup
	for _, poolerAddrStr := range poolerAddrList {
		// create and provide pool by pooler
		// pool in white list randomly
		wg.Add(1)
		go emt.CreateFarmPoolWithRandomTokenAndProvide(&wg, poolerAddrStr, latestHeight, utils.GetRandomBool())
	}
	wg.Wait()
}

func (emt *Emitter) ProvidePools(pools []gosdk.FarmPool, currentHeight int64) {
	times := len(pools)/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = len(pools)
		}

		//fmt.Printf("Group %d: provide pools %d ~ %d ...\n", i, index1, index2)
		emt.providePoolsByGroup(pools[index1:index2], currentHeight)
		time.Sleep(emt.sleepTime)
	}
}

func (emt *Emitter) providePoolsByGroup(pools []gosdk.FarmPool, currentHeight int64) {
	var wg sync.WaitGroup
	for _, pool := range pools {
		// provide pool by pooler
		wg.Add(1)
		go emt.ProvidePool(&wg, pool.Owner.String(), pool, currentHeight)
	}
	wg.Wait()
}

func (emt *Emitter) DestroyPools(pools []gosdk.FarmPool) {
	times := len(pools)/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = len(pools)
		}

		//fmt.Printf("Group %d: destroy pools %d ~ %d ...\n", i, index1, index2)
		emt.destroyPoolsByGroup(pools[index1:index2])
		time.Sleep(emt.sleepTime)
	}
}

func (emt *Emitter) destroyPoolsByGroup(pools []gosdk.FarmPool) {
	var wg sync.WaitGroup
	for _, pool := range pools {
		// destroy pool by pooler
		wg.Add(1)
		go emt.DestroyPool(&wg, pool.Owner.String(), pool.Name)
	}
	wg.Wait()
}

func (emt *Emitter) UnlockFromExpiredPools() error {
	// 1.get the expired pools on current height
	expiredPools, _, err := utils.GetExpiredPoolsOnCurrentHeight()
	if err != nil {
		return err
	}

	if len(expiredPools) != 0 {
		// 2.lockers unlock all from the expired pools
		emt.unlockAllTokensFromPools(expiredPools)
	}

	return nil
}

func (emt *Emitter) unlockAllTokensFromPools(pools []gosdk.FarmPool) {
	loggerUnlockAllTokensFromPools(pools)

	times := len(pools)/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = len(pools)
		}

		//fmt.Printf("Group %d: unlock pool %d ~ %d ...\n", i, index1, index2)
		emt.unlockAllTokensFromPoolsByGroup(pools[index1:index2])
		time.Sleep(emt.sleepTime)
	}
}

func loggerUnlockAllTokensFromPools(pools []gosdk.FarmPool) {
	// print title
	fmt.Printf(`
============================================================
| lockers unlock all locked token from the %d expired pools |
============================================================

`, len(pools))

	for _, pool := range pools {
		fmt.Println(pool.Name)
	}
	fmt.Printf("\n")
}

func (emt *Emitter) unlockAllTokensFromPoolsByGroup(pools []gosdk.FarmPool) {
	for _, pool := range pools {
		// get locker info
		lockerAddrs, err := client.CliManager.GetClient().Farm().QueryAccountsLockedTo(pool.Name)
		if err != nil {
			return
		}

		var wg sync.WaitGroup
		// locker unlocks the token
		for _, lockerAddr := range lockerAddrs {
			wg.Add(1)
			go emt.UnlockAllFromOnePool(&wg, lockerAddr.String(), pool.Name)
		}
		wg.Wait()
	}
}

func (emt *Emitter) LockToRandomPoolsByRandomLockers(pools []gosdk.FarmPool) {
	loggerLockToRandomPoolsByRandomLockers()

	times := emt.pickedNumOneRound/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = emt.pickedNumOneRound
		}

		//fmt.Printf("Group %d: lock to random pools by random locker %d ~ %d ...\n", i, index1, index2)
		emt.lockToRandomPoolsByRandomLockersByGroup(emt.selectedAddrsList[index1:index2], pools)
		time.Sleep(emt.sleepTime)
	}
}

func (emt *Emitter) LockToTargetPool(pool gosdk.FarmPool) {
	loggerLockToTargetPool()

	times := emt.pickedNumOneRound/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = emt.pickedNumOneRound
		}

		emt.lockToRandomPoolsByRandomLockersByGroup(emt.selectedAddrsList[index1:index2], []gosdk.FarmPool{pool})
		time.Sleep(emt.sleepTime)
	}
}

func loggerLockToTargetPool() {
	// print title
	fmt.Printf(`
============================================================
|    picked lockers lock tokens on the target pool         |
============================================================

`)
}

func (emt *Emitter) lockToRandomPoolsByRandomLockersByGroup(selectedAddrsList []string, pools []gosdk.FarmPool) {
	maxRange := len(pools)
	var wg sync.WaitGroup
	for _, lockerAddrStr := range selectedAddrsList {
		// a picked locker locks on a random pool
		rand.Seed(time.Now().UnixNano())
		luckyNum := rand.Intn(maxRange)
		wg.Add(1)
		go emt.LockOnOnePoolByOneLocker(&wg, lockerAddrStr, pools[luckyNum])
	}
	wg.Wait()
}

func loggerLockToRandomPoolsByRandomLockers() {
	// print title
	fmt.Printf(`
============================================================
|    picked lockers lock tokens on picked pools randomly   |
============================================================

`)
}

func (emt *Emitter) UnlockFromPoolsLockedBeforeByRandomLockers() {
	loggerUnlockFromPoolsLockedBeforeByRandomLockers()

	times := emt.pickedNumOneRound/emt.txNumOneTime + 1
	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * emt.txNumOneTime
		if i != times-1 {
			index2 = (i + 1) * emt.txNumOneTime
		} else {
			index2 = emt.pickedNumOneRound
		}

		//fmt.Printf("Group %d: unlock from pools that were locked on before by random locker %d ~ %d ...\n", i, index1, index2)
		emt.unlockFromPoolsLockedBeforeByRandomLockersByGroup(emt.selectedAddrsList[index1:index2])
		time.Sleep(emt.sleepTime)
	}
}

func (emt *Emitter) unlockFromPoolsLockedBeforeByRandomLockersByGroup(selectedAddrsList []string) {
	var wg sync.WaitGroup
	for _, lockerAddrStr := range selectedAddrsList {
		// get pool names that the locker locked on before
		poolNames, err := client.CliManager.GetClient().Farm().QueryAccount(lockerAddrStr)
		if err != nil {
			return
		}

		if len(poolNames) == 0 {
			// the locker didn't lock on any pools before
			return
		}

		// choose a pool randomly
		rand.Seed(time.Now().UnixNano())
		luckyNum := rand.Intn(len(poolNames))

		wg.Add(1)
		go emt.UnlockPartFromOnePool(&wg, lockerAddrStr, poolNames[luckyNum])
	}
	wg.Wait()

}

func loggerUnlockFromPoolsLockedBeforeByRandomLockers() {
	// print title
	fmt.Printf(`
============================================================
|   picked lockers unlock tokens from pools locked before  |
============================================================

`)
}
