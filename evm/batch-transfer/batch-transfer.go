package batch_transfer

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/client"
	"github.com/okex/adventure/evm/constant"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	privateKey  string
	addressFile string
)

func batchTransfer(cmd *cobra.Command, args []string) {
	// 0.1 load env parameters
	cli, privateKey, addrs := loadEnv()
	// query nonce
	nonce, err := cli.QueryNonce(common.GetEthAddressFromPK(privateKey).String())
	if err != nil {
		log.Println(fmt.Errorf("failed to query nonce, error", err))
		return
	}

	// 1. deploy
	contractAddr, err := deploy(cli, privateKey, nonce)
	if err != nil {
		log.Println(fmt.Errorf("failed to deploy, error: %s", err))
		return
	}
	time.Sleep(time.Second * 5)

	// 2. transfers
	if err := transfers(cli, privateKey, nonce+1, contractAddr, sdk.MustNewDecFromStr(args[0]).Int, addrs); err != nil {
		log.Println(fmt.Errorf("failed to transfer, error: %s", err))
		return
	}
}

func loadEnv() (client.Client, *ecdsa.PrivateKey, []ethcmn.Address) {
	ips := viper.GetStringSlice(constant.FlagIPs)
	if len(ips) == 0 {
		panic(fmt.Errorf("ip list is nil, please set them in flag %s", constant.FlagIPs))
	}
	cli := client.NewClient(ips[0])

	privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to unencrypted private key [%s]: %s", privateKey, err))
	}

	addresses := constant.HexAddresses
	if addressFile != "" {
		addresses = common.ReadDataFromFile(addressFile)
	}
	hexAddrs := make([]ethcmn.Address, len(addresses), len(addresses))
	for i, addr := range addresses {
		hexAddrs[i] = ethcmn.HexToAddress(addr)
	}

	return cli, privateKey, hexAddrs
}

func deploy(cli client.Client, privateKey *ecdsa.PrivateKey, nonce uint64) (ethcmn.Address, error) {
	// deploy contract BatchTransfer contract
	txhash, err := cli.CreateContract(privateKey, nonce, nil, 300000, evmtypes.DefaultGasPrice, ethcmn.Hex2Bytes(constant.BatchTransferHex))
	if err != nil {
		return ethcmn.Address{}, err
	}
	contractAddr := crypto.CreateAddress(common.GetEthAddressFromPK(privateKey), nonce)
	log.Printf("caller: %s, nonce: %d, contract: %s, txhash: %s\n", common.GetEthAddressFromPK(privateKey), nonce, contractAddr, txhash)

	return contractAddr, nil
}

func transfers(cli client.Client, privateKey *ecdsa.PrivateKey, nonce uint64, to ethcmn.Address, amount *big.Int, addrs []ethcmn.Address) error {
	// load abi
	tABI, err := abi.JSON(strings.NewReader(constant.BatchTransferABI))
	if err != nil {
		return fmt.Errorf("failed to initialize BatchTransfer abi, error: %s", err)
	}

	batchNum := 200 // 40,000 gas per address
	totalAmount := big.NewInt(1).Mul(amount, big.NewInt(int64(batchNum)))
	for i := 0; i <= len(addrs)/batchNum && i*batchNum < len(addrs); i++ {
		start, end := i*batchNum, (i+1)*batchNum
		if end > len(addrs) {
			end = len(addrs)
		}
		txdata, err := tABI.Pack("transfers", addrs[start:end], amount)
		if err != nil {
			return fmt.Errorf("failed to pack BatchTransfer parameters, error: %s", err)
		}

		txhash, err := cli.SendTx(privateKey, nonce, to, totalAmount, uint64(41000*batchNum), evmtypes.DefaultGasPrice, txdata)
		if err != nil {
			return err
		}
		log.Printf("caller: %s, nonce: %d, to[%d:%d] txhash: %s\n", common.GetEthAddressFromPK(privateKey), nonce, start, end-1, txhash)

		nonce++
		time.Sleep(time.Second)
	}

	return nil
}
