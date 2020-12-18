package ERC721

import (
	"github.com/okex/okexchain-go-sdk/utils"
)

var (
	ERC721Builder utils.PayloadBuilder
)

func Init() {
	var err error

	// 1. init builders
	ERC721Builder, err = utils.NewPayloadBuilder(ERC721Bin, ERC721ABI)
	if err != nil {
		panic(err)
	}
}

func BuildERC721ContractPayload(name, symbol string) []byte {
	payload, err := ERC721Builder.Build("", name, symbol)
	if err != nil {
		panic(err)
	}
	return payload
}
