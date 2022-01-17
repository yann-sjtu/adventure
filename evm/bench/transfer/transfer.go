package transfer

import (
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/evm/bench/utils"
	"github.com/okex/adventure/evm/constant"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// used for flags
	fixed bool
)

func transfer(cmd *cobra.Command, args []string) {
	amount := sdk.MustNewDecFromStr("0.00001").Int
	fixedAddr := ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64)))
	var toAddrs []ethcmm.Address
	if !fixed {
		toAddrs = generateAddress()
	}

	utils.RunTxs(
		utils.DefaultBaseParamFromFlag(),
		func(addr ethcmm.Address) []utils.TxParam {
			to := fixedAddr
			if !fixed {
				to = toAddrs[rand.Intn(len(toAddrs))]
			}
			return []utils.TxParam{utils.NewTxParam(to, amount, 21000, evmtypes.DefaultGasPrice, nil)}
		},
	)
}

func generateAddress() []ethcmm.Address {
	privateKeys := constant.PrivateKeys
	if privateKeyFile := viper.GetString(constant.FlagPrivateKeyFile); privateKeyFile != "" {
		privateKeys = common.ReadDataFromFile(privateKeyFile)
	}

	leng := len(privateKeys)
	addrs := make([]ethcmm.Address, leng, leng)
	for i := 0; i < leng; i++ {
		pk, err := crypto.HexToECDSA(privateKeys[i])
		if err != nil {
			panic(err)
		}
		addrs[i] = common.GetEthAddressFromPK(pk)
	}
	return addrs
}
