package evm_transfer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	adcomm "github.com/okex/adventure/common"
	"github.com/okex/exchain-ethereum-compatible/utils"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/rand"
)

var (
	concurrencyTx   int
	sleepTimeTx     int
	privkPath       string
	rest_hosts    []string

	chainID       *big.Int
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

	flags.StringVar(&privkPath, "privkeyPath", "", "")
	flags.StringSliceVar(&rest_hosts, "rest-hosts", []string{}, "")
	return cmd
}

func transferTx(cmd *cobra.Command, args []string) {
	chainID = getChainId()
	log.Println("chain id", chainID.Uint64())

	privkeys := adcomm.GetPrivKeyFromPrivKeyFile(privkPath)
	var wg sync.WaitGroup
	wg.Add(concurrencyTx)
	for i := 0; i < concurrencyTx; i++ {
		go func(index int, privkey string) {
			defer wg.Done()

			restHost := rest_hosts[index%len(rest_hosts)]
			transfer(privkey, restHost)

		}(i, privkeys[i])
	}
	wg.Wait()
}

func getChainId() *big.Int {
	client, err := ethclient.Dial(rest_hosts[0])
	if err != nil {
		panic(err)
	}

	id, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	return id
}

func transfer(privkey string, host string) {
	privateKey, _ := crypto.HexToECDSA(privkey)
	client, err := ethclient.Dial(host)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), getEthAddress(privateKey))
	if err != nil {
		log.Fatalf("failed to fetch noce: %+v", err)
	}

	for {
		tx := types.NewTransaction(
			nonce,
			common.BytesToAddress(crypto.Keccak256(rand.Bytes(64))),
			big.NewInt(1),
			uint64(21000),
			big.NewInt(1e9),
			nil)
		signedTx, err := types.SignTx(
			tx,
			types.NewEIP155Signer(chainID),
			privateKey,
		)
		if err != nil {
			log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
		}

		// 3. send rawTx
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Printf("[%s] err: %s\n", getEthAddress(privateKey), err)
			continue
		}
		txhash, _ := utils.Hash(signedTx)
		fmt.Println(txhash)
		nonce++
		time.Sleep(time.Second * time.Duration(sleepTimeTx))
	}
}

func getEthAddress(privateKey *ecdsa.PrivateKey) common.Address {
	pubkeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if ok != true {
		panic(fmt.Errorf("convert into pubkey failed"))
	}
	fromAddress := crypto.PubkeyToAddress(*pubkeyECDSA)
	return fromAddress
}