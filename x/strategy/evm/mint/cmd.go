package mint

import (
	"log"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/TTotken"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

func MintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   mint,
	}

	flags := cmd.Flags()
	flags.IntVarP(&GoroutineNum, "goroutine-num", "g", 1, "set Goroutine Num of deploying contracts")
	flags.StringVarP(&PrivKeysPath, "privkeys-path", "p", "", "set the PrivKeysPath path")
	flags.StringVarP(&TTokenAddr, "token-contract", "c", "", "set the ttoken contract addr")

	return cmd
}

var (
	GoroutineNum = 1
	PrivKeysPath = ""
	TTokenAddr   = ""
)

func mint(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(PrivKeysPath)
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.0005okt", types.BroadcastSync, 500000)

	//succ, fail := tools.NewCounter(-1), tools.NewCounter(-1)
	var wg sync.WaitGroup
	for i := 0; i < GoroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for {
				// 1. get one of the eth private key
				privkey := privkeys[index]
				info, err := utils.CreateAccountWithPrivateKey(privkey, "acc", common.PassWord)
				if err != nil {
					panic(err)
				}
				ethAddr, err := utils.ToHexAddress(info.GetAddress().String())
				if err != nil {
					panic(err)
				}
				// 2. get cli
				cli := clients.GetClient()
				// 3. get acc number
				acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
				if err != nil {
					log.Println(err)
					continue
				}
				seqNum, offset := acc.GetSequence(), uint64(0)

				// Let Us GO GO GO !!!!!!
				// 1. mint
				payload := TTotken.BuildTTokenMintPayload(ethAddr.String(), sdk.NewDec(1).Int)
				for {
					//res, err :=
					res, err := cli.Evm().SendTxEthereum(privkey, TTokenAddr, "", ethcommon.Bytes2Hex(payload), 500000, seqNum+offset)
					if err != nil {
						log.Printf("[%s] %s failed to mint in %s: %s\n", res.TxHash, ethAddr, TTokenAddr, err)
						continue
					} else {
						log.Printf("[%s] %s mint in %s \n", res.TxHash, ethAddr, TTokenAddr)
						offset++
					}
					time.Sleep(time.Second*2)
				}
			}
		}(i)
	}
	wg.Wait()
}
