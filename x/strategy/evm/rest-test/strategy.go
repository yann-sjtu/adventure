package rest_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/x/strategy/evm/rest-test/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var (
	fromAddrStr string
)

func strategyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-test",
		Short: "run rest test",
		Args:  cobra.NoArgs,
		RunE:  runStrategyCmd,
	}

	cmd.Flags().StringVarP(&HostUrl, "url", "u", "http://localhost:8545", "host url")

	return cmd
}

func runStrategyCmd(cmd *cobra.Command, args []string) error {
	fromAddr, err := utils.GetAddress("")
	if err != nil {
		return err
	}

	fromAddrStr = fromAddr.Hex()

	go transfer()
	estimateGas()
	time.Sleep(time.Minute)
	return nil
}

func transfer() {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["value"] = (*hexutil.Big)(sdk.OneDec().BigInt()).String()

	for {
		receiverAddr := utils.GetReceiverAddrRandomly()
		param[0]["to"] = receiverAddr
		_, err := utils.CallWithError("eth_sendTransaction", param, HostUrl)
		if err != nil {
			continue
		}

		log.Printf("%s transfers 1okt to %s successfully\n", fromAddrStr, receiverAddr)
		time.Sleep(500 * time.Millisecond)
	}
}

func estimateGas() {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["value"] = (*hexutil.Big)(sdk.OneDec().BigInt()).String()

	for {
		receiverAddr := utils.GetReceiverAddrRandomly()
		param[0]["to"] = receiverAddr
		_, err := utils.CallWithError("eth_estimateGas", param, HostUrl)
		if err != nil {
			continue
		}

		log.Printf("%s estimate gas with transferring 1okt to %s successfully\n", fromAddrStr, receiverAddr)
	}
}
