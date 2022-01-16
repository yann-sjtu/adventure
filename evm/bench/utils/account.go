package utils

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/client"
)

type EthAccount struct {
	lock       *sync.Mutex
	nonce      uint64
	queried    bool
	privateKey *ecdsa.PrivateKey
}

func generateAccounts(privkeys []string) (accounts []*EthAccount) {
	for _, p := range privkeys {
		privateKey, err := crypto.HexToECDSA(p)
		if err != nil {
			panic(err)
		}

		accounts = append(accounts, &EthAccount{new(sync.Mutex), 0, false, privateKey})
	}
	return
}

func (a *EthAccount) Lock() {
	a.lock.Lock()
}

func (a *EthAccount) Unlock() {
	a.lock.Unlock()
}

func (a *EthAccount) SetNonce(cli client.Client) error {
	if a.queried {
		return nil
	}

	nonce, err := cli.QueryNonce(common.GetEthAddressFromPK(a.privateKey).String())
	if err != nil {
		return err
	}

	a.nonce = nonce
	a.queried = true
	return nil
}

func (a *EthAccount) AddNonce() {
	a.nonce += 1
}

func (a *EthAccount) GetNonce() uint64 {
	return a.nonce
}

func (a *EthAccount) GetPrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}
