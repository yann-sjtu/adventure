package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Data struct {
	Vals               Validators
	OurTotalShares     sdk.Dec
	EnemyTotalShares   sdk.Dec
	EnemyLowestShares  sdk.Dec
	Top21SharesMap     map[string]sdk.Dec
	TargetValSharesMap map[string]sdk.Dec
}
