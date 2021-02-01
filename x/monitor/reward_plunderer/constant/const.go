package constant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	RoundInterval            = 30 * time.Second
	IntervalAfterTxBroadcast = 1 * time.Minute
	QueryInverval            = 1 * time.Second
)

var (
	ReservedFee = sdk.OneDec()
)
