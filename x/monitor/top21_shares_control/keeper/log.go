package keeper

import "fmt"

func (k *Keeper) logInit() {
	fmt.Printf(`
============================================================
               %d target validators    	
============================================================

`, len(k.targetValAddrs))
	for i, valAddrStr := range k.targetValAddrs {
		fmt.Println(i, valAddrStr)
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

percentage to dominate:
           %s
`, k.dominationPct.String(),
	)
}
