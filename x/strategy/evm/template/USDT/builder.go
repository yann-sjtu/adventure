package USDT

import (
	"math/big"

	"github.com/okex/okexchain-go-sdk/utils"
)

var (
	USDTBuilder utils.PayloadBuilder
)

func Init() {
	var err error

	// 1. init builders
	USDTBuilder, err = utils.NewPayloadBuilder(USDTBin, USDTABI)
	if err != nil {
		panic(err)
	}
}

func BuildUSDTContractPayload(initalSupply, decimals *big.Int, name, symbol string) []byte {
	payload, err := USDTBuilder.Build("", initalSupply, name, symbol, decimals)
	if err != nil {
		panic(err)
	}
	return payload
}