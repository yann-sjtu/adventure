package bench

import (
	"crypto/ecdsa"
	"fmt"

	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain-go-sdk/utils"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var (
	concurrency int
	sleepTime   int
	privkPath   string
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

func getCosmosAddress(privateKey *ecdsa.PrivateKey) sdk.Address {
	cosmosAddr, err := utils.ToCosmosAddress(getEthAddress(privateKey).String())
	if err != nil {
		panic(err)
	}
	return cosmosAddr
}