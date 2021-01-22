package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func (w Worker) String() string {
	return fmt.Sprintf("%s with index %d", w.accAddr.String(), w.index)
}

func (w *Worker) GetAccAddr() sdk.AccAddress {
	return w.accAddr
}

func (w *Worker) GetIndex() int {
	return w.index
}
