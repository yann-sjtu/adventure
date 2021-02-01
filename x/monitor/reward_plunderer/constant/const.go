package constant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	RoundInterval            = 1 * time.Minute
	IntervalAfterTxBroadcast = 2 * time.Minute
)

var (
	ReservedFee = sdk.OneDec()
)
