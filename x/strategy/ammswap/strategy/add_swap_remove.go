package strategy

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/tools/account"
	"github.com/okex/adventure/x/strategy/token"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/okex/okexchain-go-sdk/utils"
	"github.com/spf13/cobra"
)

const passWd = common.PassWord

var (
	IssueMnemonic = ""
	IssueNum      = uint64(10)

	MnemonicPath = ""

	GoroutineNum = 10
)

func addSwapRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loop",
		Short: "start tests among add->swap->remove in an infinite loop",
		Args:  cobra.NoArgs,
		RunE:  addSwapRemoveScripts,
	}

	flags := cmd.Flags()
	flags.StringVarP(&IssueMnemonic, "issue_mnemonic", "i", "", "the mnemonic used for issuing tokens")
	flags.Uint64VarP(&IssueNum, "issue_num", "n", 10, "set num of how many tokens to be issued")
	flags.StringVarP(&MnemonicPath, "mnemonic_path", "p", "", "the mnemonic file path for testing swap tx")
	flags.IntVarP(&GoroutineNum, "goroutine_num", "g", 10, "set goroutine number")
	//flags.Uint64VarP(&num, "num", "n", 1000, "set num of issusing token")

	return cmd
}

var (
	// map[baseToken] -> [quoteToken1, quoteToken2...]
	tokenPairsMap1 = make(map[string][]string)
	// map[quoteToken] -> [baseToken1, baseToken2...]
	tokenPairsMap2 = make(map[string][]string)

	// map[addr] -> [token1, token2, token3...]
	addrTokenMap = make(map[string][]string)
)

type RoundCounter struct {
	round int
	lock  *sync.RWMutex
}

func NewRoundConter() *RoundCounter {
	return &RoundCounter{
		round: 0,
		lock:  new(sync.RWMutex),
	}
}

func (rc *RoundCounter) AddRound() int {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	rc.round++
	return rc.round
}

