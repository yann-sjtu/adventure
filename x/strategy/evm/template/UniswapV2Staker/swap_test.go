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
	routerAddr = "0x0653a68B22b18663F69a7103621F7F3EB59191F1"

	usdtAddr = "0xffea71957a3101d14474a3c358ede310e17c2409"
)

func TestBuilder(t *testing.T) {
	UniswapV2.Init()
	clients := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	// 1. get one of the eth private key
	info, _, err := utils.CreateAccountWithMnemo("nose lend select ball vocal box speed custom energy caution order hole", fmt.Sprintf("acc%d", 1), "12345678")
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

	toEthAddress := utils.EthAddress(utils.GetEthAddressStrFromCosmosAddr(info.GetAddress()))
	payload, err := UniswapV2.RouterBuilder.Build("addLiquidityETH",
		utils.EthAddress(usdtAddr),
		Uint256(sdk.MustNewDecFromStr("4.062232370071723288")), Uint256(sdk.MustNewDecFromStr("3.041921208221365886")),
		Uint256(sdk.MustNewDecFromStr("90.5000000000000000")),
		utils.EthAddress(toEthAddress.String()), big.NewInt(time.Now().Add(time.Hour).Unix()),
	)
	res, err := cli.Evm().SendTx(info, common.PassWord, routerAddr, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	if err != nil {
		panic(err)
	}
	log.Println(res)
}

func Uint256(d sdk.Dec) *big.Int {
	return d.BigInt()
}