package types

type AnalyseResult struct {
	code              int
	valAddrsToPromote []string
}

func NewAnalyseResult(code int, valAddrsToPromote []string) AnalyseResult {
	return AnalyseResult{
		code,
		valAddrsToPromote,
	}
}

func (ar *AnalyseResult) GetCode() int {
	return ar.code
}

func (ar *AnalyseResult) GetValAddrsToPromote() []string {
	return ar.valAddrsToPromote
}
