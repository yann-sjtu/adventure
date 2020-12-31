package common

import (
	"math/rand"
	"sync"
	"time"

	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/types"
)

type ClientManager struct {
	i       int
	clients []*gosdk.Client
	sum     int
	lock    *sync.Mutex
}

func (r *ClientManager) GetClient() *gosdk.Client {
	r.lock.Lock()
	defer r.lock.Unlock()
	k := r.i
	r.i = (r.i + 1) % r.sum
	return r.clients[k]
}

func (r *ClientManager) GetRandomClient() *gosdk.Client {
	rand.Seed(time.Now().UnixNano())
	return r.clients[rand.Intn(r.sum)]
}

func NewClientManager(hosts []string, fee string, gas ...uint64) *ClientManager {
	clients := getAllInvariantClients(hosts, fee, gas...)
	control := &ClientManager{
		0,
		clients,
		len(clients),
		new(sync.Mutex),
	}
	return control
}

func getAllInvariantClients(hosts []string, fee string, gas ...uint64) []*gosdk.Client {
	var clients []*gosdk.Client
	for i := 0; i < len(hosts); i++ {
		cfg := initClientConfig(fee, hosts[i], gas...)
		cli := gosdk.NewClient(cfg)
		clients = append(clients, &cli)
	}
	return clients
}

func initClientConfig(fee string, host string, gas ...uint64) (cfg types.ClientConfig) {
	if fee == AUTO {
		cfg, _ = types.NewClientConfig(host, "okexchaintestnet-1", types.BroadcastBlock, "", 20000000, 1.5, "0.0000001"+NativeToken)
	} else {
		if len(gas) != 0 {
			cfg, _ = types.NewClientConfig(host, "okexchaintestnet-1", types.BroadcastBlock, fee, gas[0], 0, "")
		} else {
			cfg, _ = types.NewClientConfig(host, "okexchaintestnet-1", types.BroadcastBlock, fee, 20000000, 0, "")
		}
	}
	return
}