func addSwapRemoveScripts(cmd *cobra.Command, args []string) error {
	clis := common.NewClientManager(common.Cfg.Hosts, common.AUTO)
	infos := common.GetAccountManagerFromFile(MnemonicPath)

	// 0. Issue tokens & create token pairs in ammswap, when IssueMnemonic is not empty
	if IssueMnemonic != "" {
		//if IssueNum% 2 == 1 {
		//	return fmt.Errorf("[phase 0] IssueNum %d is invaild, it must be divisible by 2", IssueNum)
		//}

		// 0.1 create account info
		info, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(IssueMnemonic), "issueAcc", common.PassWord)
		if err != nil {
			return fmt.Errorf("[phase 0] issue mnemonic invaild: %s", err)
		}
		cli := clis.GetRandomClient()
		// 0.2 issue tokens
		err = IssueTokens(cli, info)
		if err != nil {
			return fmt.Errorf("[phase 0] isuue tokens failed: %s", err)
		}
		// 0.3 create swap pairs
		time.Sleep(time.Second * 5)
		err = CreateTokenPairs(cli, info)
		if err != nil {
			return fmt.Errorf("[phase 0] create swap pairs in ammswap failed: %s", err)
		}
		// 0.4 send swap coins to addrs
		err = SendCoins(cli, IssueMnemonic, infos.GetInfos())
		if err != nil {
			return fmt.Errorf("[phase 0] send swap tokens to addrs failed: %s", err)
		}
	}

	// 2. init the pair names in swap pool. (map[name1] -> name2)
	err := InitSwapPairsMap(clis.GetRandomClient())
	if err != nil {
		return err
	}

	// 3. create a number of goroutine
	successfulRound, failedRound := NewRoundConter(), NewRoundConter()
	for i := 0; i < GoroutineNum; i++ {
		go func() {
			for {

				// 3.0 get account info
				info, cli := infos.GetInfo(), clis.GetClient()
				addr := info.GetAddress().String()

				// 3.1 get tokens in map
				tokens, err := GetTokensInAddrMap(cli, addr)
				if err != nil {
					fmt.Printf("[%d] %s failed to get tokens: %s", failedRound.AddRound(), addr, err)
					continue
				}

				name1, name2 := PickTwoTokensRandomly(tokens)
				// 3.2 pick one random token
				if name1 == "" || name2 == "" {
					fmt.Printf("[%d] %s : failed to get matached with one empty token [%s <-> %s]\n", failedRound.AddRound(), addr, name1, name2)
					continue
				}

				// 3.3.1 add liquidility tx
				accNum, seqNum := getAccountInfo(cli, addr)
				_, err = cli.AmmSwap().AddLiquidity(info, passWd, "0.1", "1"+name1, "0.05"+name2, "1m", "", accNum, seqNum)
				if err != nil {
					fmt.Printf("[%d] %s failed to AddLiquidity: %s \n", failedRound.AddRound(), addr, err)
					continue
				}
				// 3.3.2 swap token tx
				_, err = cli.AmmSwap().TokenSwap(info, passWd, "0.01"+name1, "0.0000001"+name2, addr, "5m", "", accNum, seqNum+uint64(1))
				if err != nil {
					fmt.Printf("[%d] %s failed to TokenSwap: %s \n", failedRound.AddRound(), addr, err)
					continue
				}
				_, err = cli.AmmSwap().TokenSwap(info, passWd, "0.003"+name2, "0.0000001"+name1, addr, "5m", "", accNum, seqNum+uint64(2))
				if err != nil {
					fmt.Printf("[%d] %s failed to TokenSwap: %s \n", failedRound.AddRound(), addr, err)
					continue
				}
				// 3.3.3 remove liquidility tx
				_, err = cli.AmmSwap().RemoveLiquidity(info, passWd, "0.02", "0"+name1, "0"+name2, "2m", "", accNum, seqNum+uint64(3))
				if err != nil {
					fmt.Printf("[%d] %s failed to RemoveLiquidity: %s \n", failedRound.AddRound(), addr, err)
					continue
				}
				fmt.Printf("[%d] %s: finish a round of test about %s\n", successfulRound.AddRound(), addr, name1+"_"+name2)
			}
		}()
	}

	select {}
}

func IssueTokens(cli *gosdk.Client, info keys.Info) error {
	accNum, seqNum := getAccountInfo(cli, info.GetAddress().String())
	for i := uint64(0); i < IssueNum; i++ {
		name := token.GetRandomString(3)
		_, err := cli.Token().Issue(info, common.PassWord,
			name, name, "9990000000.00000000", "Used for test "+name+" "+strconv.Itoa(int(seqNum+i)),
			"test token for ammswap", true, accNum, seqNum+i)
		if err != nil {
			return err
		}
		fmt.Printf("[phase 0] issue token successfully: %s \n", name)
	}
	return nil
}

func getAccountInfo(cli *gosdk.Client, addr string) (uint64, uint64) {
	// query account info
	accInfo, err := cli.Auth().QueryAccount(addr)
	if err != nil {
		fmt.Println(addr, "querying account info is failed.")
		return 0, 0
	}
	return accInfo.GetAccountNumber(), accInfo.GetSequence()
}

func CreateTokenPairs(cli *gosdk.Client, info keys.Info) error {
	acc, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		return err
	}
	coins := acc.GetCoins()
	if coins.Empty() {
		return fmt.Errorf("there are no coins in %s", info.GetAddress().String())
	}

	accNum, seqNum := getAccountInfo(cli, info.GetAddress().String())
	seqNumOffset := 0
	for i := 0; i < len(coins); i++ { //
		if !isTokenWithSuffix(coins[i].Denom) {
			continue
		}
		for j := i + 1; j < len(coins); j++ {
			if !isTokenWithSuffix(coins[j].Denom) {
				continue
			}

			_, err := cli.AmmSwap().CreateExchange(info, common.PassWord,
				coins[i].Denom, coins[j].Denom,
				"", accNum, seqNum+uint64(seqNumOffset))
			if err != nil {
				return err
			}
			fmt.Printf("[phase 0] create swap pairs in ammswap successfully: %s <-> %s \n", coins[i].Denom, coins[j].Denom)

			seqNumOffset++
			i = j
			break
		}
	}

	return nil
}

