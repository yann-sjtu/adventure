package distribution

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/common/logger"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const withdrawRewards = common.WithdrawRewards

func WithdrawRewards(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		logger.PrintQueryAccountError(err, withdrawRewards, info)
		return
	}

	_, err = cli.Distribution().WithdrawRewards(info, common.PassWord, convertAccAddrToValAddr(info.GetAddress()),
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		logger.PrintExecuteTxError(err, withdrawRewards, info)
		return
	}
	logger.PrintExecuteTxSuccess(withdrawRewards, info)
}

func convertAccAddrToValAddr(accAddr types.AccAddress) string {
	return types.ValAddress(accAddr).String()
}
