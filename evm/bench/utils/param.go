package utils

import (
	"github.com/okex/adventure/common"
	"github.com/okex/adventure/evm/constant"
	"github.com/spf13/viper"
)

type BasepParam struct {
	sleep       int
	concurrency int
	ips         []string

	privateKeys []string
}

func NewBaseParam(sleep int, concurrency int, ips []string, privateKeyFile string) BasepParam {
	privateKeys := constant.PrivateKeys
	if privateKeyFile != "" {
		privateKeys = common.ReadDataFromFile(privateKeyFile)
	}
	return BasepParam{
		sleep,
		concurrency,
		ips,
		privateKeys,
	}
}

func DefaultBaseParamFromFlag() BasepParam {
	return NewBaseParam(
		viper.GetInt(constant.FlagSleep),
		viper.GetInt(constant.FlagConcurrency),
		viper.GetStringSlice(constant.FlagIPs),
		viper.GetString(constant.FlagPrivateKeyFile),
	)
}

func (bParam BasepParam) GetSleep() int {
	return bParam.sleep
}

func (bParam BasepParam) GetConcurrency() int {
	return bParam.concurrency
}

func (bParam BasepParam) GetIPs() []string {
	return bParam.ips
}

func (bParam BasepParam) GetPrivateKeys() []string {
	return bParam.privateKeys
}
