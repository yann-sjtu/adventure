package tools

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
)

type ContractManager struct {
	usdtContracts []*Contract
	lock          *sync.RWMutex
}

type Contract struct {
	ContractAddr string
	Owner        keys.Info
}

func NewContractManager() *ContractManager {
	return &ContractManager{
		usdtContracts: make([]*Contract, 0, 1000),
		lock:          new(sync.RWMutex),
	}
}

func (c *ContractManager) AppendUSDTAddr(addr string, key keys.Info) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.usdtContracts = append(c.usdtContracts, &Contract{addr, key})
}

func (c *ContractManager) GetOneRandomUSDTOwner() *Contract {
	c.lock.RLock()
	defer c.lock.RUnlock()
	rand.Seed(time.Now().UnixNano())
	length := len(c.usdtContracts)
	return c.usdtContracts[rand.Intn(length)]
}

func (c *ContractManager) List() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	log.Println("Total usdt contracts:", len(c.usdtContracts))
	for i, contract := range c.usdtContracts {
		fmt.Printf("[%d] contract addr: %s, owner: %s\n", i, contract.ContractAddr, contract.Owner.GetAddress().String())
	}
}
