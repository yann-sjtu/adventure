package common

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	MainnetAli = "mainnet-ali"
	TestnetAli = "testnet-ali"
	TestnetAws = "testnet-aws"
	Localnet   = "local"
)

var GlobalConfig Config

type Config struct {
	Networks map[string]Network `toml:"networks"`
}

type Network struct {
	TestCaesPath string
	Hosts        []string `toml:"hosts"`
	ChainId      string   `toml:"chain-id"`
}

func InitConfig(path string) {
	if _, err := toml.DecodeFile(path, &GlobalConfig); err != nil {
		panic(err)
	}
	return
}

func GetConfig() Config {
	return GlobalConfig
}

func DecodeConfig(path string) (cfg Config) {
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		panic(err)
	}
	return
}

func (c *Network) String() string {
	var hosts string
	for i, host := range c.Hosts {
		hosts += host + " "
		if (i+1)%4 == 0 && (i+1) != len(c.Hosts) {
			hosts += "\n       "
		}
	}
	return fmt.Sprintf(`⚙️⚙️⚙️⚙️⚙️ Golbal Config ⚙️⚙️⚙️⚙️⚙️
TestCasesPath: %s
Chain-id: %s
Hosts: %v
`, c.TestCaesPath, c.ChainId, hosts)
}
