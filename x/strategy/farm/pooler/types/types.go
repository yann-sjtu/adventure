package types

import (
	"fmt"

	gosdk "github.com/okex/exchain-go-sdk"
)

type FarmPool gosdk.FarmPool

func (fp FarmPool) String() string {
	return fmt.Sprintf(`FarmPool:
  Pool Name:  					    %s	
  Owner:							%s
  Min Lock Amount:      			%s
  Deposit Amount:                   %s
  Total Value Locked:               %s
  Yielded Token Infos:			    %v
  Total Accumulated Rewards:        %s`,
		fp.Name, fp.Owner, fp.MinLockAmount, fp.DepositAmount, fp.TotalValueLocked, fp.YieldedTokenInfos, fp.TotalAccumulatedRewards)
}
