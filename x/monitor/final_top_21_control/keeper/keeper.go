package keeper

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/final_top_21_control/constant"
	"github.com/okex/adventure/x/monitor/final_top_21_control/types"
)

type Keeper struct {
	cliManager     *common.ClientManager
	targetValAddrs []string
	workers        []mntcmn.Worker
	dominationPct  sdk.Dec
	data           types.Data
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
	percentToDominate, err := sdk.NewDecFromStr(config.PercentToDominate)
	if err != nil {
		return err
	}
	k.dominationPct = percentToDominate

	// val addr
	k.targetValAddrs = config.TargetValAddrs

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
	if len(k.targetValAddrs) != 21 {
		log.Panicf("length of targetValAddrs is not 21\n")
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
			log.Printf("worker [%s] will deposit [%s] for all target validators\n", workerAddr, tokenToDeposit.String())
			return w, nil
		}
		time.Sleep(time.Second * 3)
	}

	err = errors.New("no efficient worker already")
	return
}

func (k *Keeper) CatchTheIntruders() []string {
	// build top21Filter
	top21Filter := make(map[string]struct{})
	for i := 0; i < 21; i++ {
		top21Filter[k.data.Vals[i].OperatorAddress.String()] = struct{}{}
	}

	for _, tarValAddrStr := range k.targetValAddrs {
		delete(top21Filter, tarValAddrStr)
	}

	var intruders []string
	for k := range top21Filter {
		intruders = append(intruders, k)
	}

	if len(intruders) != 0 {
		log.Printf("WARNING! instruders %s are found\n", intruders)
	}

	return intruders
}

func (k *Keeper) SendMsgs(worker mntcmn.Worker, coin sdk.DecCoin) error {
	workerAddr := worker.GetAccAddr().String()
	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(workerAddr)
	if err != nil {
		return fmt.Errorf("worker [%s] query account failed: %s", workerAddr, err.Error())
	}

	signMsg := mntcmn.NewMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), coin, worker.GetAccAddr(),1100000)
	err = mntcmn.SendMsg(mntcmn.Staking, signMsg, worker.GetIndex())
	if err != nil {
		return err
	}
	fmt.Printf("worker %s has informed to depoist %s successfully!\n", workerAddr, coin.String())
	time.Sleep(constant.IntervalAfterTxBroadcast)
	return nil
}
