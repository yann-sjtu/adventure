package strategy

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
	"github.com/spf13/cobra"
)

const passWd = common.PassWord

var (
	mnemonicPath = ""

	goroutineNum = 10
)

func addSwapRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loop",
		Short: "start tests among add->swap->remove in an infinite loop",
		Args:  cobra.NoArgs,
		RunE:  addSwapRemoveScripts,
	}

	flags := cmd.Flags()
	flags.StringVarP(&mnemonicPath, "mnemonic_path", "p", "", "set account mnemonic")
	flags.IntVarP(&goroutineNum, "goroutine_num", "g", 10, "set goroutine number")
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

func addSwapRemoveScripts(cmd *cobra.Command, args []string) error {
	clis := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	// init the pair names in swap pool. (map[name1] -> name2)
	err := initSwapPairsMap(clis.GetRandomClient())
	if err != nil {
		return err
	}

	// create a number of goroutine
	infos := common.GetAccountManagerFromFile(mnemonicPath)
	for i := 0; i < goroutineNum; i++ {
		go func(id int) {
			round := 0
			for {
				round++

				info, cli := infos.GetInfo(), clis.GetClient()
				addr := info.GetAddress().String()

				// get tokens in map
				tokens, err := getTokensInAddrMap(cli, addr)
				if err != nil {
					fmt.Printf("[%d] round(%d) %s failed: %s", id, round, addr, err)
					continue
				}

				// pick one random token
				if name2 == "" {
					fmt.Printf("[%d] round(%d) %s: there is no specific token %s adapted in token-list %v \n", id, round, addr, name1, tokens)
					continue
				}

				// add liquidility tx
				accNum, seqNum := getAccountInfo(cli, addr)
				_, err = cli.AmmSwap().AddLiquidity(info, passWd, "0.1", "1"+name1, "0.05"+name2, "1m", "", accNum, seqNum)
				if err != nil {
					fmt.Println(addr, tokens, err)
					continue
				}
				// swap token tx
				_, err = cli.AmmSwap().TokenSwap(info, passWd, "0.01"+name1, "0.0000001"+name2, addr, "5m", "", accNum, seqNum+uint64(1))
				if err != nil {
					fmt.Println(addr, tokens, err)
					continue
				}
				_, err = cli.AmmSwap().TokenSwap(info, passWd, "0.003"+name2, "0.0000001"+name1, addr, "5m", "", accNum, seqNum+uint64(2))
				if err != nil {
					fmt.Println(addr, tokens, err)
					continue
				}
				// remove liquidility tx
				_, err = cli.AmmSwap().RemoveLiquidity(info, passWd, "0.02", "0"+name1, "0"+name2, "2m", "", accNum, seqNum+uint64(3))
				if err != nil {
					fmt.Println(addr, tokens, err)
					continue
				}
				fmt.Printf("[%d] round(%d) %s: finish a round of test about %s\n", id, round, addr, name1+"_"+name2)
				time.Sleep(time.Second * 1)
			}
		}(i)
	}

	select {}
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

func initSwapPairsMap(cli *gosdk.Client) error {
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

func getTokensInAddrMap(cli *gosdk.Client, addr string) ([]string, error) {
	if _, ok := addrTokenMap[addr]; !ok {
		acc, err := cli.Auth().QueryAccount(addr)
		if err != nil {
			return nil, err
		}

		var tokens []string
		for _, token := range acc.GetCoins() {
			name := token.Denom
			if strings.Contains(name, "-") && !strings.Contains(name, "_") {
				tokens = append(tokens, name)
			}
		}
		addrTokenMap[addr] = tokens
	}

	tokens := addrTokenMap[addr]
	if len(tokens) <= 1 {
		return nil, fmt.Errorf( "don't have more than two types of token %v\n", tokens)
	}
	return tokens, nil
}

func pickTwoTokensRandomly(tokens []string) (string, string) {
	rand.Seed(time.Now().UnixNano() + rand.Int63n(10000))

	var name1, name2 string
	name1 = tokens[rand.Intn(len(tokens))]
	if name1List, ok := tokenPairsMap1[name1]; ok {
		for _, token := range tokens {

		}
	}
	if name2 == "" {
		if tmp, ok := tokenPairsMap2[name1]; ok {
			for _, name := range tmp {
				for _, token := range tokens {
					if token == name {
						name2 = name1
						name1 = name
					}
				}
			}
		}
	}

	return compareTokenNames(name1, name2)
}

func compareTokenNames(name1, name2 string) (string, string) {
	if name1 > name2 {
		return name2, name1
	} else {
		return name1, name2
	}
}