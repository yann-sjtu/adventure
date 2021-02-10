package evm

import (
	"log"
	"sync"

	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

func uniswapTestCmd() *cobra.Command {
	InitTemplate()

	cmd := &cobra.Command{
		Use:   "uniswap-testnet-operate",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.IntVarP(&GoroutineNum, "goroutine-num", "g", 1, "set Goroutine Num of deploying contracts")
	flags.StringVarP(&MnemonicPath, "mnemonic-path", "m", "", "set the MnemonicPath path")

	return cmd
}

const (
	usdtAddr = "0xffea71957a3101d14474a3c358ede310e17c2409"
)

func testLoop(cmd *cobra.Command, args []string) {
	infos := common.GetAccountManagerFromFile(MnemonicPath)
	clients := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	//contractManager := tools.NewContractManager()
	var wg sync.WaitGroup
	for i := 0; i < GoroutineNum; i++ {
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


			}
		}()
	}
	wg.Wait()
}


