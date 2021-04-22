package types

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/staking/validators/val"
	"github.com/okex/exchain-go-sdk"
	stakingTypes "github.com/okex/exchain-go-sdk/module/staking/types"
	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/utils"
)

var (
	// Singleton Pattern
	once       sync.Once
	valManager ValManager
)

type ValManager map[string]val.Validator

func GetValManager() ValManager {
	once.Do(func() {
		// initialize the valManager

		// 1.get operator keys
		valAccounts, err := GetTestAccountsFromFile("./template/mnemonic/val_testnet_40")
		//valAccounts, err := GetTestAccountsFromFile("./cmd/advanture/staking/test/val_mnemonics_local.txt")
		if err != nil {
			panic(err)
		}

		valsLen := len(valAccounts)
		valManager = make(ValManager, valsLen)

		// 2.build validators
		for i, valAcc := range valAccounts {
			valManager[sdk.ValAddress(valAcc.GetAddress()).String()] = NewValidator(valAcc, i, ConsPubkeys[i])
		}
	})

	return valManager
}

// GetValidators gets all validators
func (vm *ValManager) GetValidators() ([]stakingTypes.Validator, error) {
	// pick a client randomly
	hosts := common.Cfg.Hosts
	luckyNum := rand.Intn(len(hosts))
	cli := gosdk.NewClient(types.ClientConfig{
		NodeURI:       hosts[luckyNum],
		BroadcastMode: gosdk.BroadcastBlock,
	})

	return cli.Staking().QueryValidators()
}

// GetIncumbentVals gets all undestroyed validator addresses
func (vm *ValManager) GetIncumbentValAddrs() (incumbentValAddrs []string, err error) {
	vals, err := vm.GetValidators()
	if err != nil {
		return
	}

	for _, v := range vals {
		if !(v.MinSelfDelegation.IsZero() || v.Jailed) {
			incumbentValAddrs = append(incumbentValAddrs, v.OperatorAddress.String())
		}
	}

	return
}

func GetTestAccountsFromFile(path string) ([]keys.Info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	var index int
	var accounts []keys.Info
	for {
		mnemonic, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		acc, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(mnemonic),
			fmt.Sprintf("acc%d", index), DefaultPasswd)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
		fmt.Println(accounts[index].GetAddress().String(), index)
		index++
	}

	return accounts, nil
}
