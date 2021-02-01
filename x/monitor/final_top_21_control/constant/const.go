package constant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	RoundInterval            = 2 * time.Minute
	IntervalAfterTxBroadcast = 1 * time.Minute
)

var (
	ReservedFee = sdk.OneDec()
)