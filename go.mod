module github.com/okex/adventure

go 1.14

require (
	github.com/ethereum/go-ethereum v1.10.8
	github.com/okex/exchain v1.2.1-0.20220414063812-5b9519a911d7
	github.com/okex/exchain-go-sdk v1.1.3-0.20220415153509-234679002b68
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/ethereum/go-ethereum => github.com/okex/go-ethereum v1.10.8-oec3
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/go-amino => github.com/okex/go-amino v0.15.1-exchain6
)
