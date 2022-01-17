module github.com/okex/adventure

go 1.14

require (
	github.com/ethereum/go-ethereum v1.10.8
	github.com/okex/exchain v1.1.4
	github.com/okex/exchain-ethereum-compatible v1.1.0
	github.com/okex/exchain-go-sdk v1.1.3-0.20220116091207-2ab92573a277
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/ethereum/go-ethereum => github.com/okex/go-ethereum v1.10.8-oec2
	github.com/tendermint/go-amino => github.com/okex/go-amino v0.15.1-exchain2
	github.com/tendermint/tm-db => github.com/okex/tm-db v0.5.2-exchain4
)
