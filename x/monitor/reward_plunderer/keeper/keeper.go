package keeper

import (
	"errors"
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
	cliManager        *common.ClientManager
	ourValAddrs       []string
	ourValAddrsFilter map[string]struct{}
	ourTop18ValAddrs  []sdk.ValAddress
	workers           []mntcmn.Worker
	plunderedPct      sdk.Dec
	data              types.Data
}

func NewKeeper() Keeper {
	return Keeper{
		ourValAddrsFilter: make(map[string]struct{}),
	}
}

func (k *Keeper) InitRound() error {
	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
	if err != nil {
		return err
	}

	k.data.Vals = vals
	// sorts vals
	sort.Sort(k.data.Vals)

	// sum total shares of all our validators/ global validators' total shares
	k.data.OurTotalShares, k.data.AllValsTotalShares = sdk.ZeroDec(), sdk.ZeroDec()
	for _, val := range k.data.Vals {
		k.data.AllValsTotalShares = k.data.AllValsTotalShares.Add(val.DelegatorShares)
		if _, ok := k.ourValAddrsFilter[val.OperatorAddress.String()]; ok {
			k.data.OurTotalShares = k.data.OurTotalShares.Add(val.DelegatorShares)
		}
	}

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

	// parse top 18 val addrs
	for _, valAddrStr := range config.OurTop18Addrs {
		valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			return err
		}

		k.ourTop18ValAddrs = append(k.ourTop18ValAddrs, valAddr)
	}

	// our val addr
	for _, accAddrStr := range config.OurValAddrs {
		accAddr, err := sdk.AccAddressFromBech32(accAddrStr)
		if err != nil {
			return err
		}
		valAddrStr := sdk.ValAddress(accAddr).String()
		k.ourValAddrs = append(k.ourValAddrs, valAddrStr)
		// add filter
		k.ourValAddrsFilter[valAddrStr] = struct{}{}
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

func (k *Keeper) PickEfficientWorker(tokenToDeposit sdk.SysCoin) (worker mntcmn.Worker, err error) {
	cli := k.cliManager.GetClient()
	for _, w := range k.workers {
		workerAddr := w.GetAccAddr().String()
		accInfo, err := cli.Auth().QueryAccount(workerAddr)
		if err != nil {
			return worker, fmt.Errorf("worker [%s] query account failed: %s", workerAddr, err.Error())
		}

		balance := accInfo.GetCoins().AmountOf(common.NativeToken)
		if balance.Sub(constant.ReservedFee).GTE(tokenToDeposit.Amount) {
			log.Printf("worker [%s] will deposit [%s] for our top 18 validators\n", workerAddr, tokenToDeposit.String())
			return w, nil
		}
		time.Sleep(constant.QueryInverval)
	}

	err = errors.New("no efficient worker already")
	return
}

func (k *Keeper) GetOurTop18ValAddrs() []sdk.ValAddress {
	return k.ourTop18ValAddrs
}

func (k *Keeper) GetWorkers() []mntcmn.Worker {
	return k.workers
}

func (k *Keeper) GetCliManager() *common.ClientManager {
	return k.cliManager
}
