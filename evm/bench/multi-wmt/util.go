package multiwmt

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethcompatible "github.com/okex/exchain-ethereum-compatible/utils"
	"io/ioutil"
	"math/big"
	"time"
)

func calHash(tx *types.Transaction) common.Hash {
	if !useOldTxHash {
		return tx.Hash()
	}
	h, err := ethcompatible.Hash(tx)
	if err == nil {
		return h
	}
	return common.Hash{}
}

func getReceipt(txs []*types.Transaction) error {
	cnt := 0
	for cnt < 100 {
		time.Sleep(2000 * time.Millisecond)

		_, errFirst := client.TransactionReceipt(context.Background(), calHash(txs[0]))
		_, errLast := client.TransactionReceipt(context.Background(), calHash(txs[len(txs)-1]))
		if errFirst == nil && errLast == nil {
			return nil
		}
		cnt++
	}

	panicLog := "failed txs:"
	for _, v := range txs {
		panicLog += v.Hash().String() + "  "
	}
	return errors.New(panicLog)
}

func getPrivateKey(key string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(key)
	panicerr(err)
	return privateKey
}

func transferOkt(key string, to common.Address, nonce uint64, value *big.Int) *types.Transaction {
	privateKey := getPrivateKey(key)

	tx, err := types.SignTx(types.NewTransaction(nonce, to, value, gasLimit, gasPrice, nil), signer, privateKey)
	panicerr(err)
	return tx
}

func SignTxWithNonce(privateKey *ecdsa.PrivateKey, to common.Address, payLoad []byte, nonce uint64) *types.Transaction {
	tx, err := types.SignTx(types.NewTransaction(nonce, to, new(big.Int), gasLimit, gasPrice, payLoad), signer, privateKey)
	panicerr(err)
	return tx
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

type SwapContract struct {
	Token0          common.Address
	Token1          common.Address
	Token2          common.Address
	Token3          common.Address
	FakewethAddress common.Address
	Factory         common.Address
	Router          common.Address
	Lp1             common.Address
	Lp2             common.Address
	StakingRewards1 common.Address
	StakingRewards2 common.Address
}

func LoadContractList(file string) []SwapContract {
	data, err := ioutil.ReadFile(file)
	panicerr(err)
	cList := make([]SwapContract, 0)
	err = json.Unmarshal(data, &cList)
	panicerr(err)
	return cList
}

func keyToAcc(key string) *acc {
	privateKey := getPrivateKey(key)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("keyToAcc")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &acc{
		privateKey: key,
		ecdsaPriv:  privateKey,
		ethAddress: fromAddress,
	}
}

func SendTxs(txs []*types.Transaction) error {
	for index, v := range txs {
		//time.Sleep(200 * time.Microsecond)
		cnt := 0
		for cnt < 10 {
			cnt++
			if err := client.SendTransaction(context.Background(), v); err != nil {
				fmt.Println("index", index, err)
				return err
			} else {
				break
			}
		}
	}
	return nil
}

type wmtConfig struct {
	RPC          string
	ContractPath string
	SuperAcc     string
	WorkerPath   string
	ParaNum      int
	UseOldTxHash bool
}

func loadWMTConfig(file string) *wmtConfig {
	data, err := ioutil.ReadFile(file)
	panicerr(err)
	c := new(wmtConfig)
	err = json.Unmarshal(data, c)
	panicerr(err)
	return c
}
