package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Data struct {
	Vals               Validators
	OurTotalShares     sdk.Dec
	AllValsTotalShares sdk.Dec
}
