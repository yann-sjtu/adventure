package rest_test

import (
	"errors"
	"log"
	"time"

	"github.com/okex/adventure/x/strategy/evm-unfinish/rest-test/utils"
	"github.com/spf13/cobra"
)

const (
	// to: 0x0000000000000000000000000000000000000001
	// value: 1
	erc20TransferPayload = "0xa9059cbb00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
)

var (
	Erc20ContractAddr string
)

func erc20TransferStrategyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-erc20-transfer-test",
		Short: "run rest erc20 transfer test",
		Args:  cobra.NoArgs,
		RunE:  runErc20TransferStrategyCmd,
	}

	flag := cmd.Flags()
	flag.StringVarP(&HostUrl, "url", "u", "http://localhost:8545", "host url")
	flag.StringVarP(&Erc20ContractAddr, "contract-addr", "c", "", "contract address")

	return cmd
}

func runErc20TransferStrategyCmd(cmd *cobra.Command, args []string) error {
	if len(Erc20ContractAddr) == 0 {
		return errors.New("failed. empty contract address")
	}

	fromAddr, err := utils.GetAddress(HostUrl)
	if err != nil {
		return err
	}

	fromAddrStr = fromAddr.Hex()

	go erc20Transfer()
	for {
		go erc20EstimateGas()
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func erc20Transfer() {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["to"] = Erc20ContractAddr
	param[0]["data"] = erc20TransferPayload

	for {
		_, err := utils.CallWithError("eth_sendTransaction", param, HostUrl)
		if err != nil {
			continue
		}

		log.Printf("%s transfers 0.000000000000000001btc(erc20) to 0x0000000000000000000000000000000000000001 successfully\n", fromAddrStr)
	}
}

func erc20EstimateGas() {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["to"] = Erc20ContractAddr
	param[0]["data"] = erc20TransferPayload
	_, err := utils.CallWithError("eth_estimateGas", param, HostUrl)
	if err != nil {
		return
	}
	log.Printf("%s estimate gas with transferring 0.000000000000000001btc(erc20) to 0x0000000000000000000000000000000000000001 successfully\n", fromAddrStr)
}
