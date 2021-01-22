package keeper

import (
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/monitor/shares-control/types"
	"strconv"
	"strings"
)

type Keeper struct {
	cliManager     *common.ClientManager
	params         types.Params
	targetValAddrs []sdk.ValAddress
	workers        []types.Worker
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
	k.params = types.NewParams(valNumberInTop21, percentToPlunder)

	// val addr
	for _, addrStr := range config.TargetValAddrs {
		accAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return err
		}

		k.targetValAddrs = append(k.targetValAddrs, sdk.ValAddress(accAddr))
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
