package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_Balance(t *testing.T) {
	clients := []string{"ip1", "ip2", "ip3"}
	var accounts []string
	for i := 0; i < 10; i++ {
		accounts = append(accounts, fmt.Sprintf("acc%d", i))
	}

	concurrency := 4
	for i := 0; i < concurrency; i++ {
		go func(gIndex int) {
			time.Sleep(time.Second * time.Duration(gIndex))

			for j := 0; ; j++ {
				aIndex := (gIndex + j*concurrency) % len(accounts)
				acc := accounts[aIndex]

				cli := clients[aIndex%len(clients)]

				fmt.Println(acc, "->", cli)
				time.Sleep(time.Second * 10)
			}
		}(i)

	}

	select {}
}

//    aIndex         gIndex
//  acc0 -> ip1    acc0 -> ip1
//  acc1 -> ip2    acc1 -> ip2
//  acc2 -> ip3    acc2 -> ip3
//  acc3 -> ip1    acc3 -> ip1
//  acc4 -> ip2    acc4 -> ip1
//  acc5 -> ip3    acc5 -> ip2
//  acc6 -> ip1    acc6 -> ip3
//  acc7 -> ip2    acc7 -> ip1
//  acc8 -> ip3    acc8 -> ip1
//  acc9 -> ip1    acc9 -> ip2

//  acc0 -> ip1    acc0 -> ip3
//  acc1 -> ip2    acc1 -> ip1
//  acc2 -> ip3    acc2 -> ip1
//  acc3 -> ip1    acc3 -> ip2
//  acc4 -> ip2    acc4 -> ip3
//  acc5 -> ip3    acc5 -> ip1
//  acc6 -> ip1    acc6 -> ip1
//  acc7 -> ip2    acc7 -> ip2
//  acc8 -> ip3    acc8 -> ip3
//  acc9 -> ip1    acc9 -> ip1

//  acc0 -> ip1    acc0 -> ip1
//  acc1 -> ip2    acc1 -> ip2
//  acc2 -> ip3    acc2 -> ip3
//  acc3 -> ip1    acc3 -> ip1
//  acc4 -> ip2    acc4 -> ip1
//  acc5 -> ip3    acc5 -> ip2
//  acc6 -> ip1    acc6 -> ip3
//  acc7 -> ip2    acc7 -> ip1
//  acc8 -> ip3    acc8 -> ip1
//  acc9 -> ip1    acc9 -> ip2

//  acc0 -> ip1    acc0 -> ip3
//  acc1 -> ip2    acc1 -> ip1
