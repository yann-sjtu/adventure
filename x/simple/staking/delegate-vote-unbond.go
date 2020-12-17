package staking

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/module/auth/types"
	stakingTypes "github.com/okex/okexchain-go-sdk/module/staking/types"
)

const (
	delegate = "delegate"
	vote     = "vote"
	unbond   = "unbond"
)

func DelegateVoteUnbond(cli *gosdk.Client, info keys.Info) {
	sleepTime := config.Cfg.Staking.DelegateVoteUnbondConfig.SleepTime
	//send delegate tx
	sendTx(cli, info, delegate)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	//send vote tx
	for k := 0; k < 3; k++ {
		sendTx(cli, info, vote)
	}
	time.Sleep(time.Duration(sleepTime) * time.Second)
	//send unbond tx
	sendTx(cli, info, unbond)
	time.Sleep(time.Duration(sleepTime) * time.Second)
}

func sendTx(cli *gosdk.Client, info keys.Info, phase string) {
	var err error
	var accInfo types.Account
	addr := info.GetAddress().String()
	accInfo, err = cli.Auth().QueryAccount(addr)
	if err != nil {
		logger.PrintQueryAccountError(err, phase, info)
		return
	}

	passWd := common.PassWord
	switch phase {
	case delegate:
		delegateNum := config.Cfg.Staking.DelegateVoteUnbondConfig.DelegateNum
		_, err = cli.Staking().Deposit(info, passWd, delegateNum, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case vote:
		valAddrs := getValditorAddrs(cli)
		_, err = cli.Staking().AddShares(info, passWd, valAddrs, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case unbond:
		unbondNum := config.Cfg.Staking.DelegateVoteUnbondConfig.UnbondNum
		_, err = cli.Staking().Withdraw(info, passWd, unbondNum, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case regProxy:
		_, err = cli.Staking().RegisterProxy(info, passWd, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case unregProxy:
		_, err = cli.Staking().UnregisterProxy(info, passWd, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case unbindProxy:
		_, err = cli.Staking().UnbindProxy(info, passWd, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	}

	if err != nil {
		logger.PrintExecuteTxError(err, phase, info)
		return
	}
	logger.PrintExecuteTxSuccess(phase, info)
}

func getValditorAddrs(cli *gosdk.Client) []string {
	oldVals, err := cli.Staking().QueryValidators()
	if err != nil {
		return []string{""}
	}

	vals := shuffle(oldVals)
	num := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(30) + 1
	if len(vals) < num {
		num = len(vals)
	}

	var valAddrs []string
	for _, v := range vals {
		if !(v.MinSelfDelegation.IsZero() || v.Jailed) {
			valAddrs = append(valAddrs, v.OperatorAddress.String())
		}
		if len(valAddrs) >= num {
			return valAddrs[:num]
		}
	}
	return valAddrs[:num]
}

func shuffle(vals []stakingTypes.Validator) []stakingTypes.Validator {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]stakingTypes.Validator, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}
