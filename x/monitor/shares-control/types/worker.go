package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Worker struct {
	accAddr sdk.AccAddress
	index   int
}

func NewWorker(accAddr sdk.AccAddress, index int) Worker {
	return Worker{
		accAddr,
		index,
	}
}
