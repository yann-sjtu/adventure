package common

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/okexchain/x/ammswap"
	farmtypes "github.com/okex/okexchain/x/farm/types"
	"github.com/okex/okexchain/x/staking"
)

type MsgWithIndex struct {
	types.StdSignMsg
	AddressIndex int `json:"addressIndex" yaml:"addressIndex"`
}

func NewMsgWithIndex(msg types.StdSignMsg, index int) MsgWithIndex {
	return MsgWithIndex{
		msg,
		index,
	}
}

const coinType = 1571

// {"coinType":1571,"hex":"","relatedId":"1","txType":1}
type Object struct {
	CoinType  int    `json:"coinType"`
	Hex       string `json:"hex"`
	RelatedId string `json:"relatedId"`
	TxType    int    `json:"txType"`
}

var (
	cdc = codec.New()
)

func init() {
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	farmtypes.RegisterCodec(cdc)
	ammswap.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
}

func NewObject(msg types.StdSignMsg, addresIndex int, time string, txType int) (Object, error) {
	msgWithIndex := NewMsgWithIndex(msg, addresIndex)
	msgWithIndexBins, err := cdc.MarshalJSON(msgWithIndex)
	if err != nil {
		return Object{}, err
	}
	sortedBins := sdk.MustSortJSON(msgWithIndexBins)
	str := string(sortedBins)
	hexStr := resolveMsgStr(str)
	//{"StdSignMsg":{"chain_id":"okexchain-66","account_number":"1238","sequence":"0","fee":{"amount":[{"denom":"okt","amount":"0.002000000000000000"}],"gas":"200000"},"msgs":[{"type":"okexchain/farm/MsgLock","value":{"pool_name":"1st_pool_okt_usdt","address":"okexchain1ln38mfpx5vuugk85grljw8c4utcechdm3v55xp","amount":{"denom":"ammswap_okt_usdt-a2b","amount":"10000.000000000000000000"}}}],"memo":""},"addressIndex":"925"}
	//{"chain_id":"okexchain-66","account_number":"1238","sequence":"0","fee":{"amount":[{"denom":"okt","amount":"0.002000000000000000"}],"gas":"200000"},"msgs":[{"type":"okexchain/farm/MsgLock","value":{"pool_name":"1st_pool_okt_usdt","address":"okexchain1ln38mfpx5vuugk85grljw8c4utcechdm3v55xp","amount":{"denom":"ammswap_okt_usdt-a2b","amount":"10000.000000000000000000"}}}],"memo":""},"addressIndex":"925"}
	//fmt.Println(ss)

	return Object{
		CoinType:  coinType,
		Hex:       hexStr,
		RelatedId: time,
		TxType:    txType,
	}, nil
}

func resolveMsgStr(msg string) string {
	prefixLength := len("\"StdSignMsg\":") // {"StdSignMsg":
	tmpMsg := msg[prefixLength+1:]

	lastMatched := "},\"addressIndex" // },"addressIndex
	lastIndex := strings.Index(tmpMsg, lastMatched)
	return tmpMsg[:lastIndex] + tmpMsg[lastIndex+1:]
}
