package evm_transfer

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/rand"
)

var (
	concurrencyTx   int
	sleepTimeTx     int
	privkPath       string
	rpc_hosts    []string

	chainID       string

	fixed  bool
)

func TransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evm-transfer",
		Short: "",
		Long:  "",
		Run:   transferTx,
	}
	flags := cmd.Flags()
	flags.IntVar(&concurrencyTx, "concurrency", 10, "set the number of tx number per second")
	flags.IntVar(&sleepTimeTx, "sleepTime", 1, "")

	flags.StringVar(&privkPath, "privkey-path", "", "")
	flags.StringSliceVar(&rpc_hosts, "rpc-hosts", []string{}, "")
	flags.StringVar(&chainID, "chain-id", "", "")

	flags.BoolVar(&fixed, "fixed", false, "")
	return cmd
}

func transferTx(cmd *cobra.Command, args []string) {
	privkeys := common.GetPrivKeyFromPrivKeyFile(privkPath)
	to := ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64)))
	for i := 0; i < concurrencyTx; i++ {
		go func(index int, privkey string) {
			if !fixed {
				to = ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64)))
			}
			rpcHost := rpc_hosts[index%len(rpc_hosts)]
			transfer(privkey, rpcHost, to.String())

		}(i, privkeys[i])
	}

	select {}
}

func transfer(privkey string, host string, to string) {
	cfg, _ := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 30000000, 1.5, "0.0000000001"+common.NativeToken)
	cli := gosdk.NewClient(cfg)

	addr := getCosmosAddress(privkey)
	accInfo, err := cli.Auth().QueryAccount(addr.String())
	if err != nil {
		panic(err)
	}
	nonce := accInfo.GetSequence()

	for {
		res, err := cli.Evm().SendTxEthereum(privkey, to, "0.000000001", "",21000, nonce)
		if err != nil {
			continue
		} else {
			log.Printf("txhash: %s\n", res.TxHash)
		}

		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}

func getEthAddress(privateKey *ecdsa.PrivateKey) ethcmm.Address {
	pubkeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if ok != true {
		panic(fmt.Errorf("convert into pubkey failed"))
	}
	fromAddress := crypto.PubkeyToAddress(*pubkeyECDSA)
	return fromAddress
}

func getCosmosAddress(privkey string) sdk.Address {
	privateKey, _ := crypto.HexToECDSA(privkey)
	cosmosAddr, err := utils.ToCosmosAddress(getEthAddress(privateKey).String())
	if err != nil {
		panic(err)
	}
	return cosmosAddr
}
