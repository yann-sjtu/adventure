package keeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/final_top_21_control/types"
	"log"
	"sort"
	"strconv"
	"strings"
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

//
//// get targetAddrs with enemyAddrs filtered in bonded vals
//// addrType:   1-accAddr, 2-valAddr
//func (k *Keeper) GetTargetValsAddr(enemyAddrs []string, addrType int) (targetAddrs []string, err error) {
//	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
//	if err != nil {
//		return
//	}
//
//	filter := utils.BuildFilter(enemyAddrs)
//	for _, val := range vals {
//		if val.Status.Equal(sdk.Bonded) {
//			var addr string
//			switch addrType {
//			case 1:
//				addr = sdk.AccAddress(val.OperatorAddress).String()
//			case 2:
//				addr = val.OperatorAddress.String()
//			default:
//				return nil, fmt.Errorf("unsupported input addr type %d", addrType)
//			}
//
//			if _, ok := filter[addr]; !ok {
//				targetAddrs = append(targetAddrs, addr)
//			}
//		}
//	}
//
//	return
//}
//
//func (k *Keeper) GetEnemyValAddrs() []string {
//	return k.enemyValAddrs
//}
//
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
		log.Printf("WARNING! instrders %s are found\n", intruders)
	}

	return intruders
}

//
//func (k *Keeper) PrecheckWorker(workers []mntcmn.Worker, tokenToDeposit sdk.SysCoin) (selected []mntcmn.Worker, err error) {
//	// 1.check the balance
//	cli := k.cliManager.GetClient()
//	for _, worker := range workers {
//		workerAddr := worker.GetAccAddr().String()
//		accInfo, err := cli.Auth().QueryAccount(workerAddr)
//		if err != nil {
//			return nil, err
//		}
//
//		// tokenToDepositAmount <= balance - reservedFee
//		balance := accInfo.GetCoins().AmountOf(common.NativeToken)
//		if balance.Sub(constant.ReservedFee).LT(tokenToDeposit.Amount) {
//			log.Printf("insufficient balance of %s: %s < %s + %s. skip that.\n", workerAddr, balance.String(), tokenToDeposit.Amount.String(), constant.ReservedFee.String())
//			continue
//		}
//
//		selected = append(selected, worker)
//	}
//
//	return
//}
