package common

import (
	"github.com/cosmos/cosmos-sdk/x/auth/types"
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

func NewObject(msgStr string, time string, txType int) Object {
	return Object{
		CoinType:  coinType,
		Hex:       msgStr,
		RelatedId: time,
		TxType:    txType,
	}
}
