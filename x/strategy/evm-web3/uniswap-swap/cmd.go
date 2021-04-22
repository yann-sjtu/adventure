package uniswap_swap

import (
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/tools"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	goroutineNum int
	privkeyPath  string
	sleepTime    int
	deposit      bool
	mode         string
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uniswap-swap",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.IntVarP(&goroutineNum, "goroutine-num", "g", 1, "set Goroutine Num")
	flags.StringVarP(&privkeyPath, "private-path", "p", "", "set the Priv Key path")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 0, "set the sleep time")
	flags.BoolVarP(&deposit, "deposit", "d", true, "deposit okt or not")
	flags.StringVarP(&mode, "mode", "s", "sync", "set the mode of sync or block")
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
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkeyPath)
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.005okt", mode, 500000)
	succ, fail := tools.NewCounter(-1), tools.NewCounter(-1)

	var wg sync.WaitGroup
	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			privkey := privkeys[index]
			cli := clients.GetClient()
			info, err := utils.CreateAccountWithPrivateKey(privkey, "acc", common.PassWord)
			if err != nil {
				panic(err)
			}

			accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
			if err != nil {
				panic(err)
			}
			seqNum := accInfo.GetSequence()
			offset := uint64(0)

			ethAddr, err  := utils.ToHexAddress(info.GetAddress().String())
			if err != nil {
				panic(err)
			}
			//fmt.Println(privkey)
			//fmt.Println(info.GetAddress().String())
			//fmt.Println(ethAddr)

			//init various payload
			depositPayloadStr := hexutil.Encode(UniswapV2.BuildWethDepositPayload())
			approvePayloadStr := hexutil.Encode(UniswapV2.BuildWethApprovePayload(
				routerAddr, 9999999999,
			))
			swapPayloadStr := hexutil.Encode(UniswapV2.BuildSwapExactTokensForTokensPayload(
				big.NewInt(1000), big.NewInt(0),
				[]string{WethAddr,UsdtAddr}, "0x2B5Cf24AeBcE90f0B8f80Bc42603157b27cFbf47",
				time.Now().Add(time.Hour*8640).Unix(),
			))

			// deposit weth
			if deposit {
				res, err := cli.Evm().SendTxEthereum(privkey, WethAddr, "1",depositPayloadStr, 500000, seqNum+offset)
				if err != nil {
					panic(fmt.Errorf("[TxHash: %s] %s failed to deposit %sokt on %s: %s", res.TxHash, ethAddr, sdk.NewDec(1), WethAddr, err))
				}
				log.Printf("[TxHash: %s] %s deposit %sokt on %s \n", res.TxHash, ethAddr, sdk.NewDec(1), WethAddr)

				offset++
			}

			// approve tx
			res, err := cli.Evm().SendTxEthereum(privkey, WethAddr, "", approvePayloadStr, 500000, seqNum+offset)
			if err != nil {
				panic(fmt.Errorf("[TxHash: %s] %s failed to approve 9999999999coin from %s to %s: %s", res.TxHash, ethAddr, WethAddr, routerAddr, err))
			}
			log.Printf("[TxHash: %s] %s approve 9999999999coin from %s to %s\n", res.TxHash, ethAddr, WethAddr, routerAddr)
			offset++

			for {
				res, err = cli.Evm().SendTxEthereum(privkey, routerAddr, "", swapPayloadStr, 500000, seqNum+offset)
				if err != nil {
					log.Printf("(%d)[TxHash: %s] %s failed to swap weth for usdt: %s\n", fail.Add(), res.TxHash, ethAddr, err)
					continue
				}
				log.Printf("(%d)[TxHash: %s] %s swap 1000*10^-18 weth %s for usdt %s\n", succ.Add(), res.TxHash, ethAddr, WethAddr, UsdtAddr)
				offset++
				time.Sleep(time.Duration(sleepTime)*time.Second)
			}

		}(i)
	}
	wg.Wait()
}