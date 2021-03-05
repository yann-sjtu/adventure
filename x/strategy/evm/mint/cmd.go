package mint

import (
	"log"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/TTotken"
	"github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain-go-sdk/utils"
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
	flags.StringVarP(&MnemonicPath, "mnemonic-path", "m", "", "set the MnemonicPath path")

	return cmd
}

var (
	GoroutineNum  = 1

	MnemonicPath = ""
)

const TTokenAddr = "0x7F7715DC893A7504D3a89c8784bC4cFa19db8cc0"

func mint(cmd *cobra.Command, args []string) {
	infos := common.GetAccountManagerFromFile(MnemonicPath)
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.0005okt", types.BroadcastSync,500000)

	//succ, fail := tools.NewCounter(-1), tools.NewCounter(-1)
	var wg sync.WaitGroup
	for i := 0; i < GoroutineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// 1. get one of the eth private key
				info := infos.GetInfo()
				// 2. get cli
				cli := clients.GetClient()
				// 3. get acc number
				acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
				if err != nil {
					log.Println(err)
					continue
				}
				accNum, seqNum := acc.GetAccountNumber(), acc.GetSequence()
				offset := uint64(0)
				ethAddr, _ := utils.ToHexAddress(info.GetAddress().String())

				for {
					// Let Us GO GO GO !!!!!!
					// 1. mint
					payload := TTotken.BuildTTokenMintPayload(ethAddr.String(), sdk.NewDec(1).Int)
					//res, err :=
					res, err := cli.Evm().SendTx(info, common.PassWord, TTokenAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum+offset)
					if err != nil {
						log.Printf("[%s] %s failed to mint in %s: %s\n", res.TxHash, ethAddr, TTokenAddr, err)
						continue
					} else {
						log.Printf("[%s] %s mint in %s \n",  res.TxHash, ethAddr, TTokenAddr)
						offset++
					}
				}
			}
		}()
	}
	wg.Wait()
}