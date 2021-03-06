package uniswap

import (
	"sync"

	deploy_contracts "github.com/okex/adventure/x/strategy/evm/deploy-contracts"
	"github.com/spf13/cobra"
)

var (
	goroutineNum int
	privkey      string
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
	flags.StringVarP(&privkey, "private-path", "p", "", "set the Priv Key path")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 0, "set the sleep time")

	return cmd
}

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
	var wg sync.WaitGroup
	for i := 0; i < deploy_contracts.GoroutineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {

			}
		}
	}
	wg.Done()
}
