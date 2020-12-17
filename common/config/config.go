package config

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/toml"
)

var (
	GlobalConfigPath = "/Users/oker/go/src/github.com/okex/adventure/config.toml"
	TxConfigPath     = ""
)

var Cfg *Config

type Config struct {
	TxConfigPath string
	Hosts        []string `toml:"hosts"`
	LogLevel     int8     `toml:"log_level"`
	Order        Order        `toml:"order"`
	Distribution Distribution `toml:"distribution"`
	Staking      Staking      `toml:"staking"`
	Token        Token        `toml:"token"`
}

func init() {
	if _, err := toml.DecodeFile(GlobalConfigPath, &Cfg); err != nil {
		log.Fatal(err)
	}
	Cfg.TxConfigPath = TxConfigPath
}

func GetConfig() *Config {
	return Cfg
}

func pickPort() int {
	defaultPort := 2333
	for port := defaultPort; port <= 65535; port++ {
		checkStatement := fmt.Sprintf("lsof -i:%d ", port)
		output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
		if len(output) > 0 {
			continue
		}
		return port
	}
	log.Fatal(fmt.Errorf("there is no vancant port"))
	return defaultPort
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
TxConfigPath: %s
Hosts: %v
`, c.TxConfigPath, hosts)
}
