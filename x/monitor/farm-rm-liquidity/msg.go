package farm_rm_liquidity

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/adventure/common"
	"github.com/okex/okexchain-go-sdk/utils"
	ammswaptypes "github.com/okex/okexchain/x/ammswap/types"
	"time"
)

func newMsgRemoveLiquidity(accNum, seqNum uint64, liquidity sdk.Dec, minBaseAmount, minQuoteAmount sdk.SysCoin, deadline int64, addr string) authtypes.StdSignMsg {
	cosmosAddr, err := utils.ToCosmosAddress(addr)
	if err != nil {
		panic(err)
	}

	msg := ammswaptypes.NewMsgRemoveLiquidity(liquidity, minBaseAmount, minQuoteAmount, deadline, cosmosAddr)
	signMsg := authtypes.StdSignMsg{
		ChainID:       "okexchain-66",
		AccountNumber: accNum,
		Sequence:      seqNum,
		Memo:          "",
		Msgs:          []types.Msg{msg},
		Fee:           authtypes.NewStdFee(200000, sdk.NewDecCoinsFromDec(common.NativeToken, sdk.NewDecWithPrec(2, 4))),
	}
	return signMsg
}

func getDeadline() int64 {
	duration, err := time.ParseDuration("24h")
	if err != nil {
		panic(fmt.Errorf("this should never happen: %s", err))
	}
	deadline := time.Now().Add(duration).Unix()
	return deadline
}
