package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Params struct {
	valNumberInTop21 sdk.Dec
	percentToPlunder sdk.Dec
}

func NewParams(valNumberInTop21, percentToPlunder sdk.Dec) Params {
	return Params{
		valNumberInTop21,
		percentToPlunder,
	}
}
