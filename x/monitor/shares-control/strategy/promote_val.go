package strategy

var _ Strategy = PromoteValStrategy{}

type PromoteValStrategy struct {
	valsToPromote []string
}

func NewPromoteValStrategy(valsToPromote []string) PromoteValStrategy {
	return PromoteValStrategy{
		valsToPromote,
	}
}

func (p PromoteValStrategy) GetSolution() {
	panic("implement me")
}
