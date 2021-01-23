package types

import "github.com/okex/adventure/x/monitor/shares-control/strategy"

type AnalyseResult struct {
	code     int
	strategy strategy.Strategy
}

func NewAnalyseResult(code int, strategy strategy.Strategy) AnalyseResult {
	return AnalyseResult{
		code,
		strategy,
	}
}

func (ar *AnalyseResult) GetCode() int {
	return ar.code
}
