package uniswap_operate

import (
	"log"
	"math/rand"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/deploy-contracts"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2Staker"
	"github.com/okex/adventure/x/strategy/evm/tools"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func UniswapTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uniswap-testnet-operate",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.IntVarP(&deploy_contracts.GoroutineNum, "goroutine-num", "g", 1, "set Goroutine Num of deploying contracts")
	flags.StringVarP(&deploy_contracts.MnemonicPath, "mnemonic-path", "m", "", "set the MnemonicPath path")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 5, "set the sleep time")

	return cmd
}

var (
	sleepTime int
)

var (
	LPAddrs    = [4]string{OktUsdtLPAddr, OktDotkLPAddr, OktBtckLPAddr, OktEthkLPAddr}
	PoolAddrs  = [4]string{OktUsdtPoolAddr, OktDotkPoolAddr, OktBtckPoolAddr, OktEthkPoolAddr}
	TokenAddrs = [4]string{UsdtAddr, DotkAddr, BtckAddr, EthkAddr}
)

const (
	routerAddr = "0x2CA0E1278B9D7A967967d3C707b81C72FC180CaF"

	OktUsdtLPAddr   = "0xe922FF7B02672bB59A64b90864FC5e511AD4d5fa"
	OktUsdtPoolAddr = "0x5aFC0E1ddDd7a5151d83a3385C01e6159539a37C"

	OktDotkLPAddr   = "0x1908839fF3292314Cf1B18D1EF72AF54a0c7F6FE"
	OktDotkPoolAddr = "0x844f80e679BA02C7408319E87FDAe8bEB128c831"

	OktBtckLPAddr   = "0x73Da05c587ECA1b36dD07e293AC00FEc9D887C88"
	OktBtckPoolAddr = "0xc5B011Ef3b5Bad391dd34Af2Da67Af0a7b8d5730"

	OktEthkLPAddr   = "0x45ca0ae81c65249a93a9f7b60BDE707B26217E5D"
	OktEthkPoolAddr = "0x4D8bC6D21E478BB34F72548906303BaD60f2a560"

	UniUsdtLPAddr   = "0xfc56c01730f1d47cd187253353521d3dc2218a82"
	UniUsdtPoolAddr = "0xaAFd4b09e0c275b3EC35B3cacB99D6DA9Ca96E33"

	UsdtAddr = "0xee666e967293094007d7c3718044e07565b1f8a9"
	WoktAddr = "0x2789Fdc29D0f1D2ddaC362B2cb79F7799A5fbdAF"
	UniAddr  = "0x0A1D36fCD446Df6bA0bA326bec5291417B97757d"
	OkbAddr  = "0xa860E9929B7DE53218c9B0a555680587D3542882"
	EthkAddr = "0x01490F1bAfE4ab9eE1F61454Bb295502ab5c3fDD"
	BtckAddr = "0xFd71e3597462ed133Ce5CDfB62041D164d1EBC99"
	UsdcAddr = "0x7B334746E0B9f7fbD94AD9f4eA9e304e1d2dF0DA"
	FilkAddr = "0x33c548B01c04D195E99c16C6dC1D4E9252EE45ea"
	DotkAddr = "0xe2017Ea8AE91108B968685cF743F2ED8Da178A13"
	LtckAddr = "0xA51E71874112cd7fa7885C23D403525Ee0F73c80"
	UsdkAddr = "0xcBCc53b501A799Dd90D6546aa5319cF87a3E66fa"
)

