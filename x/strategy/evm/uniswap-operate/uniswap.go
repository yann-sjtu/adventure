package uniswap_operate

import (
	"log"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/deploy-contracts"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
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

	return cmd
}

var (
	LPAddrs    = [4]string{OktUsdtLPAddr, OktDotkLPAddr, OktBtckLPAddr, OktEthkLPAddr}
	PoolAddrs  = [4]string{OktUsdtPoolAddr, OktDotkPoolAddr, OktBtckPoolAddr, OktEthkPoolAddr}
	TokenAddrs = [4]string{UsdtAddr, DotkAddr, BtckAddr, EthkAddr}
)

const (
	routerAddr = "0x0653a68B22b18663F69a7103621F7F3EB59191F1"

	OktUsdtLPAddr   = "0x7068B191ff97e32D6Fbba3204408877A9007BBd1"
	OktUsdtPoolAddr = "0x0Bd475f8b27EA57158291372667aD1e7eeD5C174"

	OktDotkLPAddr   = "0x4425e0ca22949f8ed75102283bb83e15180bc579"
	OktDotkPoolAddr = "0x87Bf79788c9dBa2298e172E25fBBBE20ABD673A4"

	OktBtckLPAddr   = "0xef7a562bbb388766551bb15138d7cc1c7e4c85d8"
	OktBtckPoolAddr = "0x1f714Ab64C9B2A86fFD1E2A7Ed5e9a0515A5688A"

	OktEthkLPAddr   = "0xe76e553eab658f0eef2ea8d29b9719b9875786e6"
	OktEthkPoolAddr = "0xb8F6fb2A7A469d21A993421bD397Dfe213c4E6AD"

	UniUsdtLPAddr   = "0xfc56c01730f1d47cd187253353521d3dc2218a82"
	UniUsdtPoolAddr = "0xaAFd4b09e0c275b3EC35B3cacB99D6DA9Ca96E33"

	UsdtAddr = "0xffea71957a3101d14474a3c358ede310e17c2409"
	WoktAddr = "0x2789Fdc29D0f1D2ddaC362B2cb79F7799A5fbdAF"
	UniAddr  = "0x0A1D36fCD446Df6bA0bA326bec5291417B97757d"
	OkbAddr  = "0xa860E9929B7DE53218c9B0a555680587D3542882"
	EthkAddr = "0x8d06747eD6EEc4b4dB9a365EF55a4aE4A0Cc0c27"
	BtckAddr = "0xd9A7425dCD77DF8A388E404b38E0D539B4d2D742"
	UsdcAddr = "0x7B334746E0B9f7fbD94AD9f4eA9e304e1d2dF0DA"
	FilkAddr = "0x33c548B01c04D195E99c16C6dC1D4E9252EE45ea"
	DotkAddr = "0x6Cf49E69a54C42cFf0f8bdD7fEf75cfB0cF965Ac"
	LtckAddr = "0xA51E71874112cd7fa7885C23D403525Ee0F73c80"
	UsdkAddr = "0xcBCc53b501A799Dd90D6546aa5319cF87a3E66fa"
)

func testLoop(cmd *cobra.Command, args []string) {
	//lpAddr, poolAddr, tokenAddr := LPAddrs[0], PoolAddrs[0], TokenAddrs[0]

	infos := common.GetAccountManagerFromFile(deploy_contracts.MnemonicPath)
	clients := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

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

				// Let Us GO GO GO !!!!!!
				// 1. add liquididy
				payload := UniswapV2.BuildAddLiquidOKTPayload(
					"tokenAddr", utils.GetEthAddressStrFromCosmosAddr(info.GetAddress()),
					6000000000000000000,1,1,
					int(time.Now().Add(time.Hour*24).Unix()),
					)
				res, err := cli.Evm().SendTx(info, common.PassWord, routerAddr, "1", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("[%s] %s add liquidity in %s \n", res.TxHash, )
				}

				time.Sleep(time.Second*5)
				// get acc number again
				acc, err = cli.Auth().QueryAccount(info.GetAddress().String())
				if err != nil {
					log.Println(err)
					continue
				}
				accNum, seqNum = acc.GetAccountNumber(), acc.GetSequence()
			}
		}()
	}
	wg.Wait()
}
