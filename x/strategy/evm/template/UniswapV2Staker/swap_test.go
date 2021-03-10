package UniswapV2Staker

import (
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/okexchain-go-sdk/utils"
)

const (
	routerAddr = "0x2CA0E1278B9D7A967967d3C707b81C72FC180CaF"

	oktUsdtPool = "0x5aFC0E1ddDd7a5151d83a3385C01e6159539a37C"
	oktUsdtLP = "0xe922FF7B02672bB59A64b90864FC5e511AD4d5fa"
	usdtAddr = "0xee666e967293094007d7c3718044e07565b1f8a9"
)

func TestBuilder(t *testing.T) {
	UniswapV2.Init()
	Init()
	clients := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	info, _, err := utils.CreateAccountWithMnemo("plunge silk glide glass curve cycle snack garbage obscure express decade dirt", fmt.Sprintf("acc%d", 1), "12345678")
	if err != nil {
		panic(err)
	}
	// 2. get cli
	cli := clients.GetClient()
	// 3. get acc number
	acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		panic(err)
	}
	accNum, seqNum := acc.GetAccountNumber(), acc.GetSequence()

	//payload, err := UniswapV2.PairBuilder.Build("approve", utils.EthAddress(routerAddr), sdk.NewDec(10000000000000000).Int)
	//if err != nil {
	//	panic(err)
	//}
	//res, err := cli.Evm().SendTx(info, common.PassWord, oktUsdtLP, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(res.TxHash)

	//payload := UniswapV2.BuildAddLiquidOKTPayload(
	//	usdtAddr, utils.GetEthAddressStrFromCosmosAddr(info.GetAddress()),
	//	sdk.MustNewDecFromStr("1").Int, sdk.MustNewDecFromStr("0.0001").Int, sdk.MustNewDecFromStr("0.0001").Int,
	//	time.Now().Add(time.Hour*24).Unix(),
	//)
	//res, err := cli.Evm().SendTx(info, common.PassWord, routerAddr, "0.1", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(res.TxHash)

	payload := UniswapV2.BuildRemoveLiquidOKTPayload(
		usdtAddr, "utils.ToHexAddress(info.GetAddress())",
		sdk.MustNewDecFromStr("900").Int, sdk.MustNewDecFromStr("8").Int, sdk.MustNewDecFromStr("0.109147").Int,
		time.Now().Add(time.Hour*24).Unix(),
	)
	res, err := cli.Evm().SendTx(info, common.PassWord, routerAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	log.Println(res.TxHash)
	if err != nil {
		panic(err)
	}

	//payload := BuildStakePayload(512775580224501)
	//payload := BuildWithdrawPayload(500000000)
	//payload := BuildGetRewardPayload()
	//payload := BuildExitPayload()
	//res, err := cli.Evm().SendTx(info, common.PassWord, oktUsdtPool, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(res.TxHash)
}

func Uint256(d sdk.Dec) *big.Int {
	return d.BigInt()
}