func testLoop(cmd *cobra.Command, args []string) {
	_, poolAddr, tokenAddr := LPAddrs[0], PoolAddrs[0], TokenAddrs[0]

	infos := common.GetAccountManagerFromFile(deploy_contracts.MnemonicPath)
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.005okt", types.BroadcastBlock,5000000)

	succ, fail := tools.NewCounter(-1), tools.NewCounter(-1)
	var wg sync.WaitGroup
	for i := 0; i < deploy_contracts.GoroutineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// 1. get one of the eth private key
				info := infos.GetInfo()
				// 2. get cli
				cli := clients.GetClient()
				// 3. get acc number
				acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
				if err != nil {
					log.Println(err)
					continue
				}
				accNum, seqNum := acc.GetAccountNumber(), acc.GetSequence()
				offset := uint64(0)
				ethAddr, _ := utils.ToHexAddress(info.GetAddress().String())

				// Let Us GO GO GO !!!!!!
				// 1. add liquididy
				toEthAddr, _ :=  utils.ToHexAddress(info.GetAddress().String())
				payload := UniswapV2.BuildAddLiquidOKTPayload(
					tokenAddr, toEthAddr.String(),
					sdk.MustNewDecFromStr("0.1").Int, sdk.MustNewDecFromStr("0.0001").Int, sdk.MustNewDecFromStr("0.0001").Int,
					time.Now().Add(time.Hour*24).Unix(),
				)
				res, err := cli.Evm().SendTx(info, common.PassWord, routerAddr, "0.001", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
				if err != nil {
					log.Printf("(%d)[%s] %s failed to add liquidity in %s: %s\n", fail.Add(), res.TxHash, ethAddr, routerAddr, err)
					continue
				} else {
					log.Printf("(%d)[%s] %s add liquidity in %s \n", succ.Add(), res.TxHash, ethAddr, routerAddr)
					offset++
				}

				// 2.1 stake
				payload = UniswapV2Staker.BuildStakePayload(1500000000000000)
				res, err = cli.Evm().SendTx(info, common.PassWord, poolAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum+offset)
				if err != nil {
					log.Printf("(%d)[%s] %s failed to stake lp in %s: %s\n", fail.Add(), res.TxHash, ethAddr, poolAddr, err)
					continue
				} else {
					log.Printf("(%d)[%s] %s stake lp in %s \n", succ.Add(), res.TxHash, ethAddr, poolAddr)
					offset++
				}

				// 2.2 withDraw randomly
				payload = UniswapV2Staker.BuildWithdrawPayload(500000000)
				res, err = cli.Evm().SendTx(info, common.PassWord, poolAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum+offset)
				if err != nil {
					log.Printf("(%d)[%s] %s failed to withdraw lp from %s: %s\n", fail.Add(), res.TxHash, ethAddr, poolAddr, err)
					continue
				} else {
					log.Printf("(%d)[%s] %s withdraw lp from %s \n", succ.Add(), res.TxHash, ethAddr, poolAddr)
					offset++
				}

				// 2.3 get Reward randomly
				rand.Seed(time.Now().Unix())
				if rand.Intn(10) <= 3 {
					payload = UniswapV2Staker.BuildGetRewardPayload()
					res, err = cli.Evm().SendTx(info, common.PassWord, poolAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum+offset)
					if err != nil {
						log.Printf("(%d)[%s] %s failed get reward from %s: %s\n", fail.Add(), res.TxHash, ethAddr, poolAddr, err)
						continue
					} else {
						log.Printf("(%d)[%s] %s get reward from %s \n", succ.Add(), res.TxHash, ethAddr, poolAddr)
						offset++
					}
				}

				// 2.4 Exit randomly
				rand.Seed(time.Now().Unix())
				if rand.Intn(10) <= 2 {
					payload = UniswapV2Staker.BuildExitPayload()
					res, err = cli.Evm().SendTx(info, common.PassWord, poolAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum+offset)
					if err != nil {
						log.Printf("(%d)[%s] %s failed to exit all lp from %s: %s\n", fail.Add(), res.TxHash, ethAddr, poolAddr, err)
						continue
					} else {
						log.Printf("(%d)[%s] %s exit all lp from %s \n", succ.Add(), res.TxHash, ethAddr, poolAddr)
						offset++
					}
				}

				time.Sleep(time.Second*time.Duration(sleepTime))
			}
		}()
	}
	wg.Wait()
}
