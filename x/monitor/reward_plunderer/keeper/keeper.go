package keeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/reward_plunderer/constant"
	"github.com/okex/adventure/x/monitor/reward_plunderer/types"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Keeper struct {
	cliManager       *common.ClientManager
	ourValAddrs      []string
	ourTop18ValAddrs []string
	workers          []mntcmn.Worker
	plunderedPct     sdk.Dec
	data             types.Data
}

func NewKeeper() Keeper {
	return Keeper{}
}

func (k *Keeper) InitRound() error {
	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
	if err != nil {
		return err
	}

	k.data.Vals = vals
	// sorts vals
	sort.Sort(k.data.Vals)
	return nil
}

func (k *Keeper) Init(configFilePath string) (err error) {
	// cli
	k.cliManager = common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	// config from toml
	var config types.Config
	if _, err = toml.DecodeFile(configFilePath, &config); err != nil {
		return
	}

	if err = k.parseConfig(&config); err != nil {
		return
	}

	k.logInit()
	return nil
}

func (k *Keeper) parseConfig(config *types.Config) error {
	// decimal
	percentToPlunder, err := sdk.NewDecFromStr(config.PercentToPlunder)
	if err != nil {
		return err
	}
	k.plunderedPct = percentToPlunder

	k.ourTop18ValAddrs = config.OurTop18Addrs
	// our val addr
	for _, accAddrStr := range config.OurValAddrs {
		accAddr, err := sdk.AccAddressFromBech32(accAddrStr)
		if err != nil {
			return err
		}

		k.ourValAddrs = append(k.ourValAddrs, sdk.ValAddress(accAddr).String())
	}

	// worker info
	for _, workerInfoStr := range config.WorkersAccInfo {
		strs := strings.Split(workerInfoStr, ",")
		if len(strs) != 2 {
			return fmt.Errorf("length of item in config.workers_infos is not 2")
		}

		accAddr, err := sdk.AccAddressFromBech32(strs[0])
		if err != nil {
			return err
		}

		index, err := strconv.Atoi(strs[1])
		if err != nil {
			return err
		}

		k.workers = append(k.workers, mntcmn.NewWorker(accAddr, index))
	}

	// sanity check
	if len(k.ourValAddrs) != 31 {
		log.Panicf("length of ourValAddrs is not 31\n")
	}

	if len(k.ourTop18ValAddrs) != 18 {
		log.Panicf("length of ourTop18ValAddrs is not 18\n")
	}

	return nil
}

func (k *Keeper) SendMsgs(worker mntcmn.Worker, coin sdk.DecCoin) error {
	workerAddr := worker.GetAccAddr().String()
	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(workerAddr)
	if err != nil {
		return fmt.Errorf("worker [%s] query account failed: %s", workerAddr, err.Error())
	}

	signMsg := mntcmn.NewMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), coin, worker.GetAccAddr())
	err = mntcmn.SendMsg(mntcmn.Staking, signMsg, worker.GetIndex())
	if err != nil {
		return err
	}
	fmt.Printf("worker %s has informed to depoist %s successfully!\n", workerAddr, coin.String())
	time.Sleep(constant.IntervalAfterTxBroadcast)
	return nil
}
