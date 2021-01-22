package farm_control

import (
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/adventure/common"
	"github.com/okex/okexchain-go-sdk/utils"
	ammswaptypes "github.com/okex/okexchain/x/ammswap/types"
	farmtypes "github.com/okex/okexchain/x/farm/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
)

func newMsgDeposit(accNum, seqNum uint64, amount sdk.SysCoin, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := stakingtypes.NewMsgDeposit(cosmosAddr, amount)
	msgs := []types.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}

func newMsgLock(accNum, seqNum uint64, amount sdk.SysCoin, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := farmtypes.NewMsgLock(poolName, cosmosAddr, amount)
	msgs := []types.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}

func newMsgUnLock(accNum, seqNum uint64, amount sdk.SysCoin, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := farmtypes.NewMsgUnlock(poolName, cosmosAddr, amount)
	msgs := []types.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}

	return signMsg
}

func newMsgAddLiquidity(accNum, seqNum uint64, minLiquidity sdk.Dec, maxBaseAmount, quoteAmount sdk.SysCoin, deadline int64, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := ammswaptypes.NewMsgAddLiquidity(minLiquidity, maxBaseAmount, quoteAmount, deadline, cosmosAddr)
	msgs := []types.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}
	return signMsg
}

func newMsgRemoveLiquidity(accNum, seqNum uint64, liquidity sdk.Dec, minBaseAmount, minQuoteAmount sdk.SysCoin, deadline int64, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := ammswaptypes.NewMsgRemoveLiquidity(liquidity, minBaseAmount, minQuoteAmount, deadline, cosmosAddr)
	msgs := []types.Msg{msg}
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          msgs,
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 3))),
	}
	return signMsg
}
