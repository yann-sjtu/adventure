package types

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	gokeys "github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/staking/validators/val"
	gosdk "github.com/okex/exchain-go-sdk"
)

var (
	_                val.Validator = (*Validator)(nil)
	CreateSuccessCnt int
	DestroySuccessCnt int
)

type Validator struct {
	id         string
	operKey    gokeys.Info
	consPubkey string
}

func NewValidator(operKey gokeys.Info, id int, consPubKey string) *Validator {
	return &Validator{
		strconv.Itoa(id),
		operKey,
		consPubKey,
	}
}

func (v *Validator) Edit(wg *sync.WaitGroup) {
	defer wg.Done()

	hosts := common.GlobalConfig.Networks[""].Hosts
	// pick a client randomly
	luckyNum := rand.Intn(len(hosts))
	config, _ := gosdk.NewClientConfig(
		hosts[luckyNum],
		common.GlobalConfig.Networks[""].ChainId,
		gosdk.BroadcastBlock,
		"0.1"+common.NativeToken,
		2000000,
		0,
		"",
	)
	cli := gosdk.NewClient(config)

	accInfo, err := cli.Auth().QueryAccount(v.operKey.GetAddress().String())
	if err != nil {
		log.Printf("failed. val %s query before edit: %s\n",
			sdk.ValAddress(v.operKey.GetAddress()).String(), err.Error())
		return
	}

	details := fmt.Sprintf("Time now: %v", time.Now())
	if _, err := cli.Staking().EditValidator(v.operKey, DefaultPasswd, DoNotModifyDesc, DoNotModifyDesc, DoNotModifyDesc, details,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence()); err != nil {
		log.Printf("failed. val %s edit: %s\n",
			sdk.ValAddress(v.operKey.GetAddress()).String(), err.Error())
	}
}

func (v *Validator) Create(wg *sync.WaitGroup) {
	defer wg.Done()

	hosts := common.GlobalConfig.Networks[""].Hosts
	// pick a client randomly
	luckyNum := rand.Intn(len(hosts))
	config, _ := gosdk.NewClientConfig(
		hosts[luckyNum],
		common.GlobalConfig.Networks[""].ChainId,
		gosdk.BroadcastBlock,
		"0.1"+common.NativeToken,
		2000000,
		0,
		"",
	)
	cli := gosdk.NewClient(config)

	accInfo, err := cli.Auth().QueryAccount(v.operKey.GetAddress().String())
	if err != nil {
		log.Printf("failed. val %s query before create: %s\n",
			sdk.ValAddress(v.operKey.GetAddress()).String(), err.Error())
		return
	}

	valAddrStr := sdk.ValAddress(v.operKey.GetAddress()).String()
	if _, err := cli.Staking().CreateValidator(v.operKey, DefaultPasswd, v.consPubkey, v.id, "", "", "",
		"", accInfo.GetAccountNumber(), accInfo.GetSequence()); err != nil {
		log.Printf("failed. val %s create: %s\n",
			valAddrStr, err.Error())
		return
	}
	log.Printf("[create validator successfully] %s\n", valAddrStr)
	CreateSuccessCnt++
}

func (v *Validator) Destroy(wg *sync.WaitGroup) {
	defer wg.Done()

	hosts := common.GlobalConfig.Networks[""].Hosts
	/// pick a client randomly
	luckyNum := rand.Intn(len(hosts))
	config, _ := gosdk.NewClientConfig(
		hosts[luckyNum],
		common.GlobalConfig.Networks[""].ChainId,
		gosdk.BroadcastBlock,
		"0.01"+common.NativeToken,
		200000,
		0,
		"",
	)
	cli := gosdk.NewClient(config)

	accInfo, err := cli.Auth().QueryAccount(v.operKey.GetAddress().String())
	if err != nil {
		log.Printf("failed. val %s query before destroy: %s\n",
			sdk.ValAddress(v.operKey.GetAddress()).String(), err.Error())
		return
	}

	valAddrStr := sdk.ValAddress(v.operKey.GetAddress()).String()
	if _, err := cli.Staking().DestroyValidator(v.operKey, DefaultPasswd, "", accInfo.GetAccountNumber(),
		accInfo.GetSequence()); err != nil {
		log.Printf("failed. val %s destroy: %s\n",
			valAddrStr, err.Error())
		return
	}
	log.Printf("[destroy validator successfully] %s\n", valAddrStr)
	DestroySuccessCnt++
}
