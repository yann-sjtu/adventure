package utils

import (
	"encoding/json"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"log"
	"math/rand"
	"time"
)

const (
	DefaultHostUrl = "http://localhost:8545"
	receiverNum    = 10
)

var (
	receiverAddrs [receiverNum]string
)

func init() {
	for i := 0; i < receiverNum; i++ {
		receiverAddrs[i] = ethcmn.BytesToAddress([]byte{byte(i)}).Hex()
	}
}

func GetAddress(hostUrl string) (addr ethcmn.Address, err error) {
	rpcRes, err := CallWithError("eth_accounts", []string{}, hostUrl)
	if err != nil {
		return
	}

	var res []ethcmn.Address
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		return
	}

	return res[0], nil
}

func WaitForReceipt(hash ethcmn.Hash, hostUrl string) map[string]interface{} {
	for i := 0; i < 12; i++ {
		receipt, err := GetTransactionReceipt(hash, hostUrl)
		if err != nil {
			log.Println(err)
		}

		if receipt != nil {
			return receipt
		}

		time.Sleep(time.Second)
	}

	return nil
}

func GetTransactionReceipt(hash ethcmn.Hash, hostUrl string) (receipt map[string]interface{}, err error) {
	param := []string{hash.Hex()}
	rpcRes, err := CallWithError("eth_getTransactionReceipt", param, hostUrl)
	if err != nil {
		return
	}

	receipt = make(map[string]interface{})
	if err = json.Unmarshal(rpcRes.Result, &receipt); err != nil {
		return
	}

	return
}

func GetReceiverAddrRandomly() string {
	rand.Seed(time.Now().UnixNano())
	return receiverAddrs[rand.Intn(receiverNum)]
}
