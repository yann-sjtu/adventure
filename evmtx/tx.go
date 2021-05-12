package evmtx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
)

var (
	defaultGasPrice, _ = sdk.ParseDecCoin("0.000000001okt")
)

func generateTxParams(from string) []map[string]string {
	param := make([]map[string]string, 1)
	param[0] = map[string]string{
		"from":     from,
		"to":       "0xEe666E967293094007d7C3718044e07565B1f8a9", // todo
		"value":    "0x0",
		"gasPrice": (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String(),
		"data":     "0x"+ethcommon.Bytes2Hex(UniswapV2.BuildApprovePayload("0x1bF1D0fF0dBE2A3ffD4E287033E3354e9Ad5Ea80", 100)),
	}

	return param
}
