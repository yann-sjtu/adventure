module github.com/okex/adventure

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.9.25
	github.com/mitchellh/mapstructure v1.1.2
	github.com/okex/exchain v0.17.5-0.20210422073932-a732f21d7264
	github.com/okex/exchain-go-sdk v0.17.2-0.20210512105333-9b72ded383b7
	github.com/spf13/cobra v1.1.1
	github.com/tendermint/tendermint v0.33.9
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.0.0-20210418004230-26125dd4b59d
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.3-exchain
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-okexchain6.0.20210422060127-5c29d4d529c3
)
