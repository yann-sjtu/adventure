package mint

import (
	"log"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	deploy_contracts "github.com/okex/adventure/x/strategy/evm/deploy-contracts"
	"github.com/okex/adventure/x/strategy/evm/template/TTotken"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2Staker"
	"github.com/okex/adventure/x/strategy/evm/tools"
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

	return cmd
}

const TTokenAddr = "0x0d021d10ab9E155Fc1e8705d12b73f9bd3de0a36"

func mint(cmd *cobra.Command, args []string) {
	infos := common.GetAccountManagerFromFile(deploy_contracts.MnemonicPath)
	clients := common.NewClientManagerWithMode(common.Cfg.Hosts, "0.05okt", types.BroadcastBlock,50000000)

	succ, fail := tools.NewCounter(-1), tools.NewCounter(-1)
	var wg sync.WaitGroup
	for i := 0; i < deploy_contracts.GoroutineNum; i++ {
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
				ethAddr := utils.GetEthAddressStrFromCosmosAddr(info.GetAddress())

				// Let Us GO GO GO !!!!!!
				// 1. add liquididy
				payload := TTotken.BuildTTokenMintPayload("", sdk.NewDec(1).Int)
				//res, err :=
				res, err := cli.Evm().SendTx(info, common.PassWord, TTokenAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
				if err != nil {
					log.Printf("(%d)[%s] %s failed to add liquidity in %s: %s\n", fail.Add(), res.TxHash, ethAddr, routerAddr, err)
					continue
				} else {
					log.Printf("(%d)[%s] %s add liquidity in %s \n", succ.Add(), res.TxHash, ethAddr, routerAddr)
					offset++
				}
			}
		}()
	}
	wg.Wait()
}