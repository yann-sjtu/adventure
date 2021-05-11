package staking

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/utils"
)

const (
	regProxy    = "reg_proxy"
	unregProxy  = "unreg_proxy"
	bindProxy   = "bind_proxy"
	unbindProxy = "unbind_proxy"
)

var proxyMnemonics = []string{
	"usual curve false good exhibit half panda olympic seminar member physical venue",
	"spike feature valid violin indoor asthma coral stable law inherit advice lava",
	"north zone firm disorder peace fantasy hamster company next error phone sorry",
	"fringe kitchen neglect laugh powder lake service industry loyal deputy seminar spider",
	"rebel time car food slam panther label speak sphere cram car hard",
	"swallow strong upset summer arctic young address engine social brain shy planet",
	"audit saddle boring satoshi capable shoot flight state embark thing must possible",
	"suggest gravity need grant permit raven exchange area moment auction ordinary boy",
	"lend second consider harsh found cruel focus creek dutch clerk sign garden",
	"trade milk enemy gain section slush flock bubble stereo indicate floor cram",
}

var (
	once       sync.Once
	proxyInfos []keys.Info
	length     int
)

func Proxy(cli *gosdk.Client, info keys.Info) {
	once.Do(func() {
		for i, m := range proxyMnemonics {
			info, _, err := utils.CreateAccountWithMnemo(m, fmt.Sprintf("proxy%d", i), common.PassWord)
			if err != nil {
				log.Printf("create info when generating proxy account infos. error: %s", err)
			}
			proxyInfos = append(proxyInfos, info)
			sendTx(cli, info, delegate)
			sendTx(cli, info, regProxy)
		}
		length = len(proxyInfos)
	})

	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(6) {
	case 0: // bind proxy
		bindProxyTx(cli, info, proxyInfos[rand.Intn(length)])
	case 1: // delegate
		sendTx(cli, info, delegate)
	case 2: // unbind proxy
		sendTx(cli, info, unbindProxy)
		bindProxyTx(cli, info, proxyInfos[rand.Intn(length)])
	case 3: // unbond
		sendTx(cli, info, unbond)
	case 4: // proxy addr to vote
		info := proxyInfos[rand.Intn(length)]
		sendTx(cli, info, regProxy)
		sendTx(cli, info, vote)
	case 5: // proxy addr to delegate and register and unbond
		info := proxyInfos[rand.Intn(length)]
		sendTx(cli, info, delegate)
		sendTx(cli, info, regProxy)
		sendTx(cli, info, unbond)
	case 6: // unreg
		info := proxyInfos[rand.Intn(length)]
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
