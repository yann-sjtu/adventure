module github.com/okex/adventure

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.9.24
	github.com/mitchellh/mapstructure v1.1.2
	github.com/okex/okexchain v0.15.0-testnet.0.20201217024321-28df62c5c23c
	github.com/okex/okexchain-go-sdk v0.11.1-0.20201217061003-d95d95f115ad
	github.com/spf13/cobra v1.1.1
	github.com/tendermint/tendermint v0.33.9
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.3-0.20201216093610-29445e0bde54
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-okexchain1
)
