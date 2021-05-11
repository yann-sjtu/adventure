package common

import (
	"math/rand"
	"sync"
	"time"

	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/types"
)

type ClientManager struct {
	i       int
	clients []*gosdk.Client
	lock    *sync.Mutex
}

func (r *ClientManager) GetClient() *gosdk.Client {
	r.lock.Lock()
	defer r.lock.Unlock()
	k := r.i
	r.i = (r.i + 1) % len(r.clients)
	return r.clients[k]
}

func (r *ClientManager) GetRandomClient() *gosdk.Client {
	rand.Seed(time.Now().UnixNano())
	return r.clients[rand.Intn(len(r.clients))]
}

func NewClientManager(hosts []string, fee string, gas ...uint64) *ClientManager {
	clients := getAllInvariantClients(hosts, fee, gas...)
	control := &ClientManager{
		0,
		clients,
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
		cfg, _ = types.NewClientConfig(host, GlobalConfig.Networks[""].ChainId, types.BroadcastBlock, "", 350000, 1.5, "0.000000001"+NativeToken)
	} else {
		if len(gas) != 0 {
			cfg, _ = types.NewClientConfig(host, GlobalConfig.Networks[""].ChainId, types.BroadcastBlock, fee, gas[0], 0, "")
		} else {
			cfg, _ = types.NewClientConfig(host, GlobalConfig.Networks[""].ChainId, types.BroadcastBlock, fee, 200000, 0, "")
		}
	}
	return
}

func NewClientManagerWithMode(hosts []string, fee string, mode string, gas ...uint64) *ClientManager {
	clients := getAllInvariantClientsWithMode(hosts, fee, mode, gas...)
	control := &ClientManager{
		0,
		clients,
		new(sync.Mutex),
	}
	return control
}

func getAllInvariantClientsWithMode(hosts []string, fee string, mode string, gas ...uint64) []*gosdk.Client {
	var clients []*gosdk.Client
	for i := 0; i < len(hosts); i++ {
		cfg := initClientConfigWithMode(fee, hosts[i], mode, gas...)
		cli := gosdk.NewClient(cfg)
		clients = append(clients, &cli)
	}
	return clients
}

func initClientConfigWithMode(fee string, host string, mode string, gas ...uint64) (cfg types.ClientConfig) {
	if fee == AUTO {
		cfg, _ = types.NewClientConfig(host, GlobalConfig.Networks[""].ChainId, mode, "", 350000, 1.5, "0.000000001"+NativeToken)
	} else {
		if len(gas) != 0 {
			cfg, _ = types.NewClientConfig(host, GlobalConfig.Networks[""].ChainId, mode, fee, gas[0], 0, "")
		} else {
			cfg, _ = types.NewClientConfig(host, GlobalConfig.Networks[""].ChainId, mode, fee, 200000, 0, "")
		}
	}
	return
}
