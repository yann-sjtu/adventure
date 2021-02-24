package utils

import (
	"encoding/json"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"log"
	"time"
)

const (
	DefaultHostUrl = "http://localhost:8545"
)

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

func WaitForReceipt(hash ethcmn.Hash) map[string]interface{} {
	for i := 0; i < 12; i++ {
		receipt, err := GetTransactionReceipt(hash)
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

func GetTransactionReceipt(hash ethcmn.Hash) (receipt map[string]interface{}, err error) {
	param := []string{hash.Hex()}
	rpcRes, err := CallWithError("eth_getTransactionReceipt", param, "")
	if err != nil {
		return
	}

	receipt = make(map[string]interface{})
	if err = json.Unmarshal(rpcRes.Result, &receipt); err != nil {
		return
	}

	return
}
