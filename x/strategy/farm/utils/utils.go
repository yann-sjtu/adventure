package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	gokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/strategy/farm/client"
	"github.com/okex/adventure/x/strategy/farm/constants"
	"github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/utils"
	tokentypes "github.com/okex/okexchain/x/token/types"
)

const (
	PrefixAdventurePoolName = "adventure"
	richerMnemonics         = "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
	richerName              = "richer"
	RichAddr                = "okexchain1pt7xrmxul7sx54ml44lvv403r06clrdkgmvr9g"
)

var (
	richerInfo                   gokeys.Info
	validPoolBlockHeightInterval = int64(math.Ceil(constants.IssuedTokenAmountSupply/constants.YieldAmountPerBlock)) + 1
)

func init() {
	richerInfo, _, _ = utils.CreateAccountWithMnemo(richerMnemonics, richerName, common.PassWord)
}

func GetRicherKeyInfo() gokeys.Info {
	return richerInfo
}

func GetTestAddrsFromFile(path string) (addrStrs []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		addr, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		addrStrs = append(addrStrs, strings.TrimSpace(addr))
	}

	return
}

// if u want to control the number of the accounts from the file, please input the params 'num'
func GetTestAccountsFromFile(path string, num ...int) (infos []keys.Info, err error) {
	if len(num) > 1 {
		return infos, errors.New("failed. num input is invalid")
	}

	var isNumLimited bool
	if len(num) == 1 {
		isNumLimited = true
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	var index int
	var accounts []keys.Info

	for {
		if isNumLimited && index >= num[0] {
			break
		}

		mnemonic, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		acc, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(mnemonic),
			fmt.Sprintf("pooler%d", index), common.PassWord)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
		fmt.Println(accounts[index].GetAddress().String(), index)
		index++
	}

	return accounts, nil
}

func GetRandomTokenSymbol() string {
	rand.Seed(time.Now().UnixNano())
	var buffer []byte
	for i := 0; i < 3; i++ {
		buffer = append(buffer, byte(rand.Intn(26)+'a'))
	}

	return string(buffer)
}

func SortTokenSymbol(tSymbol1, tSymbol2 string, tAmount1, tAmount2 int) (string, string) {
	amount1 := fmt.Sprintf("%d%s", tAmount1, tSymbol1)
	amount2 := fmt.Sprintf("%d%s", tAmount2, tSymbol2)
	if tSymbol1 > tSymbol2 {
		return amount2, amount1
	}

	return amount1, amount2
}

func BuildLPTName(tokenName0, tokenName1 string) string {
	if tokenName0 > tokenName1 {
		return fmt.Sprintf("ammswap_%s_%s", tokenName1, tokenName0)
	}

	return fmt.Sprintf("ammswap_%s_%s", tokenName0, tokenName1)
}

func IssueStableCoin() error {
	cli := client.CliManager.GetClient()

	var (
		accInfo gosdk.Account
		err     error
	)
	for {
		accInfo, err = cli.Auth().QueryAccount(RichAddr)
		if err == nil {
			break
		}
	}

	if _, err = cli.Token().Issue(richerInfo, common.PassWord, "usdk", "usdk",
		"90000000000", "stable coin on okexchain", "", true, accInfo.GetAccountNumber(),
		accInfo.GetSequence()); err != nil {
		log.Printf("Tx error. [%s] issues %s failed: %s\n", RichAddr, common.DefaultStableCoin, err)
		return err
	}

	log.Printf("[%s] issues %s successfully\n", common.DefaultStableCoin, RichAddr)
	return nil
}

func BuildTransferUnits(accAddrs []sdk.AccAddress, coinsStr string) (transferUnits []tokentypes.TransferUnit, err error) {
	coins, err := sdk.ParseDecCoins(coinsStr)
	if err != nil {
		return
	}

	for _, accAddr := range accAddrs {
		transferUnits = append(transferUnits, tokentypes.TransferUnit{To: accAddr, Coins: coins})
	}

	return
}

func ParseAccAddrsFromAddrsStr(addrsStrs []string) (accAddrs []sdk.AccAddress, err error) {
	for _, addrStr := range addrsStrs {
		accAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, err
		}

		accAddrs = append(accAddrs, accAddr)
	}

	return
}

