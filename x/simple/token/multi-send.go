package token

import (
	"math/rand"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/config"
	gosdk "github.com/okex/okexchain-go-sdk"
	tokenTypes "github.com/okex/okexchain-go-sdk/module/token/types"
)

const (
	multiSendCoin = "0.001"
)

var (
	once     sync.Once
	accAddrs []types.AccAddress
)

func MultiSend(cli *gosdk.Client, info keys.Info) {
	once.Do(func() {
		addrs := common.GetAccountAddressFromFile(config.Cfg.Token.MultiSendConfig.ToAddrsPath)
		for _, addr := range addrs {
			accAddr, _ := types.AccAddressFromBech32(addr)
			accAddrs = append(accAddrs, accAddr)
		}
	})

	acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil || len(acc.GetCoins()) == 0 {
		common.PrintQueryTokensError(err, common.MultiSend, info)
		return
	}

	coins := acc.GetCoins()
	symbol := coins[rand.Intn(len(coins))].Denom
	topUp(info, accAddrs, multiSendCoin+symbol, cli)
}

func topUp(rich keys.Info, accAddrs []types.AccAddress, coinStr string, cli *gosdk.Client) {
	transferUnit, err := makeTransferUnit(accAddrs, coinStr)
	if err != nil {
		return
	}

	accInfo, err := cli.Auth().QueryAccount(rich.GetAddress().String())
	if err != nil {
		return
	}

	// assume a successful transfer
	_, err = cli.Token().MultiSend(rich, common.PassWord,
		transferUnit,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, common.MultiSend, rich)
		return
	}
	common.PrintExecuteTxSuccess(common.MultiSend, rich)

}

func makeTransferUnit(accAddrs []types.AccAddress, coinStr string) ([]tokenTypes.TransferUnit, error) {
	coins, err := types.ParseDecCoins(coinStr)
	if err != nil {
		return nil, err
	}

	accLen := len(accAddrs)
	transferUnits := make([]tokenTypes.TransferUnit, accLen)
	for i := 0; i < accLen; i++ {
		transferUnits[i] = tokenTypes.TransferUnit{To: accAddrs[i], Coins: coins}
	}

	return transferUnits, nil
}
