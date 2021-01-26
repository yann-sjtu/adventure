package farm_unlock

import (
	"log"
	"time"

	"github.com/okex/adventure/common"
	common2 "github.com/okex/adventure/x/monitor/common"
	farm_control "github.com/okex/adventure/x/monitor/farm-control"
	"github.com/spf13/cobra"
)

var (
	poolName   = "1st_pool_okt_usdt"
	lockSymbol = "ammswap_okt_usdt-a2b"
)

func FarmUnlockCmd() *cobra.Command {
	farmUnlockCmd := &cobra.Command{
		Use:   "farm-unlock",
		Short: "farm-unlock",
		Args:  cobra.NoArgs,
		RunE:  runFarmUnlocklCmd,
	}

	//flags := farmUnlockCmd.Flags()
	return farmUnlockCmd
}

func runFarmUnlocklCmd(cmd *cobra.Command, args []string) error {
	clientManager := common.NewClientManager(common.Cfg.Hosts, common.AUTO)

	cli := clientManager.GetClient()
	for i := 0; i < len(addrs); i++ {
		index := i + 901
		lockinfo, err := cli.Farm().QueryLockInfo(poolName, addrs[i])
		if err != nil {
			log.Printf("[%d]%s failed to query lock info. err: %s", index, addrs[i], err.Error())
			continue
		}

		if !lockinfo.Amount.IsZero() {
			acc, err := cli.Auth().QueryAccount(addrs[i])
			if err != nil {
				log.Printf("[%d]%s failed to query account info. err: %s", index, addrs[i], err.Error())
				continue
			}

			unlockMsg := farm_control.NewMsgUnLock(acc.GetAccountNumber(), acc.GetSequence(), lockinfo.Amount, addrs[i])
			err = common2.SendMsg(common2.Unfarmlp, unlockMsg, index)
			if err != nil {
				log.Printf("[%d] %s failed to unlock: %s\n", index, addrs[i], err)
				continue
			}
			log.Printf("[%d] %s send lock msg: %+v\n", index, addrs[i], unlockMsg.Msgs[0])
		}

		time.Sleep(time.Duration(10) * time.Second )
	}
	return nil
}

