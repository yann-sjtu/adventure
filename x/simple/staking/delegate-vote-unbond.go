package staking

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/module/auth/types"
	stakingTypes "github.com/okex/okexchain-go-sdk/module/staking/types"
)

const (
	delegate = "delegate"
	vote     = "vote"
	unbond   = "unbond"

	sleepTime   = 3
	delegateNum = "0.01" + common.NativeToken
	unbondNum   = "0.005" + common.NativeToken

	passWd = common.PassWord
)

func DelegateVoteUnbond(cli *gosdk.Client, info keys.Info) {
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
		common.PrintQueryAccountError(err, phase, info)
		return
	}

	switch phase {
	case delegate:
		_, err = cli.Staking().Deposit(info, passWd, delegateNum, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case vote:
		valAddrs := getValditorAddrs(cli)
		_, err = cli.Staking().AddShares(info, passWd, valAddrs, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case unbond:
		_, err = cli.Staking().Withdraw(info, passWd, unbondNum, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case regProxy:
		_, err = cli.Staking().RegisterProxy(info, passWd, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case unregProxy:
		_, err = cli.Staking().UnregisterProxy(info, passWd, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	case unbindProxy:
		_, err = cli.Staking().UnbindProxy(info, passWd, "", accInfo.GetAccountNumber(), accInfo.GetSequence())
	}

	if err != nil {
		common.PrintExecuteTxError(err, phase, info)
		return
	}
	common.PrintExecuteTxSuccess(phase, info)
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
