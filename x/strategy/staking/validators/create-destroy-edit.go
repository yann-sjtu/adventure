package validators

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/okex/adventure/x/strategy/staking/validators/types"
	gotypes "github.com/okex/exchain-go-sdk/module/staking/types"
	"github.com/spf13/cobra"
)

func valsLoopTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "val-loop",
		Short: "create -> edit -> destroy validators circularly",
		Args:  cobra.NoArgs,
		RunE:  runValsLoop,
	}

	return cmd
}

const (
	maxToDestroy = 3
	minIncumbent = 35

	// time to wait when the number of incumbent val is less than minIncumbent
	timeWaitingInsufficientIncumbent = time.Minute * 30

	// round interval
	timeInterval = time.Minute * 60
)

func runValsLoop(cmd *cobra.Command, args []string) error {

	// get valManager
	valManager := types.GetValManager()

	var round int

	for {
		round++
		logging(fmt.Sprintf("Round %d", round), nil)
		log.Printf("[create success] %d 				[destroy success] %d",
			types.CreateSuccessCnt, types.DestroySuccessCnt)

		// do nothing if the number of healthy validators is less than 35
		// reason: votes are not fully withdrawn from the destroyed validators
		icbValAddrs, err := valManager.GetIncumbentValAddrs()
		if err != nil {
			log.Println(err)
		}

		var wg sync.WaitGroup
		// TX: edit-validator
		logging("EDIT validators", icbValAddrs)

		for _, valAddr := range icbValAddrs {
			wg.Add(1)
			go valManager[valAddr].Edit(&wg)
		}
		wg.Wait()

		// TX: create validator
		vals, err := valManager.GetValidators()
		if err != nil {
			log.Println(err)
		}

		rebirthValAddrs := pickRebirth(vals, valManager)
		logging("CREATE validators", rebirthValAddrs)

		for _, valAddr := range rebirthValAddrs {
			wg.Add(1)
			go valManager[valAddr].Create(&wg)
		}
		wg.Wait()

		if len(icbValAddrs) <= minIncumbent {
			log.Printf("==== WARNING: less incumbent validator: %d val ====\n", len(icbValAddrs))
			time.Sleep(timeWaitingInsufficientIncumbent)
			continue
		}

		// TX: destroy validator randomly
		rand.Seed(time.Now().UnixNano())
		numToDestroy := rand.Intn(maxToDestroy)
		quitValAddrs := pickQuitVals(icbValAddrs, numToDestroy)
		logging("DESTROY validators", quitValAddrs)

		for _, valAddr := range quitValAddrs {
			wg.Add(1)
			go valManager[valAddr].Destroy(&wg)
		}
		wg.Wait()

		logging("INTERVAL", []string{})
		time.Sleep(timeInterval)
	}

}

// log
func logging(title string, valAddrs []string) {
	log.Printf("==================== %s ====================\n", title)
	for _, addr := range valAddrs {
		log.Println(addr)
	}
}

// pickDiffSet returns the slice of string val operator address that are removed
func pickRebirth(curVals []gotypes.Validator, manager types.ValManager) (vals []string) {
	// get a new filter(check which validator has been removed)
	filter := newFilter(manager)
	for _, val := range curVals {
		delete(filter, val.OperatorAddress.String())
	}

	for addr := range filter {
		vals = append(vals, addr)
	}

	return
}

// getFilter returns a filter to check which validator has been remove
func newFilter(manager types.ValManager) map[string]struct{} {
	filter := make(map[string]struct{}, len(manager))
	for k := range manager {
		filter[k] = struct{}{}
	}

	return filter
}

// pickQuitVals returns the slice of string val operator address with a specific number to destroy
func pickQuitVals(allValAddrs []string, num int) []string {
	addrs := shuffle(allValAddrs)
	valAddrs := make([]string, num)
	//for i, j := 0, 0; j < num; i++ {
	//	if allValAddrs[i] == fixValOperAddr {
	//		continue
	//	}
	//	valAddrs[j] = allValAddrs[i]
	//	j++
	//}
	for i := 0; i < num; i++ {
		valAddrs[i] = addrs[i]
	}

	return valAddrs
}

func shuffle(vals []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]string, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}
