package keeper

import (
	"github.com/okex/adventure/x/monitor/shares-control/strategy"
	gosdk "github.com/okex/okexchain-go-sdk"
)

// case 1
func (k *Keeper) genStrategyToPromoteValidators(valAddrsToPromote []string, vals []gosdk.Validator) strategy.Strategy {
	stg := strategy.NewPromoteValStrategy(valAddrsToPromote)
	//TODO
	return stg

}
