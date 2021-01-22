package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Params struct {
	expectedValNumberInTop21  sdk.Dec
	expectedPercentToPlunder  sdk.Dec
	expectedPercentToDominate sdk.Dec
}

func NewParams(valNumberInTop21, percentToPlunder, percentToDominate sdk.Dec) Params {
	return Params{
		valNumberInTop21,
		percentToPlunder,
		percentToDominate,
	}
}

func (p *Params) GetExpectedValNumberInTop21() sdk.Dec {
	return p.expectedValNumberInTop21
}

func (p *Params) GetExpectedPercentToPlunder() sdk.Dec {
	return p.expectedPercentToPlunder
}

func (p *Params) GetExpectedPercentToDominate() sdk.Dec {
	return p.expectedPercentToDominate
}
