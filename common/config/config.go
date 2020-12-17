package config

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap/zapcore"
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
	LogListenUrl string
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
	Cfg.LogListenUrl = fmt.Sprintf("localhost:%d", pickPort())
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
LogLevel:     %s          (set the logger dynamically: curl -XPUT --data '{"level":"error"}' http://%s/handle/level )
LogListenUrl: %s
`, c.TxConfigPath, hosts, getLevel(c.LogLevel), c.LogListenUrl, c.LogListenUrl)
}

func getLevel(i int8) string {
	switch i {
	case int8(zapcore.DebugLevel):
		return "debug"
	case int8(zapcore.InfoLevel):
		return "info"
	case int8(zapcore.WarnLevel):
		return "warn"
	case int8(zapcore.ErrorLevel):
		return "error"
	case int8(zapcore.DPanicLevel):
		return "dpanic"
	case int8(zapcore.PanicLevel):
		return "panic"
	case int8(zapcore.FatalLevel):
		return "fatal"
	}
	return ""
}
