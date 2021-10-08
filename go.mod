module github.com/okex/adventure

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.10.8
	github.com/mitchellh/mapstructure v1.1.2
	github.com/okex/exchain v0.19.9
	github.com/okex/exchain-ethereum-compatible v1.0.2
	github.com/okex/exchain-go-sdk v0.19.0
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/tendermint/tendermint v0.33.9
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-exchain17
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.3-exchain2
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-exchain13
	github.com/tendermint/tm-db => github.com/okex/tm-db v0.5.2-exchain1
)
