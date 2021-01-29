package types

import (
	"fmt"
	gosdk "github.com/okex/okexchain-go-sdk"
	"sort"
	"testing"
)

const (
	URL     = "http://okexchain-rpc1.okexcn.com:26657"
	chainID = "okexchain-66"
)

func TestValidators(t *testing.T) {
	config, err := gosdk.NewClientConfig(URL, chainID, gosdk.BroadcastBlock, "", 2000000,
		1.1, "0.00000001okt")
	if err != nil {
		panic(err)
	}

	cli := gosdk.NewClient(config)
	vals, err := cli.Staking().QueryValidators()
	if err != nil {
		panic(err)
	}

	validators:=Validators(vals)
	sort.Sort(validators)
	fmt.Println(validators)
}