var addrs = []string{
	"okexchain1hej6n80vs262wmrpvmkf9zmugqgjfqzydxa8et",
	"okexchain1mcut93efcpdjgjkpfvemfvvqju0zs8ccj3u0xh",
	"okexchain1wc49x2rnml97d537vsv4wqj0e67y90xe0y4nzm",
	"okexchain1992paeg9p7v4kkvqkp6qynk7ujy3220gw7l9ps",
	"okexchain10gclna4ds064tpxq74s3jfyyjvnylkdkjftvyd",
	"okexchain1u56l3ehphegvgx0557wyawqggcw7m8j36sfwwa",
	"okexchain1g82hlllygaf6rnnsaxqdl0xxmue2fwt2j9hdkf",
	"okexchain1hg5synr7qxqsyc0gj2r0hvtdf0kntfsl73xp23",
	"okexchain1ah6fu38g6nm9rmksa7uc6hn4qyu4nah8335900",
	"okexchain19s7zpd3jc8wtt7uwz67n7703fzudsrpeam2jte",
	"okexchain15uu9nutl9x0kyvwam3yz9hxkz83cglggnxzel5",
	"okexchain120rxcq6w9y46qa2ewyzuxvx0t74y2ch506kl2h",
	"okexchain1nlve0y2tjgwskpxrr4c3lywyklx2xdylmxsdl7",
	"okexchain12mq8qa5vfqmwmjhekq47a8hfxlp7ft2vl077kd",
	"okexchain1mwmgfge472httc7n2vznlmu3q6vfd7h3ft3ewy",
	"okexchain1kd9j6gm05c3txf4ndfcqr4raqq347g2e2tq90q",
	"okexchain1c7rkel67fucyq7j2zuzh42762tsn7s9gtpwyuz",
	"okexchain1zyufn2zm8az8664ed07muskktz2qmkuymf3ye3",
	"okexchain1ztp8hz25hqt5pcvym698vnqglfu8kfxhp30x05",
	"okexchain17nxv3awzz2g8u99xntwqgdy9j0x5scs6wt25xv",
	"okexchain1czwtauev2j03ltgwal3gpts95r84k26ktj282t",
	"okexchain1qzl0wl5kr7qk25wn5evr0zzfqrsjcrj023wc6n",
	"okexchain1grkf3mtfzlvx6kdsdt48jtyt4v0pa5yvc7dxt5",
	"okexchain1r99n9mrda79jdktr0xxmza87k97v5eh0aq93fm",
	"okexchain1ln38mfpx5vuugk85grljw8c4utcechdm3v55xp",
	"okexchain15v8k8gfp2paxrpaw98mnf9pfycgr4xard3u8yr",
	"okexchain1sey3syrw0xsqsvkdlk86j72w0hfnq8j4yqk9e3",
	"okexchain183quv4zattal80tkq3ccnxntv3yrpt6yyjt6sa",
	"okexchain1jgjvsut6lpwhkm2zrmvdl80k2qy8ljx6qnxdkg",
	"okexchain1ydse3j6l6rjj30anu0jajmr8gtdyzpgehe5unq",
	"okexchain1rm2hey38tdgupmf0lfqy9zhvd9aqp4x20kus0e",
	"okexchain1dhmj6639ctse5cqyxnrwghpjakktqxvefq9rg2",
	"okexchain1mqhh0tfqaxm87mvd7szzry80ll2ydrlt2ndczk",
	"okexchain1tpg9mmvqw4j97z0kgwuz4xum7swva9s8j5qzlc",
	"okexchain1tts22z7zrd80shyelpxw3h4p8vgfxlp0y9vaf0",
	"okexchain12rkvgnsge0py6dfcyhnu82zn334nne3skv589m",
	"okexchain1lhk92ptwwrncsfxlg4cflczyp56k29ffd38hkn",
	"okexchain1g6nemwdsgzn6t8ylkmrual75s3rr9qw8an9dyc",
	"okexchain1xradykga8gea73hv9ypu2ch9nw0kln7p6eehft",
	"okexchain190e2jqcfp8de6v46eh8qqxk4v5lz506whwm4d3",
	"okexchain10l4ja2vzenvqrwdfklsexnymkrvryxqsc03lx2",
	"okexchain1kwxnhcngphdnm4qm5ly9wg3x686pzfvl2lsvv9",
	"okexchain1ug9rkhu95xy4m30dj2h0l88rltp2ga7cp8c43z",
	"okexchain1hqxqmlfs9khgqmef57xq242c7kdha8nsjeyxcv",
	"okexchain157p3dta442g9cav3l0g5ws4rr79al3rpvls0ju",
	"okexchain17rkgqreruk9wchyf4a62n32g82sngnp6sjc0dz",
	"okexchain1rqpsn8na63t07dj76kdmye2cym085kp0hqehmx",
	"okexchain1k9gf34evr9frnkw5vcnl675h4uz5yt7509z0uu",
	"okexchain1t07asmpda3ekk88x8cu9rww4g25t9er0uvep0k",
	"okexchain1qnysv44vfkp8745kt076ryqv38zwrqlv56gmgj",
	"okexchain189sq8hphj3kzp8a302kk48r7m4f2kq4z2vu0u7",
	"okexchain1e7ltr6mp8de7dznua9xua2carf27nh3pqlg9l7",
	"okexchain1h4t9z7amss2tmy07efngjez3zrpe7zrg4k95kp",
	"okexchain1ujtt9q5v9xtl9k64ush999c78v3nepr4x63qpd",
	"okexchain1h5hxzjd2v2p3h4plae9y9qkvlvn4vx8c0xakzy",
	"okexchain1cv2vv36kk8adk2rve0766lwr6q50qsg7se2x03",
	"okexchain13stegj0ycrtcqzwuntdynk75x3qcq92579qfus",
	"okexchain1w2pmj4ghzc3d8svas20xcxhwnyxtna9q4f4vzc",
	"okexchain1wauwe45eckj38cn6d9f0whzqveswelgpxhe2hg",
	"okexchain1gm7vrklxvkl62wtzuy3j02uc33hr0xmyd8jjns",
	"okexchain19z2jzft3y8dlkeaxpnraccrdfxn0uz079kwfvy",
	"okexchain1q7a7kemn60tlsw2qt4tft0az3ffw0l03fchumy",
	"okexchain146dh4fw7a9qycqhagd7zwkj2n833n0tx8gtwy9",
	"okexchain13tughvphe8mvqrkkugwzp06glhrgkpzacfm8ag",
	"okexchain1esyfgc0dmtnmp3gkrcthl8lu2zsn5gmvuq6d5v",
	"okexchain1nydktjzyfynd97rm4x9gdp3hywjyevdu4red6g",
	"okexchain1ylnnp9r8tg4w0fdynrfyr2gq790c9u2crhqtpy",
	"okexchain173qndaykw3y3x4fr9x275age9aryesla7jz4r4",
	"okexchain1p9uqj3rzmk9r9pkx05fdnre7d2kwhtan7tl3xm",
	"okexchain1kz39t79evg5y7fkq8m4yexnyuzt8lzx40gk7es",
	"okexchain177mqkht4fqn9q35tlcc2qnpn0ak5xvjcfpq662",
	"okexchain1czrqkgcrf75mejth7wm9sxdfw2ykwj9amd6prj",
	"okexchain1my2vm88cnet0gjq2a7dw6zrhqm5656c352ym4a",
	"okexchain1rmpk5rmsyagakdxx7t8xny8eglu6lp5dvj8g4w",
	"okexchain1yw6qx8dudxpkeghdh7n8z300e4svxtzrk2qc6j",
	"okexchain17zgmngpkc07f9ug5xee5aj7y5duev2k42mh58y",
	"okexchain10m7am5qtxutps94jvwejwklcfw0z9g9nxzj7jz",
	"okexchain17v7hx5pryjqnpqvljhvxayvd2xd2ljh0hjw200",
	"okexchain186v5e9ke80jaj8je0dnvh2q7w0dl46ayyulv2s",
	"okexchain1l94x89s6d65ffzzt2ns8lyr9d958u7dkm25zc0",
	"okexchain140sy6yc4mx7equhrzavm20snzm5vd9va6f5rfw",
	"okexchain16ru7xvv39wjyppvczkqvtl5hykkm527uxdc05z",
	"okexchain10r3njv6pgln59narmklrdswfh2yj7hkhrwf76a",
	"okexchain1wffed7zcp30f9rzzxss4wdk523t9j29kljrnqq",
	"okexchain1xymyqst3dxl0p8ffy9mfeuqwc79kf6rc5t8kcn",
	"okexchain1exutduraqy5gpg7m3ptrw2tc99sktm72m606mn",
	"okexchain1rfhhsq0rqfa8erlnxn3cmtt9clrf22kcy02pje",
	"okexchain18pkm7sr0g72l3ekmqjj90y3fceflkz547juezh",
	"okexchain1cgrd77q5f0wg78vhhk38pasl2gpnqlpg5ew4wm",
	"okexchain1jr06l7fce79xx2pt7pzquhqt355fm0g920r3y7",
	"okexchain1fuh850ml6vcrn8lkckpewpgpau2dfuyjcyqvaq",
	"okexchain1mpmw2jx0dfnrh73tj0xzcc79n5ddsp5s323y03",
	"okexchain1fqtg6rs73jfs6gz6zndnc0ufy8sm0wk8d0vzuu",
	"okexchain1xjjw2z20zaym8h3grfz4gnf8l2w740npwptjqg",
	"okexchain17996fzeu9psv7n4yvp227m3rjmdrag5eyg8lkf",
	"okexchain1edsq5rfum9pkq96pt0efs7jesc9xdcyq3hx4a2",
	"okexchain1h47k29q4vthecc8hnkyrwjf5pxx8cwcj5jfw3r",
	"okexchain18x02ulv7mceswsrfjgt4mscc8nkxqqy0r8626s",
	"okexchain12wp7k6vk6rfzpldk8qs936sm5p5u36qzsjjjrp",
	"okexchain1pa345e0qmn7hrhfyv4a7ys3zhshcnetlx008qs",
}