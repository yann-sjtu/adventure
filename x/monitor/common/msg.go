package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/adventure/common"
	stakingtypes "github.com/okex/exchain/x/staking/types"
)

var gasPrice sdk.Dec

func init() {
	gasPrice = sdk.NewDecWithPrec(1, 9)
}

func NewMsgAddShares(accNum uint64, seqNum uint64, valAddrs []sdk.ValAddress, accAddr sdk.AccAddress, gas uint64) authtypes.StdSignMsg {
	msg := stakingtypes.NewMsgAddShares(accAddr, valAddrs)
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          []sdk.Msg{msg},
		Fee:           authtypes.NewStdFee(gas, sdk.NewDecCoinsFromDec(common.NativeToken, gasPrice.MulInt64(int64(gas)))),
	}

	return signMsg
}

func NewMsgDeposit(accNum, seqNum uint64, amount sdk.SysCoin, accAddr sdk.AccAddress, gas uint64) authtypes.StdSignMsg {
	msg := stakingtypes.NewMsgDeposit(accAddr, amount)
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          []sdk.Msg{msg},
		Fee:           authtypes.NewStdFee(gas, sdk.NewDecCoinsFromDec(common.NativeToken, gasPrice.MulInt64(int64(gas)))),
	}

	return signMsg
}
