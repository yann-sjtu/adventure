package UniswapV2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain-go-sdk/utils"
)

var (
	FactoryBuilder utils.PayloadBuilder
	WethBuilder    utils.PayloadBuilder
	RouterBuilder  utils.PayloadBuilder
)

func Init() {
	var err error

	// 1. init builders
	FactoryBuilder, err = utils.NewPayloadBuilder(FactoryBin, FactorABI)
	if err != nil {
		panic(err)
	}

	WethBuilder, err = utils.NewPayloadBuilder(WethBin, WethABI)
	if err != nil {
		panic(err)
	}

	RouterBuilder, err = utils.NewPayloadBuilder(RouterBin, RouterABI)
	if err != nil {
		panic(err)
	}
}

func BuildFactoryContractPayload(feeToSetter common.Address) []byte {
	payload, err := FactoryBuilder.Build("", feeToSetter)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildWethContractPayload() []byte {
	payload, err := WethBuilder.Build("")
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildRouterContractPayload(factoryAddress, wethAddress common.Address) []byte {
	payload, err := RouterBuilder.Build("", factoryAddress, wethAddress)
	if err != nil {
		panic(err)
	}
	return payload
}
