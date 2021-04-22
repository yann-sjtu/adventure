package transfer

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/common"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	round        int
	privkey      string
	contractAddr string
	sleepTime    int
	mode         string
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   testLoop,
	}

	flags := cmd.Flags()
	flags.StringVarP(&privkey, "private-key", "p", "", "set the Priv Key")
	flags.StringVarP(&contractAddr, "contract", "c", "", "set the Contract addr")
	flags.IntVarP(&sleepTime, "sleep-time", "t", 0, "set the sleep time")
	flags.IntVarP(&round, "round", "r", 0, "set the round number")

	return cmd
}

func testLoop(cmd *cobra.Command, args []string) {
	Init()
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.015okt", types.BroadcastSync, 1500000)
	cli := clients.GetClient()
	info, err := utils.CreateAccountWithPrivateKey(privkey, "acc", common.PassWord)
	if err != nil {
		panic(err)
	}
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	offset := uint64(0)

	ethAddr, err := utils.ToHexAddress(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	approvePayloadStr := hexutil.Encode(approvePayload)
	res, err := cli.Evm().SendTxEthereum(privkey, contractAddr, "", approvePayloadStr, 2000000, accInfo.GetSequence()+offset)
	if err != nil {
		log.Printf("[%s] %s failed to approve in %s: %s\n", res.TxHash, ethAddr, contractAddr, err)
		return
	} else {
		log.Printf("[%s] %s approve successfull in %s \n", res.TxHash, ethAddr, contractAddr)
		offset++
	}
	time.Sleep(time.Second * time.Duration(sleepTime))

	transferPayloadStr := hexutil.Encode(transferPayload)
	for i := 0; i < round; i++ {
		res, err := cli.Evm().SendTxEthereum(privkey, contractAddr, "", transferPayloadStr, 2000000, accInfo.GetSequence()+offset)
		if err != nil {
			log.Printf("[%s] %s failed to transfer in %s: %s\n", res.TxHash, ethAddr, contractAddr, err)
			continue
		} else {
			log.Printf("[%s] %s transfer successfull in %s \n", res.TxHash, ethAddr, contractAddr)
			offset++
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}