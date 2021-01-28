package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
)

type Data struct {
	Vals             []stakingtypes.Validator
	OurTotalShares   sdk.Dec
	EnemyTotalShares sdk.Dec
}
