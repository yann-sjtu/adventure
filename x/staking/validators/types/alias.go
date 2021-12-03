package types

const (
	DefaultPasswd   = "12345678"
	DoNotModifyDesc = "[do-not-modify]"
)

var (
	ConsPubkeys = []string{
		"okchainvalconspub1zcjduepqtv2yy90ptjegdm34vfhlq2uw9eu39hjrt98sffj7yghl4s47xv7swuf0dx",
		"okchainvalconspub1zcjduepq6ml2709hf3j6hjddes4ns548l6muhkxlphlm2t3tf377he62ckeqjt5x9r",
		"okchainvalconspub1zcjduepq0epr5evlnzvtj63vvzw8jmw8vxxmpkjh2rld0j62ynnfv8lxcuzqr3w6k3",
		"okchainvalconspub1zcjduepqs6xfkt92nll9knum5an2uuffzm8wwwhk4hu2nxthf8dmjq5cm9jqcg5eal",
		"okchainvalconspub1zcjduepqyjszm3x3j5m66hhzy4pf0jxn4r53zww5kv4snj65stt9270aw9nqm3ezk2",
		"okchainvalconspub1zcjduepqz7npagdg6pw53ez26227yj6qh0hdkw5wekgqnkrgw0q4afycllhspxq7u6",
		"okchainvalconspub1zcjduepqsxzsjn43tmy7hrymllq2kr88yhnr8f3vn3rk3wefma46rhxhyntqjmpczl",
		"okchainvalconspub1zcjduepqw0v9ehcqz8f2zcmx2raf3judmx0h0qk6hwcaf6m3rqw8za6xvcpqzxtvve",
		"okchainvalconspub1zcjduepq94k2ekkkx8e44qgmwhmwsnh9gxmq6pujzsx7tqrf789h0xjk046qndjarq",
		"okchainvalconspub1zcjduepqm5xp65zns2pepzyg3ry58newrlnxnd70hj23njzadngcnua8v35qfdh9w3",
		"okchainvalconspub1zcjduepqw2rqc8227ng24cjd0ccjvzfxn35fc8fnd4cg2n6kyfrgd82r83lqpk3u4k",
		"okchainvalconspub1zcjduepqyyy4g578mg08nxr6uf6vl9l0wzq2k25u86pmfuf9gyfuqt5gv3qqhhduct",
		"okchainvalconspub1zcjduepq8jr56mhvmjh9xeqfmwf3w3mwfwzp732murf9erfgsxnvpecc7aqqw6pjnk",
		"okchainvalconspub1zcjduepqdtckh0697kdghss24hekfnvwpzhu4ac4dj5uygzef042jl7f6u2q0h499j",
		"okchainvalconspub1zcjduepqqqhnky9y2a94mt6fxdz4wt2qwgnkthzfwh6d93g6wemqzcscppfq4evkdk",
		"okchainvalconspub1zcjduepq9eeda9uryxrk4904hqj8mx83y4750zl6t0klt46jp79pepuz7p7qz794vr",
		"okchainvalconspub1zcjduepqwh5jdduhq2udpyn79qd0kev0ar2q5q9j7nx7k3d3p0rlshs0sunqkvev5n",
		"okchainvalconspub1zcjduepqppt6nxrh5jnqrffyj8f7wq0zzge7qh6lcf5ruu704pzd8qr8zprsafhetk",
		"okchainvalconspub1zcjduepq9x68382ue84xds566tmsnfrk250dzwx7uyw9vty0r6qrrcyn5weqe0peq4",
		"okchainvalconspub1zcjduepqc4r4sccsrfkjse79lf8gumd8mej24vv9q2lt0za7ut0tdde3l4fsw45hv6",
		"okchainvalconspub1zcjduepq8hl6slpltrkfw7vlk27epkasvuensh0zkjk44c00wtj6tgk2dy8shm6qc4",
		"okchainvalconspub1zcjduepqhl0w656tm7ynpzdhg6y0prl2khpykz242c4wuxr8u59r2jdvf4rs8r5gp0",
		"okchainvalconspub1zcjduepq6a0qxr524jr4dg9yfptezv0l5j302uqrmamdwcs3pevu5kd77eyqgg9mph",
		"okchainvalconspub1zcjduepqp4pruuucsj5ag3kpn7xq95nendh5kyulv8nd6qwm3zh6dlvyxsnqaud000",
		"okchainvalconspub1zcjduepq24cxfcm6vnd5mmw50qyzclljfc3valm5af5te3v05jhqccfyljks5v4a5r",
		"okchainvalconspub1zcjduepq6sdy84x32d0nfrmxele69h8ccsdvnay4np3z9t5s3gyjl730033s0xkunw",
		"okchainvalconspub1zcjduepqh5xu3cp6rc528dyd76d2m8j7cvfhanlzxnkut605padyktxq8vlsnesnlq",
		"okchainvalconspub1zcjduepqsksyr5zu7v8nxfpkzmt6qelajpf0dz56fu75d306t4jkhkeref2qdapet3",
		"okchainvalconspub1zcjduepqrad3u2g2x32h9nqweh3kapd7ft6f8kwjgqq5uz8mykckkyvqqdps26ltya",
		"okchainvalconspub1zcjduepqquz9wwvzlmhm0dx5f773nzvwy7l03wqus48azv60gapv9397cklqty98fg",
		"okchainvalconspub1zcjduepq6x93hg9ucc4q6kud5ncxfcj2mf5e7nwf484qyx4krx8n30h7sacqg7xty5",
		"okchainvalconspub1zcjduepqntchpqcg0zvn6984rrkga25994ahfrvp4r4pfjfgk5gravegxn9qzv6h58",
		"okchainvalconspub1zcjduepqptxcr020t7tkh62nz98mv9qdujn30kmp23pvvpl0wq8kdujrrumszvtf63",
		"okchainvalconspub1zcjduepqeet9wljspnncayggf9vkx40ng8frwwgmlcymagrzsh39zrx45w4sx0xx8f",
		"okchainvalconspub1zcjduepqyt9zck777mn9v7ml2l5txje84m8eddud90u8hvgg29rgu4pvs92qke2gkv",
		"okchainvalconspub1zcjduepqjsycnggt0t7zrt8yaahad2lq5v7rrttky24kd5kfg6hqm4m0srrs3kf607",
		"okchainvalconspub1zcjduepq83ajw3lsxdf7t3ca8lyxcftgjapm70d7fxznt7jelphv6dj4lq3sekxy6k",
		"okchainvalconspub1zcjduepqtuh6esph5ja2h72fpvn02948wr239fmgyn8lyxra4ce4ytydp5xsaczmrw",
		"okchainvalconspub1zcjduepqfjwunern3qynmmglfzhdrtc9f5ghha46kmfxchw4grfqts0hk7zs5lhnr8",
		"okchainvalconspub1zcjduepq6xn876chsep9l3suehl4ycz7njnyj6davn56n7spdmfp79zcvv9qcrgyes",
	}
)