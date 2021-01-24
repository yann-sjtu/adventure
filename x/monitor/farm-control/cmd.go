package farm_control

import (
	"fmt"
	"log"
	"time"

	"github.com/okex/adventure/common"
	"github.com/spf13/cobra"
)

func FarmControlCmd() *cobra.Command {
	farmControlCmd := &cobra.Command{
		Use:   "farm-control",
		Short: "farm controlling in the first miner on a farm pool",
		Args:  cobra.NoArgs,
		RunE:  runFarmControlCmd,
	}

	//flags := farmControlCmd.Flags()

	return farmControlCmd
}

const (
	poolName   = "1st_pool_okt_usdt"
	lockSymbol = "ammswap_okt_usdt-a2b"

	baseCoin  = "okt"
	quoteCoin = "usdt-a2b"
)

func runFarmControlCmd(cmd *cobra.Command, args []string) error {
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	if err := refreshFarmAccounts(clientManager.GetClient()); err != nil {
		return err
	}

	for i := 0; ; i++ {
		// 0. sleep 60 seconds, or so
		time.Sleep(time.Second * 120)
		log.Printf("\n======================== Round %d ========================\n", i)
		cli := clientManager.GetClient()
		if i%10 == 0 && i != 0  { // todo: used for refreshing accounts cache storged in local, this judgement might be removed
			time.Sleep(time.Second * 60)
			for j := 0; j < 10; j++ {
				if 	err := refreshFarmAccounts(cli); err != nil {
					fmt.Printf("[Phase0 Refresh %d] failed: %s\n", j, err.Error())
					continue
				}
				break
			}
		}

		// 1. check the ratio of (our_total_locked_lpt / total_locked_lpt), then return how many lpt to be replenished
		requiredToken, err := calculateReuiredAmount(cli)
		if err != nil {
			fmt.Printf("[Phase1 Calculate] failed: %s\n", err.Error())
			continue
		}

		// 2. judge if the requiredToken is zero or not
		if requiredToken.IsZero() {
			// 2.1 our_total_locked_lpt / total_locked_lpt > 80%, then do nothing
			fmt.Printf("This Round doesn't need to lock more %s \n", lockSymbol)
		} else {
			// 2.1 our_total_locked_lpt / total_locked_lpt < 80%, then promote the ratio over 81%
			replenishLockedToken(cli, requiredToken)
		}
	}
}