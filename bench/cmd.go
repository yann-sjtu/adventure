package bench

import (
	"crypto/ecdsa"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

var (
	concurrency int
	sleepTimeTx int
	privkPath       string
	rpc_hosts    []string

	chainID       string
)

func BenchCmd() *cobra.Command {
	var evmCmd = &cobra.Command{
		Use:   "evm-bench",
		Short: "evm web3 cli of test strategy",
	}

	evmCmd.AddCommand(
		InitStorageCmd(),
		OperateCmd(),
	)
	return evmCmd
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