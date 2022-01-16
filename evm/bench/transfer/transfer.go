package transfer

import (
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/evm/bench/utils"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
	"github.com/spf13/cobra"
)

var (
	// used for flags
	fixed bool
)

func transfer(cmd *cobra.Command, args []string) {
	amount := sdk.MustNewDecFromStr("0.00001").Int
	fixedAddr := ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64)))

	utils.RunTxs(
		utils.DefaultBaseParamFromFlag(),
		func(addr ethcmm.Address) []utils.TxParam {
			to := fixedAddr
			if !fixed {
				to = ethcmm.BytesToAddress(crypto.Keccak256(addr.Bytes())[12:])
			}
			return []utils.TxParam{utils.NewTxParam(to, amount, 21000, evmtypes.DefaultGasPrice, nil)}
		},
	)
}
