package UniswapV2Staker

import (
	"fmt"
	"log"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/evm/template/UniswapV2"
	"github.com/okex/okexchain-go-sdk/utils"
)

const (
	routerAddr = "0x0653a68B22b18663F69a7103621F7F3EB59191F1"

	oktUsdtPool = "0x0Bd475f8b27EA57158291372667aD1e7eeD5C174"
	oktUsdtLP = "0x7068B191ff97e32D6Fbba3204408877A9007BBd1"
	usdtAddr = "0xffea71957a3101d14474a3c358ede310e17c2409"
)

func TestBuilder(t *testing.T) {
	UniswapV2.Init()
	Init()
	clients := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	info, _, err := utils.CreateAccountWithMnemo("expand present amused regret normal gallery coconut brass speed survey giraffe couple", fmt.Sprintf("acc%d", 1), "12345678")
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

	//payload, err := UniswapV2.PairBuilder.Build("approve", utils.EthAddress(oktUsdtPool), sdk.NewDec(10000000000000000).Int)
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
	//	6000000000000000,1,1,
	//	int(time.Now().Add(time.Hour*24).Unix()),
	//)
	//res, err := cli.Evm().SendTx(info, common.PassWord, routerAddr, "0.001", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(res.TxHash)

	//payload := BuildStakePayload(512775580224501)
	//payload := BuildWithdrawPayload(500000000)
	//payload := BuildGetRewardPayload()
	payload := BuildExitPayload()
	res, err := cli.Evm().SendTx(info, common.PassWord, oktUsdtPool, "", ethcommon.Bytes2Hex(payload), "", accNum, seqNum)
	if err != nil {
		panic(err)
	}
	log.Println(res.TxHash)
}

func Uint256(d sdk.Dec) *big.Int {
	return d.BigInt()
}