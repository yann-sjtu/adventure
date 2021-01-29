package keeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/final_top_21_control/types"
	"log"
	"strconv"
	"strings"
)

type Keeper struct {
	cliManager       *common.ClientManager
	targetValAddrs   []string
	targetValsFilter map[string]struct{}
	workers          []mntcmn.Worker
	dominationPct    sdk.Dec
	data             types.Data
}

func NewKeeper() Keeper {
	return Keeper{
		targetValsFilter: make(map[string]struct{}),
	}
}

//func (k *Keeper) InitRound() error {
//	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
//	if err != nil {
//		return err
//	}
//
//	k.data.Vals = vals
//	// sorts vals
//	sort.Sort(k.data.Vals)
//
//	enemyFilter, oursFilter := utils.BuildFilter(k.enemyValAddrs), utils.BuildFilter(k.targetValAddrs)
//	k.data.EnemyTotalShares, k.data.OurTotalShares = sdk.ZeroDec(), sdk.ZeroDec()
//	k.data.Top21SharesMap = make(map[string]sdk.Dec)
//	k.data.TargetValSharesMap = make(map[string]sdk.Dec)
//	var enemyCounter, oursCounter int
//	for i, val := range k.data.Vals {
//		valAddrStr := val.OperatorAddress.String()
//		// check whether target val
//		if _, ok := k.targetValsFilter[valAddrStr]; ok {
//			k.data.TargetValSharesMap[valAddrStr] = val.DelegatorShares
//		}
//
//		// top 21 vals
//		if i < 21 {
//			k.data.Top21SharesMap[valAddrStr] = val.DelegatorShares
//			if _, ok := oursFilter[valAddrStr]; ok {
//				k.data.OurTotalShares = k.data.OurTotalShares.Add(val.DelegatorShares)
//				oursCounter++
//				continue
//			}
//
//			if _, ok := enemyFilter[valAddrStr]; ok {
//				if enemyCounter == 0 {
//					k.data.EnemyLowestShares = val.DelegatorShares
//				} else if val.DelegatorShares.LT(k.data.EnemyLowestShares) {
//					k.data.EnemyLowestShares = val.DelegatorShares
//				}
//
//				k.data.EnemyTotalShares = k.data.EnemyTotalShares.Add(val.DelegatorShares)
//				enemyCounter++
//			}
//		}
//	}
//
//	if enemyCounter+oursCounter != 21 {
//		log.Println("Warning! the sum of enemies and ours is not 21.")
//	}
//
//	return nil
//}

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
	for _, valAddrStr := range k.targetValAddrs {
		// add to filter
		k.targetValsFilter[valAddrStr] = struct{}{}
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
	if len(k.targetValAddrs) != len(k.targetValsFilter) {
		log.Panicf("different length with targetValAddrs and targetValsFilter\n")
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
//func (k *Keeper) CatchTheIntruders() []string {
//	// build top21Filter
//	top21Filter := make(map[string]struct{})
//	for valAddrStr := range k.data.Top21SharesMap {
//		top21Filter[valAddrStr] = struct{}{}
//	}
//
//	for _, targetValAddrStr := range k.targetValAddrs {
//		delete(top21Filter, targetValAddrStr)
//	}
//
//	for _, enemyValAddrStr := range k.enemyValAddrs {
//		delete(top21Filter, enemyValAddrStr)
//	}
//
//	var intruders []string
//	for intruder := range top21Filter {
//		intruders = append(intruders, intruder)
//	}
//
//	if len(intruders) != 0 {
//		log.Printf("WARNING! instrders %s are found\n", intruders)
//	}
//
//	return intruders
//}
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
