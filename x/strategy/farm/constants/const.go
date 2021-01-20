package constants

import "time"

const (
	BlockInterval = 5 * time.Second
	// add liquidity params
	MinLiquidity     = "0.001"
	IssuedAmount     = 100000
	StableCoinAmount = 1
	// get lpt
	Times    = 3000
	Duration = "30s"

	// create farm pool
	MinLockAmount = "0.00000001"
	// provide farm pool
	IssuedTokenAmountSupply  = 11111.11111
	YieldAmountPerBlock      = 111.11
	StartYieldHeightInterval = 150
	RandomRange              = 500
	// ----------------- local test param ----------------- //
	//YieldAmountPerBlock      = 1111.11
	//StartYieldHeightInterval = 10
	//RandomRange              = 5

	// strategy pooler
	SleepSecondPerRoundStrategyPooler       = 10 * time.Second
	SleepSecondAfterOperationOfExpiredPools = 10 * time.Second
	TxNumOneTime                            = 50
	SleepTimeBtwGroupBroadcast              = 500 * time.Millisecond

	// strategy locker
	SleepSecondPerRoundStrategyLockAndUnlock       = 5 * time.Second
	PickedNumOneRound                              = 50
	AmountToLock                                   = 0.001
	AmountDividerToUnlock                    int64 = 2
)
