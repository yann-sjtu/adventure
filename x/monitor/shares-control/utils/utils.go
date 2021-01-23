package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewFilter(targetValAddrs []sdk.ValAddress) map[string]struct{} {
	filter := make(map[string]struct{}, len(targetValAddrs))
	for _, valAddr := range targetValAddrs {
		filter[valAddr.String()] = struct{}{}
	}

	return filter
}

//func GetTokenFromShares(shares sdk.Dec)