package common

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	GlobalConfigPath = "" //TODO
)

var Cfg *Config

type Config struct {
	TestCaesPath string
	Hosts        []string `toml:"hosts"`
	ChainId      string   `toml:"chain-id"`
}

func init() {
	if _, err := toml.DecodeFile(GlobalConfigPath, &Cfg); err != nil {
		log.Fatal(err)
	}
}

func GetConfig() *Config {
	return Cfg
}

func (c *Config) String() string {
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
