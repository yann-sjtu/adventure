package types

import (
	stakingtypes "github.com/okex/okexchain/x/staking/types"
)

type Validators []stakingtypes.Validator

func (vs Validators) Len() int {
	return len(vs)
}

func (vs Validators) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs Validators) Less(i, j int) bool {
	return vs[i].DelegatorShares.GT(vs[j].DelegatorShares)
}
