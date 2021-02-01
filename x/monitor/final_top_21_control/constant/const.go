package constant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	RoundInterval            = 2 * time.Second
	IntervalAfterTxBroadcast = 2 * time.Minute
)

var (
	ReservedFee = sdk.OneDec()
)