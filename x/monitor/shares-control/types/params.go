package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Params struct {
	valNumberInTop21  sdk.Dec
	percentToPlunder  sdk.Dec
	percentToDominate sdk.Dec
}

func NewParams(valNumberInTop21, percentToPlunder, percentToDominate sdk.Dec) Params {
	return Params{
		valNumberInTop21,
		percentToPlunder,
		percentToDominate,
	}
}

func (p *Params) GetValNumberInTop21() sdk.Dec {
	return p.valNumberInTop21
}

func (p *Params) GetPercentToPlunder() sdk.Dec {
	return p.percentToPlunder
}

func (p *Params) GetPercentToDominate() sdk.Dec {
	return p.percentToDominate
}
