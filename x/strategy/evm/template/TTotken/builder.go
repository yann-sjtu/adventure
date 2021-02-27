package TTotken

import (
	"math/big"

	"github.com/okex/okexchain-go-sdk/utils"
)

var (
	TTokenBuilder utils.PayloadBuilder
)

func Init() {
	var err error

	// 1. init builders
	TTokenBuilder, err = utils.NewPayloadBuilder(TTokenBin, TTokenABI)
	if err != nil {
		panic(err)
	}
}

func BuildTTokenContractPayload(name, symbol string) []byte {
	payload, err := TTokenBuilder.Build("", name, symbol)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildTTokenMintPayload(addr string, amount *big.Int) []byte {
	payload, err := TTokenBuilder.Build("mint", utils.EthAddress(addr), utils.Uint256(amount))
	if err != nil {
		panic(err)
	}
	return payload
}