func isTokenWithSuffix(name string) bool {
	if name == common.NativeToken {
		return false // okt false
	}
	if strings.Contains(name, "-") && !strings.Contains(name, "_") {
		return true // xxb-abc true
	}
	return false // ammswap-xxb1_xxb2 false
}

func SendCoins(cli *gosdk.Client, mnemonic string, toAddrInfos []keys.Info) error {
	// query coins
	info, _, err := utils.CreateAccountWithMnemo(strings.TrimSpace(mnemonic), "issueAcc", common.PassWord)
	if err != nil {
		return err
	}

	var coinStr string
	if acc, err := cli.Auth().QueryAccount(info.GetAddress().String()); err != nil {
		return err
	} else {
		for _, coin := range acc.GetCoins() {
			if isTokenWithSuffix(coin.Denom) {
				coinStr += "100000" + coin.Denom + ","
			}
		}
	}

	addrs := make([]string, len(toAddrInfos), len(toAddrInfos))
	for i := 0; i < len(toAddrInfos); i++ {
		addrs[i] = toAddrInfos[i].GetAddress().String()
	}

	err = account.SendCoins(cli, addrs, strings.TrimRight(coinStr, ","), mnemonic)
	if err != nil {
		return err
	}

	fmt.Printf("[phase 0] send coins %s to %d address successfully \n", coinStr, len(addrs))
	return nil
}

func InitSwapPairsMap(cli *gosdk.Client) error {
	// init the pair names in swap pool. (map[name1] -> name2)
	swapPairs, err := cli.AmmSwap().QuerySwapTokenPairs()
	if err != nil {
		return err
	}
	for _, swapPair := range swapPairs {
		name1, name2 := swapPair.BasePooledCoin.Denom, swapPair.QuotePooledCoin.Denom
		tokenPairsMap1[name1] = append(tokenPairsMap1[name1], name2)
		tokenPairsMap2[name2] = append(tokenPairsMap2[name2], name1)
	}
	return nil
}

func GetTokensInAddrMap(cli *gosdk.Client, addr string) ([]string, error) {
	if _, ok := addrTokenMap[addr]; !ok {
		acc, err := cli.Auth().QueryAccount(addr)
		if err != nil {
			return nil, err
		}

		var tokens []string
		for _, token := range acc.GetCoins() {
			name := token.Denom
			if isTokenWithSuffix(name) {
				tokens = append(tokens, name)
			}
		}
		addrTokenMap[addr] = tokens
	}

	tokens := addrTokenMap[addr]
	if len(tokens) <= 1 {
		return nil, fmt.Errorf("don't have more than two types of token %v\n", tokens)
	}
	return tokens, nil
}

func PickTwoTokensRandomly(tokens []string) (string, string) {
	rand.Seed(time.Now().UnixNano() + rand.Int63n(10000))

	var name1, name2 string
	name1 = tokens[rand.Intn(len(tokens))]
	if name2List, ok := tokenPairsMap1[name1]; ok {
		for _, token := range tokens {
			if index := stringsContains(name2List, token); index != -1 {
				name2 = name2List[index]
				break
			}
		}
	}
	if name2 == "" {
		if name1List, ok := tokenPairsMap2[name1]; ok {
			for _, token := range tokens {
				if index := stringsContains(name1List, token); index != -1 {
					name2 = name1List[index]
					break
				}
			}
		}
	}

	return compareTokenNames(name1, name2)
}

func stringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

func compareTokenNames(name1, name2 string) (string, string) {
	if name1 > name2 {
		return name2, name1
	} else {
		return name1, name2
	}
}
