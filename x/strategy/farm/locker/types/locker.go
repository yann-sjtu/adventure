package types

import (
	"fmt"
	"log"
	"sync"

	gokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/farm/client"
	"github.com/okex/adventure/x/strategy/farm/constants"
	gosdk "github.com/okex/okexchain-go-sdk"
)

type Locker struct {
	id      int
	key     gokeys.Info
	accAddr string
}

func NewLocker(info gokeys.Info, id int) *Locker {
	return &Locker{
		id:      id,
		accAddr: info.GetAddress().String(),
		key:     info,
	}
}

func (l *Locker) UnlockFromOnePool(wg *sync.WaitGroup, poolName string, isAll bool) {
	defer wg.Done()
	// get the lock info from query
	cli := client.CliManager.GetClient()
	lockInfo, err := cli.Farm().QueryLockInfo(poolName, l.accAddr)
	if err != nil {
		return
	}

	// get accInfo
	accInfo, err := cli.Auth().QueryAccount(l.accAddr)
	if err != nil {
		return
	}

	// unlock locked token from the pool
	if !isAll {
		lockInfo.Amount.Amount = lockInfo.Amount.Amount.QuoInt64(constants.AmountDividerToUnlock)
	}
	amountStr := lockInfo.Amount.String()
	if _, err := cli.Farm().Unlock(l.key, common.PassWord, poolName, amountStr, "", accInfo.GetAccountNumber(),
		accInfo.GetSequence()); err != nil {
		log.Printf("Tx error. %s unlocks %s from pool %s: %s\n", l.accAddr, amountStr, poolName, err)
		return
	}

	log.Printf("%s unlocks %s from pool %s successfully\n", l.accAddr, amountStr, poolName)
}

func (l *Locker) LockOnOnePoolByOneLocker(wg *sync.WaitGroup, pool gosdk.FarmPool) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	// get accInfo
	accInfo, err := cli.Auth().QueryAccount(l.accAddr)
	if err != nil {
		return
	}

	// lock some token on the pool
	amountStr := fmt.Sprintf("%f%s", constants.AmountToLock, pool.MinLockAmount.Denom)
	if _, err := cli.Farm().Lock(l.key, common.PassWord, pool.Name, amountStr, "", accInfo.GetAccountNumber(),
		accInfo.GetSequence()); err != nil {
		log.Printf("Tx error. %s locks %s on pool %s: %s\n", l.accAddr, amountStr, pool.Name, err)
		return
	}

	log.Printf("%s locks %s on pool %s successfully\n", l.accAddr, amountStr, pool.Name)
}