func MultiSendByGroup(wg *sync.WaitGroup, senderKeyInfo gokeys.Info, coinsStr string, accAddrs []sdk.AccAddress, addrNumOneTime int) (err error) {
	if wg != nil {
		defer wg.Done()
	}

	cli := client.CliManager.GetClient()

	transferUnits, err := BuildTransferUnits(accAddrs, coinsStr)
	if err != nil {
		return
	}

	var accInfo gosdk.Account
	senderAddrStr := senderKeyInfo.GetAddress().String()

	for {
		accInfo, err = cli.Auth().QueryAccount(senderAddrStr)
		if err == nil {
			break
		}
	}
	var times int
	if len(transferUnits)%addrNumOneTime != 0 {
		times = len(transferUnits)/addrNumOneTime + 1
	} else {
		times = len(transferUnits) / addrNumOneTime
	}

	for i := 0; i < times; i++ {
		var index2 int
		index1 := i * addrNumOneTime
		if i != times-1 {
			index2 = (i + 1) * addrNumOneTime
		} else {
			index2 = len(transferUnits)
		}

		if _, err = cli.Token().MultiSend(senderKeyInfo, common.PassWord, transferUnits[index1:index2], "",
			accInfo.GetAccountNumber(), accInfo.GetSequence()+uint64(i)); err != nil {
			log.Printf("Tx error. [%s] allocates [%s] to each of accounts with index %d to %d failed: %s\n",
				senderAddrStr, coinsStr, index1, index2, err)
			continue
		}

		log.Printf("[%s] allocates [%s] to each of accounts with index %d to %d successfully\n",
			senderAddrStr, coinsStr, index1, index2)
		time.Sleep(constants.SleepTimeBtwGroupBroadcast)
	}

	return nil
}

func GetPoolsRandomly() (pickedPools []gosdk.FarmPool, err error) {
	pools, err := client.CliManager.GetClient().Farm().QueryPools()
	if err != nil {
		return
	}

	maxRange := len(pools)
	if maxRange == 0 {
		return pickedPools, errors.New("no pools are queried from OKEXChain")
	}
	rand.Seed(time.Now().UnixNano())
	// get two random number as lucky num and selected pools' numbers
	selectedPoolsNum := rand.Intn(maxRange)
	luckyNum := rand.Intn(maxRange)
	pickedPools = make([]gosdk.FarmPool, selectedPoolsNum)

	for i := 0; i < selectedPoolsNum; i++ {
		pickedPools[i] = pools[(luckyNum+i)%maxRange]
	}

	fmt.Printf(`
============================================================
|               %d pools are picked randomly                |
============================================================

`, selectedPoolsNum)
	for _, pool := range pickedPools {
		fmt.Println(pool.Name)
	}

	return
}

func GetExpiredPoolsOnCurrentHeight() (expiredPools []gosdk.FarmPool, currentHeight int64, err error) {
	// 1.get current height
	currentHeight, err = queryLatestBlockHeight()
	if err != nil {
		return
	}

	// 2.get the info of pools that are controlled by ADVENTURE
	pools, err := QueryAllAdventurePools()
	if err != nil {
		return
	}

	// 3.get the pools that are expired
	for _, pool := range pools {
		if isPoolExpired(&pool, currentHeight) {
			expiredPools = append(expiredPools, pool)
		}
	}

	return
}

func queryLatestBlockHeight() (height int64, err error) {
	block, err := client.CliManager.GetClient().Tendermint().QueryBlock(0)
	if err != nil {
		return
	}

	return block.Height, err
}

func QueryAllAdventurePools() (adventurePools []gosdk.FarmPool, err error) {
	pools, err := client.CliManager.GetClient().Farm().QueryPools()
	if err != nil {
		return
	}

	// filter the pools that are controlled by ADVENTURE
	for _, pool := range pools {
		if strings.HasPrefix(pool.Name, PrefixAdventurePoolName) {
			adventurePools = append(adventurePools, pool)
		}
	}

	return
}

func isPoolExpired(pool *gosdk.FarmPool, currentHeight int64) bool {
	return currentHeight-pool.YieldedTokenInfos[0].StartBlockHeightToYield >= validPoolBlockHeightInterval
}

func GetRandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		return false
	}

	return true
}
