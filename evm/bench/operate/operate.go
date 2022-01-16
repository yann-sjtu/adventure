package operate

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/evm/bench/utils"
	"github.com/okex/adventure/evm/constant"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	"github.com/spf13/cobra"
)

var (
	// used for flags
	contract string
	direct   bool
	id       int64
	opts     []int64
	times    int64

	// global variables
	eParam utils.TxParam
)

func operate(cmd *cobra.Command, args []string) {
	eParam = utils.NewTxParam(
		ethcmm.HexToAddress(contract),
		nil,
		2000000,
		evmtypes.DefaultGasPrice,
		generateTxData(),
	)

	utils.RunTxs(
		utils.DefaultBaseParamFromFlag(),
		func(_ ethcmm.Address) []utils.TxParam {
			return []utils.TxParam{eParam}
		},
	)
}

func generateTxData() []byte {
	bigOpts := make([]*big.Int, len(opts))
	for i := 0; i < len(opts); i++ {
		bigOpts[i] = big.NewInt(opts[i])
	}

	if direct {
		operateABI, err := abi.JSON(strings.NewReader(constant.OperateABI))
		if err != nil {
			panic(err)
		}
		txdata, err := operateABI.Pack("operate", bigOpts, big.NewInt(times))
		if err != nil {
			panic(err)
		}
		return txdata
	} else {
		routerABI, err := abi.JSON(strings.NewReader(constant.OperateRouterABI))
		if err != nil {
			panic(err)
		}
		txdata, err := routerABI.Pack("operate", big.NewInt(id), bigOpts, big.NewInt(times))
		if err != nil {
			panic(err)
		}
		return txdata
	}
}
