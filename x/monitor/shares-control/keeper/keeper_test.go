package keeper

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"testing"
)

func TestTomlStruct(t *testing.T) {
	var config TestConfig
	str, _ := os.Getwd()
	fmt.Println(str)
	if _, err := toml.DecodeFile("./test_config.toml", &config); err != nil {
		return
	}

	fmt.Println(config)

}

type TestConfig struct {
	List [][]string `toml:"list"`
}
