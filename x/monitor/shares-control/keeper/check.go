package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/x/monitor/shares-control/utils"
	gosdk "github.com/okex/okexchain-go-sdk"
)

func (k *Keeper) checkValNumInTop21(vals []gosdk.Validator) (warning bool, valAddrsToPromote []string) {
	// check which target validator is out of top 21
	filter := utils.NewFilter(k.targetValAddrs)
	for _, val := range vals {
		delete(filter, val.OperatorAddress.String())
	}

	if len(filter) == 0 {
		return
	}

	fmt.Println(" WARNING!!!!! targets validators missed:")
	for addr := range filter {
		valAddrsToPromote = append(valAddrsToPromote, addr)
		fmt.Printf("\t\t%s\n", addr)
	}

	return true, valAddrsToPromote
}

func (k *Keeper) checkPercentToDominate(targetTotal, bonedTotal sdk.Dec) bool {
	curPercentToDominate := targetTotal.Quo(bonedTotal)
	fmt.Printf(`           Percentage to dominate
	current: [%s]    expected: [%s]
`,
		curPercentToDominate,
		k.expectedParams.GetExpectedPercentToDominate(),
	)

	isDangerous := curPercentToDominate.LT(k.expectedParams.GetExpectedPercentToDominate())
	if isDangerous {
		fmt.Println(" WARNING!!!!! current DOMINATION PERCENTAGE is less than THE EXPECTED")
	}

	return isDangerous
}

func (k *Keeper) checkPercentToPlunder(vals []gosdk.Validator, targetTotal, globalTotal sdk.Dec) bool {
	curPercentToPlunder := k.calculatePercentToPlunder(vals, targetTotal, globalTotal)
	fmt.Printf(`           Percentage to plunder
	current: [%s]    expected: [%s]
`,
		curPercentToPlunder,
		k.expectedParams.GetExpectedPercentToPlunder(),
	)

	isDangerous := curPercentToPlunder.LT(k.expectedParams.GetExpectedPercentToPlunder())
	if isDangerous {
		fmt.Println(" WARNING!!!!! current REWARD PERCENTAGE is less than THE EXPECTED")
	}

	return isDangerous
}
