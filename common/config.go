package common

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	GlobalConfigPath = "" //TODO
	TestCasesPath    = ""
)

var Cfg *Config

type Config struct {
	TestCaesPath string
	Hosts        []string            `toml:"hosts"`
}

func init() {
	if _, err := toml.DecodeFile(GlobalConfigPath, &Cfg); err != nil {
		log.Fatal(err)
	}
	Cfg.TestCaesPath = TestCasesPath
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
Hosts: %v
`, c.TestCaesPath, hosts)
}
