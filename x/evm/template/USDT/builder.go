package USDT

import (
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
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

func BuildUSDTTransferPayload(toAddr string, num int64) []byte {
	payload, err := USDTBuilder.Build("transfer", utils.EthAddress(toAddr), sdk.NewDec(num).Int)
	if err != nil {
		panic(err)
	}
	return payload
}
