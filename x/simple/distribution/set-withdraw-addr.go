package distribution

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/okex/adventure/common"
	gosdk "github.com/okex/okexchain-go-sdk"
)

const setWithdrawAddr = common.SetWithdrawAddr

func SetWithdrawAddr(cli *gosdk.Client, info keys.Info) {
	accInfo, err := cli.Auth().QueryAccount(info.GetAddress().String())
	if err != nil {
		common.PrintQueryAccountError(err, setWithdrawAddr, info)
		return
	}

	_, err = cli.Distribution().SetWithdrawAddr(info, common.PassWord,
		pickRandomAddress(),
		"", accInfo.GetAccountNumber(), accInfo.GetSequence())
	if err != nil {
		common.PrintExecuteTxError(err, setWithdrawAddr, info)
		return
	}
	common.PrintExecuteTxSuccess(setWithdrawAddr, info)
}

func pickRandomAddress() string {
	rand.Seed(time.Now().UnixNano())
	length := len(backupAddress)
	return backupAddress[rand.Intn(length)]
}

// related to mnemonic/normal_100
var backupAddress = []string{
	"okexchain1xufxgnvctu35d8wyxxl8xshzu0ny9uv38v0p9n",
	"okexchain1y6qcerk9u8whxwxyfupuc2jxanp6hfcc99dl65",
	"okexchain1lnwla4jgcka8y026kxl5q75yq7zm0x0ry6l0c5",
	"okexchain1d435zl4sjcfyf5rm3a8mvl6uv8k923nwlmef4k",
	"okexchain1l9r3ll75hnt209z83e7pn3ts0uy7dhnywrj6wj",
	"okexchain1nq30l2m6tpg2gwt7fx8wkl2rjl4mgac4ywdhe8",
	"okexchain1vgfux8cffajthh9pcgapx36lrkd3pzyum987pe",
	"okexchain1vk2k252fu4qzkd75quu3t4v0x2fsjx2y4wmaxz",
	"okexchain1gtzf0dh988we4w0cvkj0anslg3g03vml56kq4d",
	"okexchain1q7tkc38wxpvsx8m5szl5cheq0luvscvuaytmh4",
	"okexchain1r8w538zz3s0l3h89mefa84mg0h5ayws5dkc0un",
	"okexchain1dlhrhjxc733f733j0ekh83pc62ndpkfxe86kwx",
	"okexchain1fcsgfun9djmmklslpca7ena5xwnseyys3qwzr4",
	"okexchain1pwhjgcxfd2c8a73uzp53e7w8tdswp3tpys2ewd",
	"okexchain1yn29k8sch9hhmahfw5m6yar76xnsep3c04erlr",
	"okexchain1uuprelkmgvr56vjevzmf0xud6zxckysm8644dk",
	"okexchain12dw93ynalcfcj588sal280xdtde2hn4weqyerv",
	"okexchain1c05nvxe9f6vmq0tsezxc66fkd589msqyj6r4f9",
	"okexchain1aj7xm6gdmwlw28x7ckjfwxvtfuav26cc6zalmf",
	"okexchain17q0vhu9qxxm8anyns5fz2shc86z7qk3sr2xcjz",
	"okexchain1f8zllmdg9yx7w2d8xytj2zt0khu5xk7fm2t2r0",
	"okexchain14xe657a2d6znxnhf8gjxg20djcy05kccjd5k87",
	"okexchain1qssu4q3m6el3c2nrg28ajs8lgn2cl78mdreymd",
	"okexchain1klwfuz7nzydrtllxjr0xvp8cq9ws8szk666dna",
	"okexchain12envmsg45ydjre4m67qgg6hqnq3scjmmzzt3u0",
	"okexchain1deqrgdddjwthg9ps9tztxtgdufuq9hm29dkwaq",
	"okexchain1ur970v9r94lwkzyllrenymw97zxlcau9cru3ry",
	"okexchain15q47vzjsewg3h3shj26js4cpaq06yrn30c8arr",
	"okexchain13x04xqh629zky6wlz5jq0jleflqmw96ygapync",
	"okexchain1m8ehymr7xk3w590mzwq0thqffn79zur7qsd7pk",
	"okexchain19zjkt5jg4uhqc7ea6j23km7grnh8zzvkqgawqg",
	"okexchain18sw879x9w7e6jfp3fcv7httpjjxqd7mcu4q3ss",
	"okexchain1wlwpzl9k986272tzakfph7x9mqzgde0crp8fxd",
	"okexchain1k4zkdf26gk7c22e948lc4tf90pz64fu4wm6c6d",
	"okexchain1qcw5yczn7ftuq6pqww7kxzpfh3tzmx24ra5lps",
	"okexchain1xlreler0s9q02tlv2jqhcdw73s0jj267h8d9hy",
	"okexchain1qsprr2wt2kxgpxjwr3qwkxh80sdmj2thpn7uc0",
	"okexchain168v62net3hf7qsc4a0taztgy3x75aps3h64d7x",
	"okexchain1jyzgwkpq7r4m2njvy2252zu0fm82hsahfq9g0n",
	"okexchain1989hpyyazjdy3ze78yd9ef2hpxstry07sjafe8",
	"okexchain19kka2lrhh75tpcwgpwe2pqfw47ry4n45pl9t96",
	"okexchain1ljh4d043fcgcxqwfv5trgp44d97j4a9janshcu",
	"okexchain1wxhxdfsvhatx5recnm7y530dp7say60qftw30q",
	"okexchain126d84scp2z94u67ytwkcwuxrhp3479jllyjym6",
	"okexchain195lwfcat5kn3antracu3924q2qd8mv9eclfgsr",
	"okexchain1qlqnecqwwnkdld3fl7fukxj2hmxlrpm57sn0k6",
	"okexchain15lp39vgexu60tpu43s56reggt9lg98mdv85xxv",
	"okexchain1p00gdzatkslpe4uat8mjs7zqguqstm8tkzu2z8",
	"okexchain1f7crmj3af5vk0xqeqhr48untuhn9hgextwnzr3",
	"okexchain1jpd5wq5kr6vccak7dncvqvudj0f2kk5urvmggq",
	"okexchain1eteqcvnunast9uwjqv0ta5r9jq9uhcwsflnu8x",
	"okexchain149vydaa8kkgquur2cwaa4t5e55vzgyeazxd038",
	"okexchain1slkh3y6atc2k94zjtuq8xtemwdv29cp3kls29c",
	"okexchain1j0udusvmjvu46uqwatjvmgw4z8ug38xdcex9qq",
	"okexchain1sz3eqfmrx40lugufnw5mcywc59qpwwf4s3wzdm",
	"okexchain1phj40fwsun2nqyvl6w044cq8tmsvu6f4u96eg3",
	"okexchain1l2nkqva7eq6dg0uxsc4nh6g344jx2djfvslm26",
	"okexchain1y6d4qure2xs64t8n3hqjzyhecm590lnkyrmm07",
	"okexchain1pss7ymv8ww580ldk35nrzrtqq69ktyy3jx4gvj",
	"okexchain1703hy58h0mcxth7x4h73kzcj793zgyc59y767m",
	"okexchain1h2hk3kw40e85xkgefy8x39wse5wry8h3ka2v3h",
	"okexchain1np88ag804764rnyr4a2c7f4jvvhx3h98pvzd2k",
	"okexchain1e69jxghqh75ypkt4mfptmv6v4hjfa86xzse4qg",
	"okexchain1g3zsf6javc27uk8znykgrle4npp7ksgszadey8",
	"okexchain15m75wpztm7vpzwz2een9ejfpe7clvz8kcl78kv",
	"okexchain1ay8ylym54l4hlkpfepe2l9r0gpggrhj4xx00m4",
	"okexchain1kwvw6dcud4ftk23tsl8nswu4mvka6hs3ecppvd",
	"okexchain18jdgadjv9eguyh5gp05qxtfqe0v55x6394zekc",
	"okexchain1tzv0gqjhqh9n39qn7yaza86gr34rjem9x4k2ws",
	"okexchain1jepjvuy7gpwwwtusms9wy5mpgzmyushvruxgf0",
	"okexchain1089336v0ghrmcuke8c0w8d4rajdqnm3uymcum3",
	"okexchain1c327043r746vezqeecemzp7jtu9xatyl4cxddx",
	"okexchain1q828jw4w03ckk0re9mu88vuln2zq0uqc97z8yt",
	"okexchain1yt5mn7cvr9ek2afrugrdjwpsq95j28jkj3f08w",
	"okexchain1rmmn2wk3rpwg7w9p8aqh73wcwgl4lhlwckgn9u",
	"okexchain16w22chqpuv92yhmhyegw4na72kutesktsjq7ne",
	"okexchain1rxlae89k9xtnx6l036v6hdt3hy3jrtthzkl33q",
	"okexchain1nakwrvqyzn5fzw4qx5wjuvmsme66x2q22yevz8",
	"okexchain1cy8fnc2zl9rl5ae2clmng4asf6tm367n92j73h",
	"okexchain1cyh49wtelaux4c37rptntnn89haxauzk6szmdv",
	"okexchain1w6f8unw09th2p36k5up0z6p4ht0zcwwpy7ql0g",
	"okexchain1t48tq6umyy89hsq2kqwsec245490el3wlmc8kx",
	"okexchain1kfjdaqxsdscjsee35hxqft02n2xa8g0fkaye7c",
	"okexchain1fkjl25dmx6r80yj420tr8v26r369l3mkl50gp6",
	"okexchain18dgh97husfqsxu9r9axn5cxm4p6yf5amfrf77d",
	"okexchain1swmntuz98pm05ru0r4fkaz9nyy9j9n59q4cdrz",
	"okexchain1w2cpjmj3qlx4udl0kt6stzsrv2r9d0cmcxdmdw",
	"okexchain1px6nvgfep9fpla6aesyux0ztxz9en355yqst3n",
	"okexchain1vfcju4j6qfhtn8mqmjt6gxkw7eur4keukwgudp",
	"okexchain1pf6yzkmnfqc0tnlnjlwyav0g08qnpgtm6wn2tm",
	"okexchain1l86sdg55fjv05c33vfe7sfvjjs68cpq3ngsy0q",
	"okexchain1l7kr6jutvlwx4yw4kkwasa25quz0psd4ug8tz9",
	"okexchain1kqqsh3qvvgjs2yv5sqjw4cg78tr9a03ah3l6zt",
	"okexchain1gaxsng4g8ve5h0pygpx8junmaug3rme42rtkw0",
	"okexchain1rnhgp540eu749ulduypf0qdswqf4z0eglug9q2",
	"okexchain1xv8rjh80cca2pje9tu6nd9vq5qz8jlrv36c2ve",
	"okexchain1fvr22awa9pj2wkhw6nn0q4s7j3k09066zwuw3k",
	"okexchain1rmp9rv8cn8v3fq77m8hcftxv6jkpj369ttzsqc",
	"okexchain1u35lq0mufrf96n2m0yu95fyr5lj34ctk7fjlq2",
	"okexchain13pyylyymk7jlul4vz04m2ypdjvcc9pxe75x34a",
}
