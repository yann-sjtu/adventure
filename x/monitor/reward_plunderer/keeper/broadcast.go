package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/reward_plunderer/constant"
	"time"
)

func (k *Keeper) InfoToAddShares(worker mntcmn.Worker, valAddrs []sdk.ValAddress) error {
	workerAddr := worker.GetAccAddr().String()
	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(workerAddr)
	if err != nil {
		return fmt.Errorf("worker [%s] query account failed: %s", workerAddr, err.Error())
	}

	signMsg := mntcmn.NewMsgAddShares(accInfo.GetAccountNumber(), accInfo.GetSequence(), valAddrs, worker.GetAccAddr(), 1050000)
	err = mntcmn.SendMsg(mntcmn.Vote, signMsg, worker.GetIndex())
	if err != nil {
		return err
	}
	fmt.Printf("worker %s has informed to add shares to %d validators successfully!\n", workerAddr, len(valAddrs))
	time.Sleep(constant.IntervalAfterTxBroadcast)
	return nil
}

func (k *Keeper) InfoToDeposit(worker mntcmn.Worker, coin sdk.SysCoin) error {
	workerAddr := worker.GetAccAddr().String()
	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(workerAddr)
	if err != nil {
		return fmt.Errorf("worker [%s] query account failed: %s", workerAddr, err.Error())
	}

	signMsg := mntcmn.NewMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), coin, worker.GetAccAddr(), 1050000)
	err = mntcmn.SendMsg(mntcmn.Staking, signMsg, worker.GetIndex())
	if err != nil {
		return err
	}
	fmt.Printf("worker %s has informed to depoist %s successfully!\n", workerAddr, coin.String())
	time.Sleep(constant.IntervalAfterTxBroadcast)
	return nil
}
