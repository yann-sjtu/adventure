package deploy_contracts

import (
	"log"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/ERC721"
	"github.com/okex/adventure/x/strategy/evm/template/USDT"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/tools"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func DeployErc20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-erc20-tokens",
		Short: "arbitrage token from swap and orderdepthbook",
		Args:  cobra.NoArgs,
		Run:   deployErc20Tokens,
	}

	flags := cmd.Flags()
	flags.IntVarP(&Num, "num", "n", 1000, "set Num of issusing token")
	flags.IntVarP(&GoroutineNum, "goroutine-num", "g", 1, "set Goroutine Num of deploying contracts")
	flags.IntVarP(&TransferGoNum, "transfer-go-num", "t", 1, "set Goroutine Num of transfering erc20 token")
	flags.StringVarP(&MnemonicPath, "mnemonic-path", "m", "", "set the MnemonicPath path")

	return cmd
}

var (
	Num           = 1000
	GoroutineNum  = 1
	TransferGoNum = 1

	MnemonicPath = ""
)

func deployErc20Tokens(cmd *cobra.Command, args []string) {
	infos := common.GetAccountManagerFromFile(MnemonicPath)
	clients := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	contractManager := tools.NewContractManager()
	counter := tools.NewCounter(Num)
	var wg sync.WaitGroup
	for i := 0; i < GoroutineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if counter.IsOver() {
					break
				}

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

				// 4.1 deploy factory
				ethAddress := utils.EthAddress(utils.GetEthAddressStrFromCosmosAddr(info.GetAddress()))
				facPayload := UniswapV2.BuildFactoryContractPayload(ethAddress)
				_, facAddress, err := cli.Evm().CreateContract(info, common.PassWord, "", ethcommon.Bytes2Hex(facPayload), "", accNum, seqNum)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("[%d]uniswapv2.factory contract addr: %s\n", counter.Add(), facAddress)

				// 4.2 deploy weth
				wethPayload := UniswapV2.BuildWethContractPayload()
				_, wethAddress, err := cli.Evm().CreateContract(info, common.PassWord, "", ethcommon.Bytes2Hex(wethPayload), "", accNum, seqNum+1)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("[%d]uniswapv2.weth contract addr: %s\n", counter.Add(), wethAddress)

				// 4.3 deploy router
				routerPayload := UniswapV2.BuildRouterContractPayload(utils.EthAddress(facAddress), utils.EthAddress(wethAddress))
				_, routerAddress, err := cli.Evm().CreateContract(info, common.PassWord, "", ethcommon.Bytes2Hex(routerPayload), "", accNum, seqNum+2)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("[%d]uniswapv2.router contract addr: %s\n", counter.Add(), routerAddress)

				// 4.4 deploy erc721
				ERC721Payload := ERC721.BuildERC721ContractPayload("okexchain coin", "OKB")
				_, ERC721Address, err := cli.Evm().CreateContract(info, common.PassWord, "", ethcommon.Bytes2Hex(ERC721Payload), "", accNum, seqNum+3)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("[%d]erc721 contract addr: %s\n", counter.Add(), ERC721Address)

				// 4.5 deploy erc20 usdt
				USDTPayload := USDT.BuildUSDTContractPayload(utils.Uint256(12642013521397079), utils.Uint256(6), "OKEX USD", "TUSDT")
				_, USDTAddress, err := cli.Evm().CreateContract(info, common.PassWord, "", ethcommon.Bytes2Hex(USDTPayload), "", accNum, seqNum+4)
				if err != nil {
					log.Println(err)
					continue
				}
				contractManager.AppendUSDTAddr(USDTAddress, info)
				log.Printf("[%d]erc20 usdt contract addr: %s\n", counter.Add(), USDTAddress)
			}
		}()
	}
	wg.Wait()
	log.Println("total number of contracts being deployed:", counter.GetCurrentNum())

	log.Println("wait half a minute, start testing transfer")
	time.Sleep(time.Second * 30)
	transfer(clients, contractManager)
}

func getTmpClient() gosdk.Client {
	cfg, _ := types.NewClientConfig("http://localhost:26657", common.Cfg.ChainId, types.BroadcastBlock, "", 20000000, 1.5, "0.00000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)
	return cli
}
