package keeper

import (
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/top21_shares_control/types"
	"strconv"
	"strings"
)

type Keeper struct {
	cliManager     *common.ClientManager
	targetValAddrs []sdk.ValAddress
	workers        []mntcmn.Worker
	dominationPct  sdk.Dec
}

func NewKeeper() Keeper {
	return Keeper{}
}

func (k *Keeper) Init(configFilePath string) (err error) {
	// cli
	k.cliManager = common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	// params from toml
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
	for _, addrStr := range config.TargetValAddrs {
		accAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return err
		}

		valAddr := sdk.ValAddress(accAddr)
		k.targetValAddrs = append(k.targetValAddrs, valAddr)
	}

	// worker info
	for _, workerInfoStr := range config.WorkersAccInfo {
		strs := strings.Split(workerInfoStr, ",")
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

	return nil
}

//func (k *Keeper) AnalyseShares() (res types.AnalyseResult, err error) {
//	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
//	if err != nil {
//		return
//	}
//
//	targetTotal, globalTotal, bonedTotal := k.sumShares(vals)
//	log.Printf("target total: [%s]    boned total: [%s]    global total: [%s]\n",
//		targetTotal.String(), globalTotal.String(), bonedTotal.String())
//
//	// check validator number in top 21
//	if warning, valsToPromote := k.checkValNumInTop21(vals); warning {
//		strategy := k.genStrategyToPromoteValidators(valsToPromote, vals)
//		return types.NewAnalyseResult(1, strategy), nil
//	}
//
//	// check percent to dominate(gov vote)
//	if k.checkPercentToDominate(targetTotal, bonedTotal) {
//		return types.NewAnalyseResult(2, nil), nil
//	}
//
//	// check percent to plunder(distr reward)
//	if k.checkPercentToPlunder(vals, targetTotal, globalTotal) {
//		return types.NewAnalyseResult(3, nil), nil
//	}
//
//	return
//}
//
//func (k *Keeper) getPureWorker() (types.Worker, error) {
//	cli := k.cliManager.GetClient()
//	for _, worker := range k.workers {
//		delegator, err := cli.Staking().QueryDelegator(worker.GetAccAddr().String())
//		if err != nil {
//			return types.Worker{}, err
//		}
//
//		// pure worker
//		if len(delegator.ValidatorAddresses) == 0 {
//			fmt.Printf("\t worker [%s] is picked\n", worker.GetAccAddr().String())
//			return worker, nil
//		}
//	}
//
//	log.Println("Warning! There's no pure worker now")
//	return types.Worker{}, errors.New("warning! there's no pure worker now")
//}
//
//func (k *Keeper) addSharesToAllVals(worker types.Worker) error {
//	// whether worker deposit
//	if err := k.ensureWorkerDeposited(worker); err != nil {
//		return err
//	}
//
//	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(worker.GetAccAddr().String())
//	if err != nil {
//		return err
//	}
//
//	msg := types.NewMsgAddShares(accInfo.GetAccountNumber(), accInfo.GetSequence(), k.targetValAddrs, worker.GetAccAddr())
//	return mntcmn.SendMsg(mntcmn.Vote, msg, worker.GetIndex())
//}
//
//func (k *Keeper) ensureWorkerDeposited(worker types.Worker) error {
//	workerAddrStr := worker.GetAccAddr().String()
//	cli := k.cliManager.GetClient()
//	delegator, err := cli.Staking().QueryDelegator(workerAddrStr)
//	if err != nil {
//		return err
//	}
//
//	if !delegator.Tokens.Equal(sdk.ZeroDec()) {
//		// worker has deposited
//		fmt.Printf("worker [%s] has already deposited [%sokt]", workerAddrStr, delegator.Tokens)
//		return nil
//	}
//
//	// if worker not deposit
//	accInfo, err := cli.Auth().QueryAccount(workerAddrStr)
//	if err != nil {
//		return err
//	}
//
//	nativeTokenAmount := accInfo.GetCoins().AmountOf(common.NativeToken)
//	if nativeTokenAmount.LTE(sdk.OneDec()) {
//		return fmt.Errorf("worker [%s] has less than 1%s", workerAddrStr, common.NativeToken)
//	}
//
//	// depositAmount = balance - 1okt
//	depositAmount := nativeTokenAmount.Sub(sdk.OneDec())
//	return k.deposit(worker, sdk.NewDecCoinFromDec(common.NativeToken, depositAmount))
//}
//
//func (k *Keeper) deposit(worker types.Worker, amount sdk.SysCoin) error {
//	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(worker.GetAccAddr().String())
//	if err != nil {
//		return err
//	}
//
//	msg := types.NewMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), amount, worker.GetAccAddr())
//	if err := mntcmn.SendMsg(mntcmn.Staking, msg, worker.GetIndex()); err != nil {
//		return err
//	}
//
//	fmt.Printf("wait for [%s] to deposit [%s] ...\n", worker.GetAccAddr(), amount.String())
//	time.Sleep(constant.IntervalAfterTxBroadcast)
//	return nil
//}
