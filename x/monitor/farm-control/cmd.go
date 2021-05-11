package farm_control

import (
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/adventure/common"
	monitorcommon "github.com/okex/adventure/x/monitor/common"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/spf13/cobra"
)

func FarmControlCmd() *cobra.Command {
	farmControlCmd := &cobra.Command{
		Use:   "farm-control",
		Short: "farm controlling in the first miner on a farm pool",
		Args:  cobra.NoArgs,
		RunE:  runFarmControlCmd,
	}

	flags := farmControlCmd.Flags()
	flags.IntVarP(&sleepTime, "sleep_time", "s", 0, "sleep time of add-liquidity msg and lock msg")
	flags.IntVarP(&startIndex, "start_index", "i", 0, "account index")
	flags.StringVar(&poolName, "pool_name", "", "farm pool name")
	flags.StringVar(&lockSymbol, "lock_symbol", "", "token name used for locking into farm pool")
	flags.StringVar(&baseCoin, "base_coin", "", "base coin name in swap pool")
	flags.StringVar(&quoteCoin, "quote_coin", "", "quote coin name in swap pool")
	flags.BoolVarP(&toAddLiquidity, "to_add_liquidity", "a", true, "decide to add liquidity or not")
	flags.BoolVarP(&toLock, "to_lock", "l", true, "decide to lock lpt or not")

	return farmControlCmd
}

var (
	sleepTime = 0

	startIndex = 0

	poolName   = ""
	lockSymbol = ""
	baseCoin   = ""
	quoteCoin  = ""

	toAddLiquidity = true
	toLock         = true
)

func runFarmControlCmd(cmd *cobra.Command, args []string) error {
	accounts := monitorcommon.AddrsBook[startIndex/100]
	clientManager := common.NewClientManager(common.GlobalConfig.Networks[""].Hosts, common.AUTO)

	for _, account := range accounts {
		cli := clientManager.GetClient()
		if !filterAddr(cli, account.Address) {
			continue
		}
		fmt.Println()
		log.Printf("=================== %+v ===================\n", account)

		// 1. check the ratio of (our_total_locked_lpt / total_locked_lpt), then return how many lpt to be replenished
		requiredToken, err := calculateReuiredAmount(cli, accounts)
		if err != nil {
			fmt.Printf("[Phase1 Calculate] failed: %s\n", err.Error())
			continue
		}

		// 2. judge if the requiredToken is zero or not
		if requiredToken.IsZero() {
			// 2.1 our_total_locked_lpt / total_locked_lpt > 80%, then do nothing
			fmt.Printf("This Round doesn't need to lock more %s \n", lockSymbol)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		} else {
			fmt.Printf("There is %s to be replenished\n", requiredToken)
			// 2.1 our_total_locked_lpt / total_locked_lpt < 80%, then promote the ratio over 81%
			err = replenishLockedToken(cli, account.Index, account.Address)
			if err != nil {
				fmt.Printf("[Phase2 Replenish] failed: %s\n", err.Error())
				continue
			}
		}
	}

	return nil
}

func filterAddr(cli *gosdk.Client, addr string) bool {
	if i := monitorcommon.StringsContains(limitAddrs, addr); i == -1 {
		return false
	} else {
		if accInfo, err := cli.Auth().QueryAccount(addr); err != nil {
			return false
		} else if accInfo.GetCoins().AmountOf(baseCoin).LT(types.MustNewDecFromStr("1000")) {
			return false
		} else if accInfo.GetCoins().AmountOf(quoteCoin).LT(types.MustNewDecFromStr("200")) {
			return false
		}
		return true
	}
}

var limitAddrs = []string{
	"okexchain1wc49x2rnml97d537vsv4wqj0e67y90xe0y4nzm",
	"okexchain1992paeg9p7v4kkvqkp6qynk7ujy3220gw7l9ps",
	"okexchain10gclna4ds064tpxq74s3jfyyjvnylkdkjftvyd",
	"okexchain1u56l3ehphegvgx0557wyawqggcw7m8j36sfwwa",
	"okexchain19s7zpd3jc8wtt7uwz67n7703fzudsrpeam2jte",
	"okexchain1mwmgfge472httc7n2vznlmu3q6vfd7h3ft3ewy",
	"okexchain1mcut93efcpdjgjkpfvemfvvqju0zs8ccj3u0xh",
	"okexchain1tts22z7zrd80shyelpxw3h4p8vgfxlp0y9vaf0",
	"okexchain1ztp8hz25hqt5pcvym698vnqglfu8kfxhp30x05",
	"okexchain1fqtg6rs73jfs6gz6zndnc0ufy8sm0wk8d0vzuu",
	"okexchain12mq8qa5vfqmwmjhekq47a8hfxlp7ft2vl077kd",
	"okexchain1nlve0y2tjgwskpxrr4c3lywyklx2xdylmxsdl7",
	"okexchain12rkvgnsge0py6dfcyhnu82zn334nne3skv589m",
	"okexchain140sy6yc4mx7equhrzavm20snzm5vd9va6f5rfw",
	"okexchain1sey3syrw0xsqsvkdlk86j72w0hfnq8j4yqk9e3",
	"okexchain1mpmw2jx0dfnrh73tj0xzcc79n5ddsp5s323y03",
	"okexchain183quv4zattal80tkq3ccnxntv3yrpt6yyjt6sa",
	"okexchain1zyufn2zm8az8664ed07muskktz2qmkuymf3ye3",
	"okexchain120rxcq6w9y46qa2ewyzuxvx0t74y2ch506kl2h",
	"okexchain1l94x89s6d65ffzzt2ns8lyr9d958u7dkm25zc0",
	"okexchain1tpg9mmvqw4j97z0kgwuz4xum7swva9s8j5qzlc",
	"okexchain1g82hlllygaf6rnnsaxqdl0xxmue2fwt2j9hdkf",
	"okexchain189sq8hphj3kzp8a302kk48r7m4f2kq4z2vu0u7",
	"okexchain1h4t9z7amss2tmy07efngjez3zrpe7zrg4k95kp",
	"okexchain17rkgqreruk9wchyf4a62n32g82sngnp6sjc0dz",
	"okexchain15v8k8gfp2paxrpaw98mnf9pfycgr4xard3u8yr",
	"okexchain1ah6fu38g6nm9rmksa7uc6hn4qyu4nah8335900",
	"okexchain146dh4fw7a9qycqhagd7zwkj2n833n0tx8gtwy9",
	"okexchain19z2jzft3y8dlkeaxpnraccrdfxn0uz079kwfvy",
	"okexchain157p3dta442g9cav3l0g5ws4rr79al3rpvls0ju",
	"okexchain1hg5synr7qxqsyc0gj2r0hvtdf0kntfsl73xp23",
	"okexchain1cv2vv36kk8adk2rve0766lwr6q50qsg7se2x03",
	"okexchain1yw6qx8dudxpkeghdh7n8z300e4svxtzrk2qc6j",
	"okexchain1rmpk5rmsyagakdxx7t8xny8eglu6lp5dvj8g4w",
}
