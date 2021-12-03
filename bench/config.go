package bench

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
)

type Param struct {
	concurrency int

	rpcHosts []string
	chainID  string
	privKeys []string

	ethPort int
}

func RunTxs(p Param, e func(*gosdk.Client, *Account)) {
	clients := GenerateClients(p.rpcHosts, p.chainID)
	accounts := GenerateAccounts(p.privKeys)

	for i := 0; i < p.concurrency; i++ {
		go func(index int) {
			for j := 0; ; j++ {
				id := index + j*p.concurrency%len(accounts)
				account := accounts[id]
				host := p.rpcHosts[index%len(p.rpcHosts)]
				account.SetNonce(host, p.chainID, p.ethPort)

				client := clients[index%len(clients)]

				e(client, account)
			}
		}(i)
	}

	select {}
}

type Account struct {
	index      int
	nonce      uint64
	queried    bool
	privateKey *ecdsa.PrivateKey
}

func (acc *Account) SetNonce(host string, chainID string, ethPort int) {
	if acc.queried {
		return
	}

	var nonce uint64
	if ethPort == 0 {
		cfg, err := types.NewClientConfig(host, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
		if err != nil {
			panic(err)
		}
		cli := gosdk.NewClient(cfg)

		addr := getCosmosAddress(acc.privateKey)
		accInfo, err := cli.Auth().QueryAccount(addr.String())
		if err != nil {
			log.Printf("[%d] query %s error: %s\n", acc.index, addr.String(), err)
			return
		}
		nonce = accInfo.GetSequence()
	} else {
		str := strings.Split(host, ":")
		ethhost := str[0] + ":" + str[1] + ":" + strconv.Itoa(ethPort)
		cli, err := ethclient.Dial(ethhost)
		if err != nil {
			panic(err)
		}

		addr := getEthAddress(acc.privateKey)
		nonce, err = cli.PendingNonceAt(context.Background(), addr)
		if err != nil {
			log.Printf("[%d] query %s error: %s\n", acc.index, addr.String(), err)
			return
		}
	}

	acc.nonce = nonce
	acc.queried = true
}

func (acc *Account) AddNonce() {
	acc.nonce += 1
}

func (acc *Account) GetNonce() uint64 {
	return acc.nonce
}

func (acc *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return acc.privateKey
}

func GenerateAccounts(privkeys []string) (accounts []*Account) {
	for i, p := range privkeys {
		privateKey, err := crypto.HexToECDSA(p)
		if err != nil {
			panic(err)
		}

		accounts = append(accounts, &Account{i, 0, false, privateKey})
	}
	return
}

func GenerateClients(hosts []string, chainID string) (clients []*gosdk.Client) {
	for _, h := range hosts {
		cfg, err := types.NewClientConfig(h, chainID, types.BroadcastSync, "", 2000000, 1.5, "0.0000000001"+common.NativeToken)
		if err != nil {
			panic(fmt.Errorf("initialize client failed: %s", err))
		}
		cli := gosdk.NewClient(cfg)

		clients = append(clients, &cli)
	}
	return
}
