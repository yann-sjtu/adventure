module github.com/okex/adventure

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.9.25
	github.com/mitchellh/mapstructure v1.1.2
	github.com/okex/okexchain v0.16.8
	github.com/okex/okexchain-go-sdk v0.16.1-0.20210301080529-4e34315a8ea9
	github.com/spf13/cobra v1.1.1
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48
	github.com/tendermint/tendermint v0.33.9
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-okexchain7
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.1-okexchain1
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-okexchain4
)
