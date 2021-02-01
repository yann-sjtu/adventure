package keeper

import "fmt"

func (k *Keeper) logInit() {
	fmt.Printf(`
============================================================
               our %d validators    	
============================================================

`, len(k.ourValAddrs))
	for i, valAddrStr := range k.ourValAddrs {
		fmt.Println(i, valAddrStr)
	}

	fmt.Printf(`
============================================================
               our top %d validators    	
============================================================

`, len(k.ourTop18ValAddrs))
	for i, valAddr := range k.ourTop18ValAddrs {
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

percentage to plunder:
           %s
`, k.plunderedPct.String(),
	)
}
