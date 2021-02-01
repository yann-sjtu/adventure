package keeper

import "fmt"

func (k *Keeper) CheckPlunderedPctWarning() bool {
	curPlunderedPct := k.data.OurTotalShares.Quo(k.data.AllValsTotalShares)
	fmt.Printf(`our vals total shares:[%s]
global vals total shares: [%s]
current percentage: [%s]
target percentage: [%s]
`,
		k.data.OurTotalShares.String(),
		k.data.AllValsTotalShares.String(),
		curPlunderedPct.String(),
		k.plunderedPct.String(),
	)

	return curPlunderedPct.LT(k.plunderedPct)
}
