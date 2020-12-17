package distribution

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const setWithdrawAddr = common.SetWithdrawAddr

func SetWithdrawAddr(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, setWithdrawAddr, info)
		return
	}

	addrs := config.Cfg.Distribution.SetWithdrawAddrConfig.Address
	addr := addrs[rand.Intn(len(addrs))]
	_, err = cli.Distribution().SetWithdrawAddr(info, common.PassWord,
		addr,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, setWithdrawAddr, info)
		return
	}
	common.PrintExecuteTxSuccess(setWithdrawAddr, info)
}
