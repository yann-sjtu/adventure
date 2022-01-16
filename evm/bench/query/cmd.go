package query

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/okex/adventure/evm/constant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	opts []int
)

func QueryCmd() *cobra.Command {
	// add flags
	cmd := &cobra.Command{
		Use:   "query",
		Short: "benchmarking query, supports nine types of ethereum query interface",
		Run:   query,
	}
	flags := cmd.Flags()
	flags.IntSliceVarP(&opts, "opts", "o", []int{1, 1, 1, 1, 1, 1, 1, 1, 1}, "set the number of query concurrent number per second")
	return cmd
}

func query(cmd *cobra.Command, args []string) {
	if len(opts) != 9 {
		panic(fmt.Errorf("concurrent config length should be 9, acutal len: %d", len(opts)))
	}

	sleep := viper.GetInt(constant.FlagSleep)
	ips := viper.GetStringSlice(constant.FlagIPs)
	ip := ips[0]
	for r := 1; ; r++ {
		for n := 0; n < 9; n++ {
			reqType := n
			for i := 0; i < opts[reqType]; i++ {
				go func(round int, num int, typeIndex int) {
					req := generateRequest(reqType)
					call(req, typeIndex, ip)

				}(r, i, reqType)
			}
		}
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}
}

func generateRequest(index int) []byte {
	var req Request
	switch index {
	case 0:
		req = persistentBlockNumberRequest
	case 1:
		req = EthGetBalance()
	case 2:
		req = EthGetBlockByNumber()
	case 3:
		req = persistentGasPriceRequest
	case 4:
		req = persistentGetCodeReuqest
	case 5:
		req = EthGetTransactionCount()
	case 6:
		req = EthGetTransactionReceipt()
	case 7:
		req = NetVersion()
	case 8:
		req = EthCall()
	default:
		req = persistentBlockNumberRequest
	}

	postBody, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return postBody
}
