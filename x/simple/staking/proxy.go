package staking

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/utils"
)

const (
	regProxy    = "reg_proxy"
	unregProxy  = "unreg_proxy"
	bindProxy   = "bind_proxy"
	unbindProxy = "unbind_proxy"
)

var (
	once       sync.Once
	proxyInfos []keys.Info
)

func Proxy(cli *gosdk.Client, info keys.Info) {
	once.Do(func() {
		proxyMnemonics := config.Cfg.Staking.ProxyConfig.ProxyMnemonics
		for i, m := range proxyMnemonics {
			info, _, err := utils.CreateAccountWithMnemo(m, fmt.Sprintf("proxy%d", i), common.PassWord)
			if err != nil {
				log.Printf("create info when generating proxy account infos. error: %s", err)
			}
			proxyInfos = append(proxyInfos, info)
			sendTx(cli, info, delegate)
			sendTx(cli, info, regProxy)
		}
	})

	rand.Seed(time.Now().Unix())
	switch rand.Intn(6) {
	case 0: // bind proxy
		bindProxyTx(cli, info, proxyInfos[rand.Intn(len(proxyInfos))])
	case 1: // delegate
		sendTx(cli, info, delegate)
	case 2: // unbind proxy
		sendTx(cli, info, unbindProxy)
		bindProxyTx(cli, info, proxyInfos[rand.Intn(len(proxyInfos))])
	case 3: // unbond
		sendTx(cli, info, unbond)
	case 4: // proxy addr to vote
		info := proxyInfos[rand.Intn(len(proxyInfos))]
		sendTx(cli, info, regProxy)
		sendTx(cli, info, vote)
	case 5: // proxy addr to delegate and register and unbond
		info := proxyInfos[rand.Intn(len(proxyInfos))]
		sendTx(cli, info, delegate)
		sendTx(cli, info, regProxy)
		sendTx(cli, info, unbond)
	case 6: // unreg
		info := proxyInfos[rand.Intn(len(proxyInfos))]
		if (rand.Intn(20)+1)%10 == 0 {
			sendTx(cli, info, unregProxy)
		}
		sendTx(cli, info, regProxy)
	}
}

func bindProxyTx(cli *gosdk.Client, info keys.Info, proxy keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, bindProxy, info)
		return
	}

	_, err = cli.Staking().BindProxy(info, common.PassWord,
		proxy.GetAddress().String(),
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, bindProxy, info)
		return
	}
	common.PrintExecuteTxSuccess(bindProxy, info)
}
