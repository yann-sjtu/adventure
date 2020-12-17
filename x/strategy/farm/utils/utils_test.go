package utils

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/okex/adventure/common"
)

func TestGetRandomTokenSymbol(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(GetRandomBool())
	}
}

func TestBuildLPTName(t *testing.T) {
	fmt.Println(BuildLPTName("eth-112", common.NativeToken))
	fmt.Println(BuildLPTName("usdt-112", common.NativeToken))
}

func TestBuildLPTName1(t *testing.T) {
	//files, err := ioutil.ReadDir("./template/mnemonic/farm_test")
	files, err := ioutil.ReadDir("..")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name())
		}
	}
}
