package keeper

import (
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/monitor/shares-control/types"
	"log"
	"strconv"
	"strings"
)

type Keeper struct {
	cliManager       *common.ClientManager
	expectedParams   types.Params
	targetValAddrs   []sdk.ValAddress
	targetValsFilter map[string]struct{}
	workers          []types.Worker
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
	valNumberInTop21 := sdk.NewDec(int64(len(config.TargetValAddrs)))
	percentToPlunder, err := sdk.NewDecFromStr(config.PercentToPlunder)
	if err != nil {
		return err
	}
	percentToDominate, err := sdk.NewDecFromStr(config.PercentToDominate)
	if err != nil {
		return err
	}

	k.expectedParams = types.NewParams(valNumberInTop21, percentToPlunder, percentToDominate)

	// val addr
	k.targetValsFilter = make(map[string]struct{})
	for _, addrStr := range config.TargetValAddrs {
		accAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return err
		}

		valAddr := sdk.ValAddress(accAddr)
		k.targetValAddrs = append(k.targetValAddrs, valAddr)
		// add filter
		k.targetValsFilter[valAddr.String()] = struct{}{}
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

		k.workers = append(k.workers, types.NewWorker(accAddr, index))
	}

	return nil
}

func (k *Keeper) AnalyseShares() (res types.AnalyseResult, err error) {
	vals, err := k.cliManager.GetClient().Staking().QueryValidators()
	if err != nil {
		return
	}

	targetTotal, globalTotal, bonedTotal := k.sumShares(vals)
	log.Printf("target total: [%s]    boned total: [%s]    global total: [%s]\n",
		targetTotal.String(), globalTotal.String(), bonedTotal.String())

	// check validator number in top 21
	if warning, valToPromote := k.checkValNumInTop21(vals); warning {
		return types.NewAnalyseResult(1, valToPromote), nil
	}

	// check percent to dominate(gov vote)
	if k.checkPercentToDominate(targetTotal, bonedTotal) {
		return types.NewAnalyseResult(2, nil), nil
	}

	// check percent to plunder(distr reward)
	if k.checkPercentToPlunder(vals, targetTotal, globalTotal) {
		return types.NewAnalyseResult(3, nil), nil
	}

	return
}
