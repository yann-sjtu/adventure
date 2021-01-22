package keeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/monitor/shares-control/types"
	"path/filepath"
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

func (k *Keeper) Init(configFilePath string) error {
	// cli
	k.cliManager = common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	// params from toml
	var config types.Config
	path, err := filepath.Abs(filepath.Dir("config.toml"))
	if err!=nil{
		return err
	}

	fmt.Println(path)
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		return err
	}

	fmt.Println(config)

	return nil
}
