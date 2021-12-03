module github.com/okex/adventure

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.10.8
	github.com/okex/exchain v0.19.17
	github.com/okex/exchain-ethereum-compatible v1.0.2
	github.com/okex/exchain-go-sdk v0.19.17
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/cobra v1.1.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/tendermint/go-amino => github.com/okex/go-amino v0.15.2-0.20211120080923-5ccee577b032
	github.com/tendermint/tm-db => github.com/okex/tm-db v0.5.3-0.20211118074832-0d08619b97fe
)
