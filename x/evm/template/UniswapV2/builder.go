package UniswapV2

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain-go-sdk/utils"
)

var (
	FactoryBuilder utils.PayloadBuilder
	WethBuilder    utils.PayloadBuilder
	RouterBuilder  utils.PayloadBuilder
	PairBuilder    utils.PayloadBuilder
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

	PairBuilder, err = utils.NewPayloadBuilder(PairBin, PairABI)
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

func BuildAddLiquidOKTPayload(token, to string, amountTokenDesired, amountTokenMin, amountOKTMin *big.Int, deadline int64) []byte {
	payload, err := RouterBuilder.Build("addLiquidityETH",
		utils.EthAddress(token),
		utils.Uint256(amountTokenDesired), utils.Uint256(amountTokenMin),
		utils.Uint256(amountOKTMin),
		utils.EthAddress(to), big.NewInt(deadline),
	)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildAddLiquidPayload(tokenA, tokenB, to string, amountADesired, amountBDesired, amountAMin, amountBMin*big.Int, deadline int64) []byte {
	payload, err := RouterBuilder.Build("addLiquidity",
		utils.EthAddress(tokenA), utils.EthAddress(tokenB),
		utils.Uint256(amountADesired), utils.Uint256(amountBDesired),
		utils.Uint256(amountAMin), utils.Uint256(amountBMin),
		utils.EthAddress(to), big.NewInt(deadline),
	)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildSwapExactTokensForTokensPayload(amountIn, amountOut *big.Int, path []string, to string, deadline int64) []byte {
	payload, err := RouterBuilder.Build("swapExactTokensForTokens",
		amountIn, amountOut,
		utils.EthAddresses(path),
		utils.EthAddress(to), big.NewInt(deadline),
	)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildRemoveLiquidOKTPayload(token, to string, liquidity, amountTokenMin, amountOKTMin*big.Int, deadline int64) []byte {
	payload, err := RouterBuilder.Build("removeLiquidityETH",
		utils.EthAddress(token),
		utils.Uint256(liquidity),
		utils.Uint256(amountTokenMin), utils.Uint256(amountOKTMin),
		utils.EthAddress(to), big.NewInt(deadline),
	)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildRemoveLiquidPayload(tokenA, tokenB, to string, liquidity, amountAMin, amountBMin*big.Int, deadline int64) []byte {
	payload, err := RouterBuilder.Build("removeLiquidity",
		utils.EthAddress(tokenA), utils.EthAddress(tokenB),
		utils.Uint256(liquidity),
		utils.Uint256(amountAMin), utils.Uint256(amountBMin),
		utils.EthAddress(to), big.NewInt(deadline),
	)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildApprovePayload(addr string, amount int64) []byte {
	payload, err := PairBuilder.Build("approve", utils.EthAddress(addr), sdk.NewDec(amount).Int)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildGetReversesPayload() []byte {
	payload, err := PairBuilder.Build("getReserves")
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildWethApprovePayload(addr string, amount int64) []byte {
	payload, err := WethBuilder.Build("approve", utils.EthAddress(addr), sdk.NewDec(amount).Int)
	if err != nil {
		panic(err)
	}
	return payload
}

func BuildWethDepositPayload() []byte {
	payload, err := WethBuilder.Build("deposit")
	if err != nil {
		panic(err)
	}
	return payload
}