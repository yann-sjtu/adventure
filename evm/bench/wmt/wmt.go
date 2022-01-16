package wmt

import (
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/evm/bench/utils"
	"github.com/okex/adventure/evm/constant"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var (
	LPAddrs    = [4]ethcmm.Address{OktUsdtLPAddr, OktDotkLPAddr, OktBtckLPAddr, OktEthkLPAddr}
	PoolAddrs  = [4]ethcmm.Address{OktUsdtPoolAddr, OktDotkPoolAddr, OktBtckPoolAddr, OktEthkPoolAddr}
	TokenAddrs = [4]ethcmm.Address{UsdtAddr, DotkAddr, BtckAddr, EthkAddr}
)

var (
	routerAddr = ethcmm.HexToAddress("0x2CA0E1278B9D7A967967d3C707b81C72FC180CaF")

	OktUsdtLPAddr   = ethcmm.HexToAddress("0xe922FF7B02672bB59A64b90864FC5e511AD4d5fa")
	OktUsdtPoolAddr = ethcmm.HexToAddress("0x5aFC0E1ddDd7a5151d83a3385C01e6159539a37C")

	OktDotkLPAddr   = ethcmm.HexToAddress("0x1908839fF3292314Cf1B18D1EF72AF54a0c7F6FE")
	OktDotkPoolAddr = ethcmm.HexToAddress("0x844f80e679BA02C7408319E87FDAe8bEB128c831")

	OktBtckLPAddr   = ethcmm.HexToAddress("0x73Da05c587ECA1b36dD07e293AC00FEc9D887C88")
	OktBtckPoolAddr = ethcmm.HexToAddress("0xc5B011Ef3b5Bad391dd34Af2Da67Af0a7b8d5730")

	OktEthkLPAddr   = ethcmm.HexToAddress("0x45ca0ae81c65249a93a9f7b60BDE707B26217E5D")
	OktEthkPoolAddr = ethcmm.HexToAddress("0x4D8bC6D21E478BB34F72548906303BaD60f2a560")

	UniUsdtLPAddr   = ethcmm.HexToAddress("0xfc56c01730f1d47cd187253353521d3dc2218a82")
	UniUsdtPoolAddr = ethcmm.HexToAddress("0xaAFd4b09e0c275b3EC35B3cacB99D6DA9Ca96E33")

	UsdtAddr = ethcmm.HexToAddress("0xee666e967293094007d7c3718044e07565b1f8a9")
	WethAddr = ethcmm.HexToAddress("0x70c1c53E991F31981d592C2d865383AC0d212225")
	WoktAddr = ethcmm.HexToAddress("0x2789Fdc29D0f1D2ddaC362B2cb79F7799A5fbdAF")
	UniAddr  = ethcmm.HexToAddress("0x0A1D36fCD446Df6bA0bA326bec5291417B97757d")
	OkbAddr  = ethcmm.HexToAddress("0xa860E9929B7DE53218c9B0a555680587D3542882")
	EthkAddr = ethcmm.HexToAddress("0x01490F1bAfE4ab9eE1F61454Bb295502ab5c3fDD")
	BtckAddr = ethcmm.HexToAddress("0xFd71e3597462ed133Ce5CDfB62041D164d1EBC99")
	UsdcAddr = ethcmm.HexToAddress("0x7B334746E0B9f7fbD94AD9f4eA9e304e1d2dF0DA")
	FilkAddr = ethcmm.HexToAddress("0x33c548B01c04D195E99c16C6dC1D4E9252EE45ea")
	DotkAddr = ethcmm.HexToAddress("0xe2017Ea8AE91108B968685cF743F2ED8Da178A13")
	LtckAddr = ethcmm.HexToAddress("0xA51E71874112cd7fa7885C23D403525Ee0F73c80")
	UsdkAddr = ethcmm.HexToAddress("0xcBCc53b501A799Dd90D6546aa5319cF87a3E66fa")
)

func wmt(cmd *cobra.Command, args []string) {
	_, wethABI, routerABI, _, stakingRewardABI := generateABI()
	poolAddr, usdtAddr := PoolAddrs[0], TokenAddrs[0]

	depositAmount := sdk.MustNewDecFromStr("0.00001").Int
	depositData, _ := wethABI.Pack("deposit")
	approveToRouterData, _ := wethABI.Pack("approve", routerAddr, sdk.NewDec(10).Int)
	approveToPoolData, _ := wethABI.Pack("approve", poolAddr, sdk.NewDec(10).Int)
	stakeData, _ := stakingRewardABI.Pack("stake", big.NewInt(100))
	getRewardData, _ := stakingRewardABI.Pack("getReward")
	withdrawData, _ := stakingRewardABI.Pack("withdraw", big.NewInt(20))
	exitData, _ := stakingRewardABI.Pack("exit")

	gasLimit := uint64(500000)
	utils.RunTxs(
		utils.DefaultBaseParamFromFlag(),
		func(addr ethcmm.Address) []utils.TxParam {
			// 0.1 deposit okt into wokt, get wokt
			depositTxParam := utils.NewTxParam(WethAddr, depositAmount, gasLimit, evmtypes.DefaultGasPrice, depositData)
			// 0.2 approve wokt to router
			approveWoktTxParam := utils.NewTxParam(WethAddr, nil, gasLimit, evmtypes.DefaultGasPrice, approveToRouterData)
			// 0.3 swap wokt -> usdt in router
			swapData, _ := routerABI.Pack("swapExactTokensForTokens", big.NewInt(10000), big.NewInt(0), []ethcmm.Address{WethAddr, UsdtAddr}, addr, big.NewInt(time.Now().Add(time.Hour*8640).Unix()))
			swapTxParam := utils.NewTxParam(routerAddr, nil, gasLimit, evmtypes.DefaultGasPrice, swapData)

			// 1.1 approve usdt to router
			approveUsdtTxParam := utils.NewTxParam(usdtAddr, nil, gasLimit, evmtypes.DefaultGasPrice, approveToRouterData)
			// 1.2 add (usdt-wokt) liquidity
			addLiquidityData, _ := routerABI.Pack("addLiquidityETH", usdtAddr, sdk.MustNewDecFromStr("0.00000000001").Int, big.NewInt(0), big.NewInt(0), addr, big.NewInt(time.Now().Add(time.Hour*8640).Unix()))
			addLiquidityTxParam := utils.NewTxParam(routerAddr, sdk.MustNewDecFromStr("0.00000000000001").Int, gasLimit, evmtypes.DefaultGasPrice, addLiquidityData)

			// 2.1 approve uni(usdt-wokt) to staking pool
			approveUniTokenTxParam := utils.NewTxParam(OktUsdtLPAddr, nil, gasLimit, evmtypes.DefaultGasPrice, approveToPoolData)
			// 2.2 stake uni(usdt-wokt) to staking pool
			stakeTxParam := utils.NewTxParam(poolAddr, nil, gasLimit, evmtypes.DefaultGasPrice, stakeData)

			var randomTxParam []utils.TxParam
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(10) <= 3 { // 2.3 withDraw randomly
				withdrawTxParam := utils.NewTxParam(poolAddr, nil, gasLimit, evmtypes.DefaultGasPrice, withdrawData)
				randomTxParam = append(randomTxParam, withdrawTxParam)
			}
			if rand.Intn(10) <= 3 { // 2.4 get Reward randomly
				getRewardTxParam := utils.NewTxParam(poolAddr, nil, gasLimit, evmtypes.DefaultGasPrice, getRewardData)
				randomTxParam = append(randomTxParam, getRewardTxParam)
			}
			if rand.Intn(10) <= 2 { // 2.5 Exit randomly
				exitTxParam := utils.NewTxParam(poolAddr, nil, gasLimit, evmtypes.DefaultGasPrice, exitData)
				randomTxParam = append(randomTxParam, exitTxParam)
			}

			return append(
				[]utils.TxParam{depositTxParam, approveWoktTxParam, swapTxParam, approveUsdtTxParam, addLiquidityTxParam, approveUniTokenTxParam, stakeTxParam},
				randomTxParam...)
		},
	)
}

func generateABI() (factoryABI, wethABI, routerABI, pairABI, stakingRewardABI abi.ABI) {
	var err error
	factoryABI, err = abi.JSON(strings.NewReader(constant.UniswapFactorABI))
	if err != nil {
		panic(err)
	}

	wethABI, err = abi.JSON(strings.NewReader(constant.WethABI))
	if err != nil {
		panic(err)
	}

	routerABI, err = abi.JSON(strings.NewReader(constant.UniswapRouterABI))
	if err != nil {
		panic(err)
	}

	pairABI, err = abi.JSON(strings.NewReader(constant.UniswapPairABI))
	if err != nil {
		panic(err)
	}

	stakingRewardABI, err = abi.JSON(strings.NewReader(constant.StakingRewardsABI))
	if err != nil {
		panic(err)
	}

	return
}
