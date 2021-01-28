package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mntcmn "github.com/okex/adventure/x/monitor/common"
	"github.com/okex/adventure/x/monitor/top21_shares_control/constant"
	"time"
)

func (k *Keeper) deposit(worker mntcmn.Worker, amount sdk.SysCoin) error {
	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(worker.GetAccAddr().String())
	if err != nil {
		return err
	}

	msg := mntcmn.NewMsgDeposit(accInfo.GetAccountNumber(), accInfo.GetSequence(), amount, worker.GetAccAddr())
	if err := mntcmn.SendMsg(mntcmn.Staking, msg, worker.GetIndex()); err != nil {
		return err
	}

	fmt.Printf("wait for [%s] to deposit [%s] ...\n", worker.GetAccAddr(), amount.String())
	time.Sleep(constant.IntervalAfterTxBroadcast)
	return nil
}

func (k *Keeper) addShares(worker mntcmn.Worker, targetValAddrsStr []string) error {
	var valAddrs []sdk.ValAddress
	for _, valAddrStr := range targetValAddrsStr {
		valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			return err
		}

		valAddrs = append(valAddrs, valAddr)
	}

	accInfo, err := k.cliManager.GetClient().Auth().QueryAccount(worker.GetAccAddr().String())
	if err != nil {
		return err
	}

	msg := mntcmn.NewMsgAddShares(accInfo.GetAccountNumber(), accInfo.GetSequence(), valAddrs, worker.GetAccAddr())
	if err := mntcmn.SendMsg(mntcmn.Staking, msg, worker.GetIndex()); err != nil {
		return err
	}

	fmt.Printf("wait for [%s] to add shares to [%s] ...\n", worker.GetAccAddr(), targetValAddrsStr)
	time.Sleep(constant.IntervalAfterTxBroadcast)
	return nil
}
