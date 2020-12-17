package types

import (
	"fmt"
	gokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/x/farm/client"
	"github.com/okex/adventure/x/farm/constants"
	"github.com/okex/adventure/x/farm/utils"
	"github.com/okex/okexchain-go-sdk"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Pooler struct {
	id      int
	accAddr string
	key     gokeys.Info
	// issued
	issuedTokenSymbol string
	lptSymbol         string
	totalSupply       sdk.Dec
}

func NewPooler(info gokeys.Info, id int) *Pooler {
	return &Pooler{
		id:      id,
		accAddr: info.GetAddress().String(),
		key:     info,
	}
}

func (p *Pooler) GetKey() gokeys.Info {
	return p.key
}

func (p *Pooler) IssueToken(wg *sync.WaitGroup) {
	defer wg.Done()

	cli := client.CliManager.GetClient()

	var (
		accInfo gosdk.Account
		e       error
	)

	for {
		accInfo, e = cli.Auth().QueryAccount(p.accAddr)
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	orgSymbol := utils.GetRandomTokenSymbol()
	res, e := cli.Token().Issue(
		p.key,
		common.PassWord,
		orgSymbol,
		fmt.Sprintf("air token %s", orgSymbol),
		"1000000000",
		fmt.Sprintf("air token %s for test", orgSymbol),
		"",
		true,
		accInfo.GetAccountNumber(),
		accInfo.GetSequence(),
	)
	if e != nil {
		log.Printf("failed. pooler %s issue token: %s\n",
			p.key.GetAddress().String(), e.Error())
	}
	fmt.Println(res)
}

func (p *Pooler) AddLiquidity(wg *sync.WaitGroup) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	if err := p.updateIssueTokenInfo(cli, false); err != nil {
		panic(err)
	}

	// add liquidity
	var (
		accInfo gosdk.Account
		e       error
	)
	for {
		accInfo, e = cli.Auth().QueryAccount(p.accAddr)
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	baseAmount, quoteAmount := utils.SortTokenSymbol(
		p.issuedTokenSymbol,
		common.DefaultStableCoin,
		constants.IssuedAmount,
		constants.StableCoinAmount)
	_, err := cli.AmmSwap().AddLiquidity(p.key, common.PassWord, constants.MinLiquidity,
		baseAmount,
		quoteAmount,
		constants.Duration,
		"",
		accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		log.Println(err)
	}

	// get large amount of lpt
	baseAmount, quoteAmount = utils.SortTokenSymbol(
		p.issuedTokenSymbol,
		common.DefaultStableCoin,
		constants.IssuedAmount*constants.Times,
		constants.StableCoinAmount*constants.Times,
	)
	res, err := cli.AmmSwap().AddLiquidity(p.key, common.PassWord, constants.MinLiquidity,
		baseAmount,
		quoteAmount,
		constants.Duration,
		"",
		accInfo.GetAccountNumber(), accInfo.GetSequence()+1)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
	return
}

func (p *Pooler) CreateSwapPairWithUSDK(wg *sync.WaitGroup) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	if err := p.updateIssueTokenInfo(cli, false); err != nil {
		panic(err)
	}

	// create swap pair
	var (
		accInfo gosdk.Account
		e       error
	)
	for {
		accInfo, e = cli.Auth().QueryAccount(p.accAddr)
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	res, err := cli.AmmSwap().CreateExchange(p.key, common.PassWord, common.DefaultStableCoin, p.issuedTokenSymbol, "",
		accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
	return
}

func (p *Pooler) CreateFarmPoolWithRandomTokenAndProvide(wg *sync.WaitGroup, latestHeight int64, isInWhiteList bool) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	if p.updateIssueTokenInfo(cli, true) != nil {
		return
	}

	// create farm pool
	var (
		accInfo gosdk.Account
		err     error
	)
	for {
		accInfo, err = cli.Auth().QueryAccount(p.accAddr)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// choose the token to lock randomly
	lockedTokens := []string{p.issuedTokenSymbol, p.lptSymbol}
	luckyNum := rand.Intn(2)
	poolName := fmt.Sprintf("%s pool|%s|%s|%s", utils.PrefixAdventurePoolName, lockedTokens[luckyNum], p.issuedTokenSymbol, time.Now().Format("2006-01-02 15:04:05"))
	minLockAmount := fmt.Sprintf("%s%s", constants.MinLockAmount, lockedTokens[luckyNum])
	if _, err = cli.Farm().CreatePool(p.key, common.PassWord, poolName, minLockAmount, p.issuedTokenSymbol, "",
		accInfo.GetAccountNumber(), accInfo.GetSequence()); err != nil {
		log.Printf("Tx error. [%s] creates pool [%s] with min-lock [%s] failed: %s\n", p.accAddr, poolName, minLockAmount, err)
	} else {
		log.Printf("[%s] creates pool [%s] with min-lock [%s] successfully\n", p.accAddr, poolName, minLockAmount)
	}

	// keep the tx successful
	time.Sleep(constants.BlockInterval * 2)

	// provide the farm pool
	// get start height 2 yield randomly
	startHeight2Yield := latestHeight + int64(constants.StartYieldHeightInterval+rand.Intn(constants.RandomRange))
	if _, err = cli.Farm().Provide(p.key, common.PassWord,
		poolName,
		fmt.Sprintf("%f%s", constants.IssuedTokenAmountSupply, p.issuedTokenSymbol),
		fmt.Sprintf("%f", constants.YieldAmountPerBlock),
		startHeight2Yield,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence()+1); err != nil {
		log.Printf("Tx error. [%s] provides pool [%s] failed: %s\n", p.accAddr, poolName, err)
	} else {
		log.Printf("[%s] provides pool [%s] successfully\n", p.accAddr, poolName)
	}

	if !isInWhiteList {
		return
	}

	// keep the tx successful
	time.Sleep(constants.BlockInterval * 2)

	// set pool in whitelist
	// TODO
	//if _, err = cli.Farm().SetWhite(p.key, common.PassWord, poolName, "", accInfo.GetAccountNumber(), accInfo.GetSequence()+2); err != nil {
	//	log.Printf("Tx error. [%s] sets pool [%s] into white list failed: %s\n", p.accAddr, poolName, err)
	//} else {
	//	log.Printf("[%s] sets pool [%s] into white list successfully\n", p.accAddr, poolName)
	//}

	return
}

func (p *Pooler) ProvideFarmPool(wg *sync.WaitGroup, pool gosdk.FarmPool, latestHeight int) {
	defer wg.Done()

	cli := client.CliManager.GetClient()

	var (
		accInfo gosdk.Account
		e       error
	)
	for {
		accInfo, e = cli.Auth().QueryAccount(p.accAddr)
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// get start height 2 yield randomly
	startHeight2Yield := int64(latestHeight + 1000 + rand.Intn(500))
	res, err := cli.Farm().Provide(p.key, common.PassWord,
		pool.Name,
		fmt.Sprintf("%f%s", constants.IssuedTokenAmountSupply, pool.YieldedTokenInfos[0].RemainingAmount.Denom),
		fmt.Sprintf("%f", constants.YieldAmountPerBlock),
		startHeight2Yield,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence(),
	)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
	return
}

func (p *Pooler) UpdateIssueTokenInfo(cli *gosdk.Client, containsLPT bool) error {
	return p.updateIssueTokenInfo(cli, containsLPT)
}

func (p *Pooler) updateIssueTokenInfo(cli *gosdk.Client, containsLPT bool) error {
	// get the token that the pooler issued
	var (
		tokenInfo []gosdk.TokenResp
		e         error
	)
	for {
		tokenInfo, e = cli.Token().QueryTokenInfo(p.key.GetAddress().String(), "")
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if len(tokenInfo) != 1 {
		return fmt.Errorf("failed. pooler %s has issued more that 1 token",
			p.accAddr)
	}

	p.issuedTokenSymbol = tokenInfo[0].Symbol
	p.totalSupply = tokenInfo[0].TotalSupply

	// update the lpt field in pooler
	if containsLPT {
		p.lptSymbol = utils.BuildLPTName(p.issuedTokenSymbol, common.DefaultStableCoin)
	}

	return nil
}

func (p *Pooler) ProvidePool(wg *sync.WaitGroup, pool gosdk.FarmPool, currentHeight int64) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	accInfo, err := cli.Auth().QueryAccount(p.accAddr)
	if err != nil {
		return
	}

	startHeight2Yield := currentHeight + int64(constants.StartYieldHeightInterval+rand.Intn(constants.RandomRange))
	if _, err = cli.Farm().Provide(p.key, common.PassWord,
		pool.Name,
		fmt.Sprintf("%f%s", constants.IssuedTokenAmountSupply, pool.YieldedTokenInfos[0].RemainingAmount.Denom),
		fmt.Sprintf("%f", constants.YieldAmountPerBlock),
		startHeight2Yield,
		"", accInfo.GetAccountNumber(), accInfo.GetSequence()); err != nil {
		log.Printf("Tx error. [%s] provides pool [%s] failed: %s\n", p.accAddr, pool.Name, err)
	} else {
		log.Printf("[%s] provides pool [%s] successfully\n", p.accAddr, pool.Name)
	}
}

func (p *Pooler) DestroyPool(wg *sync.WaitGroup, poolName string) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	accInfo, err := cli.Auth().QueryAccount(p.accAddr)
	if err != nil {
		return
	}

	if _, err = cli.Farm().DestroyPool(p.key, common.PassWord, poolName, "", accInfo.GetAccountNumber(),
		accInfo.GetSequence()); err != nil {
		log.Printf("Tx error. [%s] destroys pool [%s] failed: %s\n", p.accAddr, poolName, err)
		return
	}

	log.Printf("[%s] destroys pool [%s] successfully\n", p.accAddr, poolName)
}

func (p *Pooler) MultiSendToAddrs(wg *sync.WaitGroup, accAddrs []sdk.AccAddress) {
	defer wg.Done()

	cli := client.CliManager.GetClient()
	if err := p.updateIssueTokenInfo(cli, false); err != nil {
		panic(err)
	}

	// create swap pair
	var (
		accInfo gosdk.Account
		e       error
	)
	for {
		accInfo, e = cli.Auth().QueryAccount(p.accAddr)
		if e == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	res, err := cli.AmmSwap().CreateExchange(p.key, common.PassWord, common.DefaultStableCoin, p.issuedTokenSymbol, "",
		accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
	return
}

func (p *Pooler) GetIssuedTokenSymbol() string {
	return p.issuedTokenSymbol
}

func (p *Pooler) GetLPTSymbol() string {
	return p.lptSymbol
}
