package keeper

import "fmt"

func (k *Keeper) logInit() {
	fmt.Printf(`
============================================================
               %d target validators    	
============================================================

`, len(k.targetValAddrs))
	for i, valAddr := range k.targetValAddrs {
		fmt.Println(i, valAddr.String())
	}

	fmt.Printf(`
============================================================
                         %d workers                        
============================================================

`, len(k.workers))
	for i, worker := range k.workers {
		fmt.Println(i, worker.String())
	}

	fmt.Printf(`
============================================================
                     expected parameters   			      
============================================================

validator's number in top 21: 	
           %s
percentage to plunder:
           %s
percentage to dominate:
           %s
`,
		k.expectedParams.GetExpectedValNumberInTop21().RoundInt(),
		k.expectedParams.GetExpectedPercentToPlunder(),
		k.expectedParams.GetExpectedPercentToDominate(),
	)
}
