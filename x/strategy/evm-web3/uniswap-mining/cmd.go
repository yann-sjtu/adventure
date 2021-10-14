package uniswap_mining

import (
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2Staker"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/spf13/cobra"
)

var (
	goroutineNum int
	privkeyPath  string
	sleepTime    int
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uniswap-testnet-operate",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.IntVarP(&goroutineNum, "goroutine-num", "g", 1, "set Goroutine Num")
	flags.StringVarP(&privkeyPath, "private-path", "p", "", "set the Priv Key path")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 0, "set the sleep time")

	return cmd
}

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
	WethAddr = "0x70c1c53E991F31981d592C2d865383AC0d212225"
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

	m := common.GetPrivKeyManager(privkeyPath)
	leng := m.Length()
	clients := common.NewClientManagerWithMode(common.GlobalConfig.Networks[common.NetworkType].Hosts, "0.005okt", types.BroadcastSync, 500000)

	depositPayloadStr := hexutil.Encode(UniswapV2.BuildWethDepositPayload())
	approveToRouterPayloadStr := hexutil.Encode(UniswapV2.BuildWethApprovePayload(routerAddr, 10))
	approveToPoolPayloadStr := hexutil.Encode(UniswapV2.BuildWethApprovePayload(OktUsdtPoolAddr, 10))

	stakePayloadStr := hexutil.Encode(UniswapV2Staker.BuildStakePayload(100))
	getRewardPayload := hexutil.Encode(UniswapV2Staker.BuildGetRewardPayload())
	withdrawPayload := hexutil.Encode(UniswapV2Staker.BuildWithdrawPayload(10))
	exitPayload := hexutil.Encode(UniswapV2Staker.BuildExitPayload())

	var wg sync.WaitGroup
	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go func(index int, cli *gosdk.Client) {
			defer wg.Done()
			for k := 0; ; k++ {
				accinfo := m.GetAccount((k*goroutineNum+index)%leng)
				ethAddr, privkey := accinfo.GetEthAddress().String(), accinfo.GetPirvkey()

				// 0.1 deposit okt
				res, err := cli.Evm().SendTxEthereum(privkey, WethAddr, "0.000000001", depositPayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s deposit       failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s deposit       done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 0.2 approve wokt
				res, err = cli.Evm().SendTxEthereum(privkey, WethAddr, "", approveToRouterPayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s approve wokt  failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s approve wokt  done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 0.3 swap wokt -> usdt
				swapPayloadStr := hexutil.Encode(UniswapV2.BuildSwapExactTokensForTokensPayload(big.NewInt(10000), big.NewInt(0), []string{WethAddr,UsdtAddr}, ethAddr, time.Now().Add(time.Hour*8640).Unix()))
				res, err = cli.Evm().SendTxEthereum(privkey, routerAddr, "", swapPayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s swap          failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s swap          done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 1.1 approve usdt
				res, err = cli.Evm().SendTxEthereum(privkey, tokenAddr, "", approveToRouterPayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s approve usdt  failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s approve usdt  done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 1.2 add liquidity (usdt-okt)
				addLiquidPayloadStr := hexutil.Encode(UniswapV2.BuildAddLiquidOKTPayload(tokenAddr, ethAddr, sdk.MustNewDecFromStr("0.0000000000001").Int, sdk.MustNewDecFromStr("0").Int, sdk.MustNewDecFromStr("0").Int, time.Now().Add(time.Hour*8640).Unix()))
				res, err = cli.Evm().SendTxEthereum(privkey, routerAddr, "0.000000001", addLiquidPayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s addLiquidity  failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s addLiquidity  done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 2.1 approve uni
				res, err = cli.Evm().SendTxEthereum(privkey, OktUsdtLPAddr, "", approveToPoolPayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s approve uni   failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s approve uni   done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 2.2 stake
				res, err = cli.Evm().SendTxEthereum(privkey, poolAddr, "", stakePayloadStr, 500000, accinfo.GetNonce(cli))
				if err != nil {
					log.Printf("[%s] %s stake         failed: %s\n", res.TxHash, ethAddr, err)
					continue
				} else {
					log.Printf("[%s] %s stake         done\n", res.TxHash, ethAddr)
					accinfo.AddNonce()
				}

				// 2.3 withDraw randomly
				rand.Seed(time.Now().UnixNano())
				if rand.Intn(10) <= 3 {
					res, err = cli.Evm().SendTxEthereum(privkey, poolAddr, "", withdrawPayload, 500000, accinfo.GetNonce(cli))
					if err != nil {
						log.Printf("[%s] %s withdraw      fail: %s\n", res.TxHash, ethAddr, err)
						continue
					} else {
						log.Printf("[%s] %s withdraw      done\n", res.TxHash, ethAddr)
						accinfo.AddNonce()
					}
				}
				// 2.4 get Reward randomly
				if rand.Intn(10) <= 3 {
					res, err = cli.Evm().SendTxEthereum(privkey, poolAddr, "", getRewardPayload, 500000, accinfo.GetNonce(cli))
					if err != nil {
						log.Printf("[%s] %s getReward     fail: %s\n", res.TxHash, ethAddr, err)
						continue
					} else {
						log.Printf("[%s] %s getReward     done\n", res.TxHash, ethAddr)
						accinfo.AddNonce()
					}
				}
				// 2.5 Exit randomly
				if rand.Intn(10) <= 3 {
					res, err = cli.Evm().SendTxEthereum(privkey, poolAddr, "", exitPayload, 500000, accinfo.GetNonce(cli))
					if err != nil {
						log.Printf("[%s] %s exit          fail: %s\n", res.TxHash, ethAddr, err)
						continue
					} else {
						log.Printf("[%s] %s exit          done\n", res.TxHash, ethAddr)
						accinfo.AddNonce()
					}
				}
				time.Sleep(time.Duration(sleepTime) * time.Second)
			}
		}(i, clients.GetClient())
	}
	wg.Wait()
}
