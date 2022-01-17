package multiwmt

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
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
		succCnt := 0
		for _, tx := range txs {
			_, err := client.TransactionReceipt(context.Background(), calHash(tx))
			if err == nil {
				succCnt++
			} else {
				break
			}
		}

		if succCnt == len(txs) {
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
	err = client.SendTransaction(context.Background(), tx)
	panicerr(err)
	return tx
}

func SendTxWithNonce(key string, to common.Address, payLoad []byte, nonce uint64) *types.Transaction {
	privateKey := getPrivateKey(key)

	tx, err := types.SignTx(types.NewTransaction(nonce, to, new(big.Int), gasLimit, gasPrice, payLoad), signer, privateKey)
	panicerr(err)
	err = client.SendTransaction(context.Background(), tx)
	panicerr(err)
	return tx
}

func SendTx(key string, to common.Address, payLoad []byte) error {
	privateKey := getPrivateKey(key)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("should panic")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	tx, err := types.SignTx(types.NewTransaction(nonce, to, new(big.Int), gasLimit, gasPrice, payLoad), signer, privateKey)
	panicerr(err)

	cnt := 0
	for true {
		cnt++
		err = client.SendTransaction(context.Background(), tx)
		if err == nil {
			break
		}
		if err != nil && cnt > 10 {
			return err
		}
	}

	return getReceipt([]*types.Transaction{tx})
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
