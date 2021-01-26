package farm_rm_liquidity

import "time"

const (
	queryInterval = 500 * time.Millisecond
	poolName      = "1st_pool_okt_usdt"
	lockSymbol    = "ammswap_okt_usdt-a2b"

	baseCoin         = "okt"
	quoteCoin        = "usdt-a2b"
	workerStartIndex = 901
)
