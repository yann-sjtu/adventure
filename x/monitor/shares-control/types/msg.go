package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/adventure/common"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
)

func NewMsgAddShares(accNum uint64, seqNum uint64, valAddrs []sdk.ValAddress, accAddr sdk.AccAddress) authtypes.StdSignMsg {
	msg := stakingtypes.NewMsgAddShares(accAddr, valAddrs)
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          []sdk.Msg{msg},
		Fee:           authtypes.NewStdFee(500000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(5, 3))),
	}

	return signMsg
}

func NewMsgDeposit(accNum, seqNum uint64, amount sdk.SysCoin, accAddr sdk.AccAddress) authtypes.StdSignMsg {
	msg := stakingtypes.NewMsgDeposit(accAddr, amount)
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          []sdk.Msg{msg},
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}
