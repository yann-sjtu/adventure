package simple

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/mitchellh/mapstructure"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/simple/dex"
	"github.com/okex/adventure/x/simple/distribution"
	"github.com/okex/adventure/x/simple/staking"
	"github.com/okex/adventure/x/simple/token"
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
	cfg := common.GetConfig()
	cfg.TestCaesPath = txConfigPath
	fmt.Println(cfg)

	//init test cases
	cases, err := common.ReadTestCases(txConfigPath)
	if err != nil {
		return err
	}
	fmt.Println(cases)

	//run test cases
	var wg sync.WaitGroup
	for _, c := range cases {
		wg.Add(1)
		go func(c common.TestCase) {
			defer wg.Done()
			switch c.RunTxMode {
			case "", common.Parallel:
				excuteTxsInParallel(c, cfg.Hosts)
			case common.Serial:
				log.Fatalln("the serial mode is still in developing")
			default:
				log.Fatalf("not support of the '%s' mode", c.RunTxMode)
			}
		}(c)
	}

	wg.Wait()
	return nil
}

func excuteTxsInParallel(c common.TestCase, hosts []string) {
	manager := common.GetAccountManagerFromFile(c.MnemonicPath)
	for _, tx := range c.Transactions {
		var arg common.BaseParam
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
			handler = token.Issue
		case common.Mint:
			handler = token.Mint
		case common.Burn:
			handler = token.Burn
		case common.Edit:
			handler = token.Edit
		case common.MultiSend:
			handler = token.MultiSend
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
		case common.DelegateVoteUnbond:
			handler = staking.DelegateVoteUnbond
		case common.Proxy:
			handler = staking.Proxy
		default:
			fmt.Printf("the types '%s' of tx is not supported now\n", tx.Type)
		}
		clientManager := common.NewClientManager(hosts, arg.Fee)
		go executeTxInLoop(manager, clientManager, arg, handler)
	}

	// until the world die away
	select {}
}

func executeTxInLoop(m *common.AccountManager, c *common.ClientManager, p common.BaseParam, handler func(*gosdk.Client, keys.Info)) {
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

func excuteTxsInSerial(accountPath string, baseParam common.BaseParam, transactions []common.Transaction) {

}
