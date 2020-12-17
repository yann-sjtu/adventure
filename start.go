package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/mitchellh/mapstructure"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/okex/adventure/common/logger"
	ammswap2 "github.com/okex/adventure/x/simple/ammswap"
	"github.com/okex/adventure/x/simple/dex"
	"github.com/okex/adventure/x/simple/distribution"
	order2 "github.com/okex/adventure/x/simple/order"
	staking2 "github.com/okex/adventure/x/simple/staking"
	token2 "github.com/okex/adventure/x/simple/token"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/spf13/cobra"
)

var (
	txConfigPath = ""
)

// TxCmd return the start cmd
func TxCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "transmit all transactions defined in tx json file",
		RunE:  RunStart,
	}

	startCmd.PersistentFlags().StringVarP(&txConfigPath, "config_path", "p", "", "the tx json config path")
	return startCmd
}

func RunStart(cmd *cobra.Command, args []string) error {
	//init config
	cfg := config.GetConfig()
	if txConfigPath != "" {
		cfg.TxConfigPath = txConfigPath
	}
	fmt.Println(cfg)
	//init logger
	logger.InitLogger(cfg.LogLevel, cfg.LogListenUrl)
	//init test cases
	cases, err := config.ReadTestCases(cfg.TxConfigPath)
	if err != nil {
		return err
	}
	fmt.Println(cases)

	//run test cases
	var wg sync.WaitGroup
	for _, c := range cases {
		wg.Add(1)
		go func(c config.TestCase) {
			defer wg.Done()
			switch c.RunTxMode {
			case "", config.Parallel:
				excuteTxsInParallel(c, cfg.Hosts)
			case config.Serial:
				log.Fatalln("the serial mode is still in developing")
			default:
				log.Fatalf("not support of the '%s' mode", c.RunTxMode)
			}
		}(c)
	}

	wg.Wait()
	return nil
}

func excuteTxsInParallel(c config.TestCase, hosts []string) {
	manager := common.GetAccountManagerFromFile(c.MnemonicPath)
	for _, tx := range c.Transactions {
		var arg config.BaseParam
		err := mapstructure.Decode(tx.Args, &arg)
		if err != nil {
			log.Fatalf("failed to decode args config. error: %s\n", err.Error())
		}
		arg.SetBaseParam(c.BaseParam)

		var handler func(*gosdk.Client, keys.Info)
		switch tx.Type {
		case common.SetWithdrawAddr:
			handler = distribution.SetWithdrawAddr
		case common.WithdrawRewards:
			handler = distribution.WithdrawRewards
		case common.Issue:
			handler = token2.Issue
		case common.Mint:
			handler = token2.Mint
		case common.Burn:
			handler = token2.Burn
		case common.Edit:
			handler = token2.Edit
		case common.MultiSend:
			handler = token2.MultiSend
		case common.Deposit:
			handler = dex.Deposit
		case common.Withdraw:
			handler = dex.Withdraw
		case common.List:
			handler = dex.List
		case common.RegisterOperator:
			handler = dex.RegisterOperator
		case common.EditOperator:
			handler = dex.EditOperator
		case common.Order:
			handler = order2.Orders
		case common.DelegateVoteUnbond:
			handler = staking2.DelegateVoteUnbond
		case common.Proxy:
			handler = staking2.Proxy
		case common.AddLiquidity:
			handler = ammswap2.AddLiquidity
		case common.RemoveLiquidity:
			handler = ammswap2.RemoveLiquidity
		case common.CreateExchange:
			handler = ammswap2.CreateExchange
		case common.SwapExchange:
			handler = ammswap2.SwapExchange
		default:
			fmt.Printf("the types '%s' of tx is not supported now\n", tx.Type)
		}
		clientManager := common.NewClientManager(hosts, arg.Fee)
		go executeTxInLoop(manager, clientManager, arg, handler)
	}

	// until the world die away
	select {}
}

func executeTxInLoop(m *common.AccountManager, c *common.ClientManager, p config.BaseParam, handler func(*gosdk.Client, keys.Info)) {
	totalRound := p.RoundNum
	concurrentNum := p.ConcurrentNum
	var wg sync.WaitGroup
	wg.Add(concurrentNum)

	for n := 0; n < concurrentNum; n++ {
		go func() {
			defer wg.Done()
			for round := 0; round < totalRound; round++ {
				info := m.GetInfo()
				cli := c.GetClient()
				handler(cli, info)

				time.Sleep(time.Duration(p.SleepTime) * time.Second)
			}
		}()
	}
	wg.Wait()
}

func excuteTxsInSerial(accountPath string, baseParam config.BaseParam, transactions []config.Transaction) {

}
