package farm_rm_liquidity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	common "github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func FarmRemoveLiquidityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-liquidity",
		Short: "remove liquidity",
		Args:  cobra.NoArgs,
		RunE:  runFarmRemoveLiquidityCmd,
	}
}

func runFarmRemoveLiquidityCmd(cmd *cobra.Command, args []string) error {
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	// build filters
	filters := make(map[string]int)
	for i, addr := range addrs {
		filters[addr] = i + workerStartIndex
	}

	fmt.Printf("%d addresses are loaded\n", len(filters))
	var counter int
	for {
		counter++
		fmt.Printf("================================ Round %d ================================\n", counter)
		var targetAddr string
		for addr := range filters {
			cli := clientManager.GetClient()
			_, err := cli.Farm().QueryLockInfo(poolName, addr)
			if err == nil {
				continue
			} else if !strings.Contains(err.Error(), "hasn't locked") {
				log.Printf("[%s] query lock info error: %s\n", addr, err.Error())
				continue
			}

			// rm liquidity
			accInfo, err := cli.Auth().QueryAccount(addr)
			if err != nil {
				log.Printf("[%s] query lock info error: %s\n", addr, err.Error())
				continue
			}

			lptAmount := accInfo.GetCoins().AmountOf(lockSymbol)
			msgRmLiquidity := newMsgRemoveLiquidity(accInfo.GetAccountNumber(), accInfo.GetSequence(),
				lptAmount, sdk.NewDecCoinFromDec(baseCoin, sdk.ZeroDec()), sdk.NewDecCoinFromDec(quoteCoin, sdk.ZeroDec()), getDeadline(),
				addr)

			index, ok := filters[addr]
			if !ok {
				panic("index of worker does not exist, " + addr)
			}

			if err = mntcmn.SendMsg(mntcmn.Undelefarm, msgRmLiquidity, index); err != nil {
				fmt.Println("failed to mntcmn.SendMsg():", err)
				continue
			}

			// assume success
			log.Printf("%s info to remove liquidity successfully\n", addr)
			targetAddr = addr
			break
		}

		delete(filters, targetAddr)
		if len(filters) == 0 {
			break
		} else {
			fmt.Printf("%d workers left\n", len(filters))
		}
	}

	fmt.Println("MISSION SUCCESS")
	return nil
}
