package DYF

import "github.com/okex/exchain-go-sdk/utils"

var (
	DYFBuilder utils.PayloadBuilder
)

func Init() {
	var err error

	// 1. init builders
	DYFBuilder, err = utils.NewPayloadBuilder(DYFBin, DYFAbi)
	if err != nil {
		panic(err)
	}
}

func BuildExcutePayload() []byte {
	payload, err := DYFBuilder.Build("go")
	if err != nil {
		panic(err)
	}
	return payload
}