package scenario

import (
	"github.com/okex/adventure/evm/bench/utils"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/rand"

	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/evm/constant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

)

var (
	// used for flags
	fixed bool
)

/**
该文件模拟用户场景，组合多个用户操作在一起，构成一个场景，放在一个协程中执行
 */

func scenario(cmd *cobra.Command, args []string) {
	amount := sdk.MustNewDecFromStr("0.00001").Int
	fixedAddr := ethcmm.BytesToAddress(crypto.Keccak256(rand.Bytes(64)))
	var toAddrs []ethcmm.Address
	if !fixed {
		toAddrs = LoadAddress()
	}

	utils.GetBalTxBal(
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


/**
遍历账户，获取账户余额， 转账，然后再次获取用户余额，验证下余额变化，余额减少
 */


func LoadAddress() []ethcmm.Address {
	privateKeyFile := viper.GetString(constant.FlagPrivateKeyFile);
	if privateKeyFile == "" {
		panic("Private key file is necessay")
	}
	privateKeys := common.ReadDataFromFile(privateKeyFile)

	n := len(privateKeys)
	if n==0 {
		panic("Private key is not found")
	}
	adds := make([]ethcmm.Address,0, n)
	for i :=0; i<n; i++ {
		pk,err := crypto.HexToECDSA(privateKeys[i])
		utils.PanicErr(err)
		adds = append(adds, common.GetEthAddressFromPK(pk))
	}
	return adds
}


