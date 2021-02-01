package constant

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	RoundInterval            = 2 * time.Minute
	IntervalAfterTxBroadcast = 2 * time.Minute
)

var (
	RedundancySharesToPromote sdk.Dec
	ReservedFee               = sdk.OneDec()
)

func init() {
	RedundancySharesToPromote, _ = sdk.NewDecFromStr("9876543210.123456789")
}
