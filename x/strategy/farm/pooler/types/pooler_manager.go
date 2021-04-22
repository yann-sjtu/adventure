package types

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/x/strategy/farm/utils"
	"github.com/okex/exchain-go-sdk"
)

var (
	// Singleton Pattern
	once        sync.Once
	poolManager PoolerManager
)

type PoolerManager map[string]*Pooler

func GetPoolerManager(mnemoPath string, num ...int) PoolerManager {
	once.Do(func() {
		fmt.Printf("Loading all poolers from %s ...\n", mnemoPath)
		// 1.get pooler keys
		poolerAccs, err := utils.GetTestAccountsFromFile(mnemoPath, num...)
		if err != nil {
			panic(err)
		}

		poolManager = make(PoolerManager, len(poolerAccs))

		// 2.build pooler
		for i, poolAcc := range poolerAccs {
			poolManager[poolAcc.GetAddress().String()] = NewPooler(poolAcc, i)
		}
	})

	return poolManager
}

func (pm PoolerManager) GetPoolerFilter() map[string]struct{} {
	filter := make(map[string]struct{}, len(pm))
	for k := range pm {
		filter[k] = struct{}{}
	}

	return filter
}

func (pm PoolerManager) DestroyPool(wg *sync.WaitGroup, senderAddrStr, poolName string) {
	if !pm.isPoolerExisted(senderAddrStr) {
		return
	}

	pm[senderAddrStr].DestroyPool(wg, poolName)
}

func (pm PoolerManager) ProvidePool(wg *sync.WaitGroup, senderAddrStr string, pool gosdk.FarmPool, currentHeight int64) {
	if !pm.isPoolerExisted(senderAddrStr) {
		return
	}

	pm[senderAddrStr].ProvidePool(wg, pool, currentHeight)
}

func (pm PoolerManager) CreateFarmPoolWithRandomTokenAndProvide(wg *sync.WaitGroup, senderAddrStr string, latestHeight int64, isInWhiteList bool) {
	if !pm.isPoolerExisted(senderAddrStr) {
		return
	}

	pm[senderAddrStr].CreateFarmPoolWithRandomTokenAndProvide(wg, latestHeight, isInWhiteList)
}

func (pm PoolerManager) isPoolerExisted(poolerAddrStr string) bool {
	if _, ok := pm[poolerAddrStr]; ok {
		return true
	}

	return false
}

func GetPoolerManagerFromFiles(path, filePrefix string) PoolerManager {
	once.Do(func() {
		fmt.Printf("Loading all poolers from files with prefix \"%s\" under the path \"%s\" ...\n", filePrefix, path)
		var poolerAccs []keys.Info

		// get all files
		filesInfo, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}

		for _, f := range filesInfo {
			if !f.IsDir() && strings.HasPrefix(f.Name(), filePrefix) {
				filePath := filepath.Join(path, f.Name())
				fmt.Printf("Loading file %s ... \n", filePath)
				// 1.get pooler keys
				partAccs, err := utils.GetTestAccountsFromFile(filePath)
				if err != nil {
					panic(err)
				}
				poolerAccs = append(poolerAccs, partAccs...)
			}
		}

		poolManager = make(PoolerManager, len(poolerAccs))

		// 2.build pooler
		for i, poolAcc := range poolerAccs {
			poolManager[poolAcc.GetAddress().String()] = NewPooler(poolAcc, i)
		}
	})

	return poolManager
}
