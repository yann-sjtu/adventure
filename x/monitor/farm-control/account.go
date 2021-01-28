package farm_control

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	gosdk "github.com/okex/okexchain-go-sdk"
)

var (
	accounts []*FarmAccount
)

type FarmAccounts []*FarmAccount

type FarmAccount struct {
	Address    string
	Index      int
	LockedCoin types.DecCoin
}

func newFarmAddrAccounts(addrs []string, startIndex int) FarmAccounts {
	farmAccounts := make([]*FarmAccount, len(addrs), len(addrs))
	for i := 0; i < len(addrs); i++ {
		farmAccounts[i] = &FarmAccount{Address: addrs[i], Index: startIndex + i, LockedCoin: types.NewDecCoinFromDec(lockSymbol, types.ZeroDec())}
	}
	return farmAccounts
}

const errMsg = "hasn't locked"

func refreshFarmAccounts(cli *gosdk.Client) error {
	for i := 0; i < len(accounts); i++ {
		lockInfo, err := cli.Farm().QueryLockInfo(poolName, accounts[i].Address)
		if err != nil {
			if strings.Contains(err.Error(), errMsg) {
				accounts[i].LockedCoin = zeroLpt
			} else {
				return fmt.Errorf("failed to query %s lock-info: %s", accounts[i].Address, err.Error())
			}
		} else {
			accounts[i].LockedCoin = lockInfo.Amount
		}
	}

	//fmt.Printf("=== accounts on %s ===\n", poolName)
	//for i := 0; i < len(accounts); i++ {
	//	fmt.Println(accounts[i].Index, accounts[i].Address, accounts[i].LockedCoin.String())
	//}
	//fmt.Printf("======================================\n")
	return nil
}
